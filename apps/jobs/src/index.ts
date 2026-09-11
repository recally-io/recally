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
  async fetch(request: Request, _env: Env): Promise<Response> {
    const url = new URL(request.url);
    if (url.pathname === "/health") return Response.json({ ok: true });
    return new Response("not found", { status: 404 });
  },

  async scheduled(_controller: ScheduledController, env: Env): Promise<void> {
    // Catch-up pass: dispatch any outbox rows whose workflow create failed or
    // never ran (§9.3). Deterministic instance ids make re-dispatch safe.
    const now = nowIso();
    for (const row of await dueOutbox(env.DB, now, 20)) {
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
      } catch {
        // Leave the row undispatched; the next cron tick retries with backoff.
      }
    }
  },
};
