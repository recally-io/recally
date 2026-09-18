import { nowIso } from "@recally/domain";
import { getJob, updateJobStatus } from "@recally/storage";
import * as Effect from "effect/Effect";
import type { JobsEnv } from "../env";

export interface WorkflowInput {
  jobId: string;
  libraryId: string;
}

export async function loadRunningJob(db: JobsEnv["DB"], libraryId: string, jobId: string) {
  const job = await getJob(db, libraryId, jobId);

  if (!job) throw new Error(`job ${jobId} not found`);

  await updateJobStatus(db, jobId, "running");

  return job;
}

export const markWorkflowError = (
  db: JobsEnv["DB"],
  jobId: string,
  runId: string | null,
  err: Error | string,
) =>
  Effect.tryPromise(async () => {
    const now = nowIso();

    const stmts = [
      db
        .prepare(
          "UPDATE jobs SET status = 'failed', result = 'workflow_error', error = ?, updated_at = ? WHERE id = ?",
        )
        .bind(err instanceof Error ? err.message.slice(0, 500) : err, now, jobId),
    ];

    if (runId) {
      stmts.push(
        db
          .prepare(
            "UPDATE capture_runs SET status = 'failed', outcome_code = 'failed', updated_at = ? WHERE id = ?",
          )
          .bind(now, runId),
      );
    }

    await db.batch(stmts);
  }).pipe(Effect.ignore);

export const markNotImplemented = (db: JobsEnv["DB"], jobId: string) =>
  updateJobStatus(db, jobId, "failed", { result: "not_implemented" });
