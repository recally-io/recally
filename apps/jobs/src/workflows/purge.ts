import { WorkflowEntrypoint, type WorkflowEvent, type WorkflowStep } from "cloudflare:workers";
import { nowIso } from "@recally/domain";
import type { Env } from "../env";

interface JobParams {
  jobId: string;
  libraryId: string;
}

// PurgeWorkflow (§15.1): tombstone-first delete. M5 scope — this skeleton
// clears search index rows and marks the job; R2/Vectorize sweep lands there.
export class PurgeWorkflow extends WorkflowEntrypoint<Env, JobParams> {
  async run(event: WorkflowEvent<JobParams>, step: WorkflowStep) {
    const { jobId, libraryId } = event.payload;
    const { item_id } = await step.do("load", async () => {
      const j = await this.env.DB.prepare("SELECT item_id, payload FROM jobs WHERE id = ?")
        .bind(jobId)
        .first<{ item_id: string; payload: string }>();
      return {
        item_id: j?.item_id ?? (JSON.parse(j?.payload ?? "{}").item_id as string),
      };
    });
    await step.do("purge-fts", async () => {
      await this.env.DB.prepare(
        "DELETE FROM search_index WHERE library_id = ? AND object_id IN (SELECT id FROM items WHERE id = ?)",
      )
        .bind(libraryId, item_id)
        .run();
    });
    await step.do("done", async () => {
      await this.env.DB.prepare("UPDATE jobs SET status = 'succeeded', updated_at = ? WHERE id = ?")
        .bind(nowIso(), jobId)
        .run();
    });
  }
}

// ExportWorkflow (§15.2): versioned manifest + JSONL + assets, written to R2
// in bounded batches. M5 scope.
export class ExportWorkflow extends WorkflowEntrypoint<Env, JobParams> {
  async run(event: WorkflowEvent<JobParams>, step: WorkflowStep) {
    const { jobId } = event.payload;
    await step.do("todo", async () => {
      await this.env.DB.prepare(
        "UPDATE jobs SET status = 'failed', result = 'not_implemented', updated_at = ? WHERE id = ?",
      )
        .bind(nowIso(), jobId)
        .run();
    });
  }
}

// DigestWorkflow: recurring library回顾 (M3+). Skeleton only.
export class DigestWorkflow extends WorkflowEntrypoint<Env, JobParams> {
  async run(event: WorkflowEvent<JobParams>, step: WorkflowStep) {
    const { jobId } = event.payload;
    await step.do("todo", async () => {
      await this.env.DB.prepare(
        "UPDATE jobs SET status = 'failed', result = 'not_implemented', updated_at = ? WHERE id = ?",
      )
        .bind(nowIso(), jobId)
        .run();
    });
  }
}
