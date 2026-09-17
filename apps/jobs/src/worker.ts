import { nowIso } from "@recally/domain";
import { dispatchJob, type WorkflowBinding } from "@recally/platform-cloudflare";
import { claimJobInstance, dueOutbox, getJob, markOutboxDispatched } from "@recally/storage";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import * as HttpServerRequest from "effect/unstable/http/HttpServerRequest";
import * as HttpServerResponse from "effect/unstable/http/HttpServerResponse";
import { definedModelVars, ModelVarsConfig } from "../../../infra/model-env";
import { ArchiveBucket, Database } from "../../../infra/resources";
import type { JobsEnv } from "./env";
import { EnrichWorkflow } from "./workflows/enrich";
import { IndexWorkflow } from "./workflows/index-wf";
import { IngestWorkflow } from "./workflows/ingest";
import { DigestWorkflow, ExportWorkflow, PurgeWorkflow } from "./workflows/purge";

const KIND_BINDINGS: Record<string, string> = {
  ingest: "ingest",
  enrich: "enrich",
  index: "index",
  digest: "digest",
  export: "export",
  purge: "purge",
};

// Adapt an alchemy workflow handle to the platform-neutral WorkflowBinding
// seam used by dispatchJob — duplicate-instance errors resolve to a ref
// instead of failing the sweep.
const toWorkflowBinding = (
  // Structural adapter over heterogeneous workflow handles.
  handle: Cloudflare.Workflows.WorkflowHandle<any, any>,
): WorkflowBinding => ({
  create: (options) => Effect.runPromise(handle.create(options)),
  get: (id) =>
    Effect.runPromise(
      handle.get(id).pipe(
        Effect.map((instance) => ({
          status: () =>
            Effect.runPromise(instance.status()).then((s) => ({ status: s.status as string })),
        })),
      ),
    ),
});

// Outbox catch-up: dispatch any rows whose workflow create failed or never
// ran (§9.3). Deterministic instance ids make re-dispatch safe.
const sweepOutbox = (db: JobsEnv["DB"], bindings: Record<string, WorkflowBinding>) =>
  Effect.tryPromise(async () => {
    const now = nowIso();
    const due = await dueOutbox(db, now, 20);
    console.log(`outbox sweep: ${due.length} due`);

    for (const row of due) {
      const job = await getJob(db, row.library_id, row.job_id);

      if (!job || job.status === "cancelled") {
        await markOutboxDispatched(db, row.id);
        continue;
      }

      const bindingKey = KIND_BINDINGS[job.kind];

      if (!bindingKey) continue;

      try {
        const instance = await dispatchJob(bindings[bindingKey]!, job.id, job.generation, {
          jobId: job.id,
          libraryId: job.library_id,
        });

        await claimJobInstance(db, job.id, instance.id);
        await markOutboxDispatched(db, row.id);
        console.log(`dispatched ${job.kind} job ${job.id} → ${instance.id}`);
      } catch (err) {
        // Leave the row undispatched; the next cron tick retries with backoff.
        console.error(`dispatch ${job.id} failed: ${err instanceof Error ? err.message : err}`);
      }
    }
  }).pipe(Effect.orDie);

// Jobs Worker: Pi runtime host + Workflows + cron outbox dispatcher (§9.3).
// The BROWSER binding lives only here — the app worker can never reach it.
export default class Jobs extends Cloudflare.Worker<Jobs>()(
  "Jobs",
  {
    // Stage-scoped physical name derived by alchemy (see apps/api worker).
    main: import.meta.url,
    compatibility: { date: "2025-10-01", flags: ["nodejs_compat"] },
  },
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);
    const aiClient = yield* Cloudflare.Workers.AI();
    const bucketClient = yield* Cloudflare.R2.ReadWriteBucket(ArchiveBucket);

    const modelVars = yield* ModelVarsConfig;
    // Workflow classes: yielding them registers the binding + export on this
    // worker and hands back the typed dispatch handle.
    const ingest = yield* IngestWorkflow;
    const enrich = yield* EnrichWorkflow;
    const index = yield* IndexWorkflow;
    const digest = yield* DigestWorkflow;
    const exp = yield* ExportWorkflow;
    const purge = yield* PurgeWorkflow;

    const bindings: Record<string, WorkflowBinding> = {
      ingest: toWorkflowBinding(ingest),
      enrich: toWorkflowBinding(enrich),
      index: toWorkflowBinding(index),
      digest: toWorkflowBinding(digest),
      export: toWorkflowBinding(exp),
      purge: toWorkflowBinding(purge),
    };

    yield* Cloudflare.Workers.cron(
      "* * * * *",
      Effect.fn(function* () {
        const db = yield* dbClient.raw;
        yield* sweepOutbox(db, bindings);
      }),
    );

    const probe =
      (env: JobsEnv) =>
      async (request: Request): Promise<Response | null> => {
        const url = new URL(request.url);

        if (url.pathname === "/health") return Response.json({ ok: true });

        if (url.pathname === "/probe") {
          const ai = env.AI as unknown as {
            run: (m: string, i: unknown) => Promise<unknown>;
            fetch?: typeof fetch;
          };

          const out: Record<string, unknown> = {};
          const model = url.searchParams.get("model") ?? "@cf/qwen/qwen3-30b-a3b-fp8";

          try {
            out.run = await ai.run(model, {
              messages: [{ role: "user", content: "say ok" }],
              max_tokens: 8,
            });
          } catch (e) {
            out.run_error = String(e);
          }

          if (url.searchParams.has("toolconv")) {
            try {
              out.toolconv = await ai.run(model, {
                messages: [
                  { role: "user", content: "what time is it" },
                  {
                    role: "assistant",
                    content: "",
                    tool_calls: [
                      {
                        id: "chatcmpl-tool-abc123",
                        type: "function",
                        function: { name: "get_time", arguments: "{}" },
                      },
                    ],
                  },
                  {
                    role: "tool",
                    tool_call_id: "chatcmpl-tool-abc123",
                    content: "2026-09-11T18:00:00Z",
                  },
                ],
                tools: [
                  {
                    type: "function",
                    function: {
                      name: "get_time",
                      description: "Return the current UTC time",
                      parameters: { type: "object", properties: {} },
                    },
                  },
                ],
                max_tokens: 64,
              });
            } catch (e) {
              out.toolconv_error = String(e);
            }
          }

          if (ai.fetch) {
            for (const [name, u] of [
              ["direct", "https://workers-binding.ai/ai/v1/chat/completions"],
              [
                "gateway",
                `https://workers-binding.ai/ai-gateway/gateways/${url.searchParams.get("gw") ?? "x"}/workers-ai/v1/chat/completions`,
              ],
            ] as const) {
              try {
                const r = await ai.fetch(u, {
                  method: "POST",
                  headers: { "content-type": "application/json" },
                  body: JSON.stringify({
                    model: "@cf/qwen/qwen3-30b-a3b-fp8",
                    messages: [{ role: "user", content: "say ok" }],
                    max_tokens: 8,
                  }),
                });

                out[name] = { status: r.status, body: (await r.text()).slice(0, 300) };
              } catch (e) {
                out[`${name}_error`] = String(e);
              }
            }
          } else {
            out.fetch = "undefined";
          }

          return Response.json(out);
        }

        if (url.pathname === "/probe-pi") {
          const { cfStreamFn, resolveCfModel } = await import("@recally/platform-cloudflare");
          const stream = cfStreamFn(env as never);

          const model = resolveCfModel(
            url.searchParams.get("model") ?? "@cf/qwen/qwen3-30b-a3b-fp8",
            env as never,
          );

          const events: string[] = [];
          let last: unknown = null;
          const withTools = url.searchParams.has("tools");

          try {
            const s = await Promise.resolve(
              stream(
                model,
                {
                  systemPrompt: withTools
                    ? "You must call get_time before answering."
                    : "You are terse.",
                  messages: [
                    {
                      role: "user",
                      content: withTools ? "what time is it" : "say ok",
                      timestamp: Date.now(),
                    },
                  ],
                  ...(withTools
                    ? {
                        tools: [
                          {
                            name: "get_time",
                            label: "Get time",
                            description: "Return the current UTC time",
                            parameters: { type: "object", properties: {} },
                          },
                        ],
                      }
                    : {}),
                } as never,
                {} as never,
              ),
            );

            for await (const ev of s) {
              events.push(ev.type);
              last = ev;
            }

            return Response.json({ events: events.slice(0, 30), last });
          } catch (e) {
            return Response.json({ events, error: String(e) });
          }
        }

        return null;
      };

    return {
      fetch: Effect.gen(function* () {
        const request = yield* HttpServerRequest.HttpServerRequest;
        const [db, r2, aiRaw] = yield* Effect.all([dbClient.raw, bucketClient.raw, aiClient.raw]);

        const env: JobsEnv = {
          DB: db,
          ARCHIVE_BUCKET: r2,
          AI: aiRaw,
          ...definedModelVars(modelVars),
        };

        const response = yield* Effect.promise(() =>
          Promise.resolve(probe(env)(request.source as unknown as Request)),
        );

        if (response) return HttpServerResponse.fromWeb(response);

        return HttpServerResponse.text("not found", { status: 404 });
      }),
    };
  }).pipe(
    Effect.provide(
      Layer.mergeAll(
        Cloudflare.D1.QueryDatabaseBinding,
        Cloudflare.R2.ReadWriteBucketBinding,
        Cloudflare.Workers.AIBinding,
        Cloudflare.Workers.CronEventSourceLive,
      ),
    ),
  ),
) {}
