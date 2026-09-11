import { nowIso } from "@recally/domain";
import { dispatchJob, type WorkflowBinding } from "@recally/platform-cloudflare";
import { claimJobInstance, dueOutbox, getJob, markOutboxDispatched } from "@recally/storage";
import type { Env } from "./env";

export { EnrichWorkflow } from "./workflows/enrich";
export { IndexWorkflow } from "./workflows/index-wf";
export { IngestWorkflow } from "./workflows/ingest";
export {
  DigestWorkflow,
  ExportWorkflow,
  PurgeWorkflow,
} from "./workflows/purge";

const KIND_BINDINGS: Record<string, keyof Env> = {
  ingest: "INGEST_WORKFLOW",
  enrich: "ENRICH_WORKFLOW",
  index: "INDEX_WORKFLOW",
  digest: "DIGEST_WORKFLOW",
  export: "EXPORT_WORKFLOW",
  purge: "PURGE_WORKFLOW",
};

// Jobs Worker: Pi runtime host + Workflows + cron outbox dispatcher (§9.3).
// The BROWSER binding lives only here — the app worker can never reach it.
export default {
  async fetch(request: Request, env: Env): Promise<Response> {
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
            `https://workers-binding.ai/ai-gateway/gateways/${new URL(request.url).searchParams.get("gw") ?? "x"}/workers-ai/v1/chat/completions`,
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
    return new Response("not found", { status: 404 });
  },

  async scheduled(_controller: ScheduledController, env: Env): Promise<void> {
    // Catch-up pass: dispatch any outbox rows whose workflow create failed or
    // never ran (§9.3). Deterministic instance ids make re-dispatch safe.
    const now = nowIso();
    const due = await dueOutbox(env.DB, now, 20);
    console.log(`outbox sweep: ${due.length} due`);
    for (const row of due) {
      const job = await getJob(env.DB, row.library_id, row.job_id);
      if (!job || job.status === "cancelled") {
        await markOutboxDispatched(env.DB, row.id);
        continue;
      }
      const bindingKey = KIND_BINDINGS[job.kind];
      if (!bindingKey) continue;
      const binding = env[bindingKey] as WorkflowBinding;
      try {
        const instance = await dispatchJob(binding, job.id, job.generation, {
          jobId: job.id,
          libraryId: job.library_id,
        });
        await claimJobInstance(env.DB, job.id, instance.id);
        await markOutboxDispatched(env.DB, row.id);
        console.log(`dispatched ${job.kind} job ${job.id} → ${instance.id}`);
      } catch (err) {
        // Leave the row undispatched; the next cron tick retries with backoff.
        console.error(`dispatch ${job.id} failed: ${err instanceof Error ? err.message : err}`);
      }
    }
  },
};
