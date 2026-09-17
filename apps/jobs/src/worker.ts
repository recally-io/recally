import { nowIso } from "@recally/domain";
import { dispatchJob, type WorkflowBinding } from "@recally/platform-cloudflare";
import { claimJobInstance, dueOutbox, getJob, markOutboxDispatched } from "@recally/storage";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import * as HttpServerRequest from "effect/unstable/http/HttpServerRequest";
import * as HttpServerResponse from "effect/unstable/http/HttpServerResponse";
import { Database } from "../../../infra/resources";
import type { JobsEnv } from "./env";
import { JobsBaseLayer } from "./bindings";
import { EnrichWorkflow } from "./workflows/enrich";
import { IndexWorkflow } from "./workflows/index-wf";
import { IngestWorkflow } from "./workflows/ingest";
import { DigestWorkflow, ExportWorkflow, PurgeWorkflow } from "./workflows/purge";

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
          status: () => Effect.runPromise(instance.status()).then((s) => ({ status: s.status })),
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

      const binding = bindings[job.kind];

      if (!binding) continue;

      try {
        const instance = await dispatchJob(binding, job.id, job.generation, {
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
    // Workflow classes: yielding them registers the binding + export on this
    // worker and hands back the typed dispatch handle.
    const ingest = yield* IngestWorkflow;
    const enrich = yield* EnrichWorkflow;
    const index = yield* IndexWorkflow;
    const digest = yield* DigestWorkflow;
    const exp = yield* ExportWorkflow;
    const purge = yield* PurgeWorkflow;

    const bindings = {
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

    return {
      fetch: Effect.gen(function* () {
        const request = yield* HttpServerRequest.HttpServerRequest;

        if (new URL(request.url, "https://jobs.internal").pathname === "/health") {
          return HttpServerResponse.jsonUnsafe({ ok: true });
        }

        return HttpServerResponse.text("not found", { status: 404 });
      }),
    };
  }).pipe(Effect.provide(Layer.mergeAll(JobsBaseLayer, Cloudflare.Workers.CronEventSourceLive))),
) {}
