import type { RunOutcome } from "@recally/agent-runtime";
import { AppError, newId, nowIso } from "@recally/domain";
import { type CaptureRunRow, getJob, updateJobStatus } from "@recally/storage";
import type { JobsEnv } from "../../env";
import type { WorkflowInput } from "../common";

export async function initializeCapture(db: JobsEnv["DB"], { jobId, libraryId }: WorkflowInput) {
  const job = await getJob(db, libraryId, jobId);

  if (!job) throw new AppError("not_found", `job ${jobId}`);

  const run = await db
    .prepare("SELECT * FROM capture_runs WHERE job_id = ? AND library_id = ?")
    .bind(jobId, libraryId)
    .first<CaptureRunRow>();

  if (!run) throw new AppError("not_found", `run for job ${jobId}`);

  const item = await db
    .prepare("SELECT deleted_at, capture_generation FROM items WHERE id = ? AND library_id = ?")
    .bind(run.item_id, libraryId)
    .first<{ deleted_at: string | null; capture_generation: number }>();

  if (!item || item.deleted_at || item.capture_generation !== run.generation) {
    await updateJobStatus(db, jobId, "failed", {
      error: "stale_generation",
    });

    return null;
  }

  const attemptId = newId();
  const now = nowIso();
  await db.batch([
    db
      .prepare(
        "INSERT INTO capture_attempts (id, library_id, run_id, attempt_no, kind, status, started_at) VALUES (?, ?, ?, 1, 'fetch', 'running', ?)",
      )
      .bind(attemptId, libraryId, run.id, now),
    db.prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?").bind(now, jobId),
    db
      .prepare("UPDATE capture_runs SET status = 'running', updated_at = ? WHERE id = ?")
      .bind(now, run.id),
  ]);

  return { run, job, attemptId };
}

export async function finalizeCapture(
  db: JobsEnv["DB"],
  jobId: string,
  runId: string,
  outcome: Exclude<RunOutcome, { kind: "continue" }>,
) {
  const now = nowIso();

  if (outcome.kind === "proposed") {
    // ArchiveService already published the snapshot.
    await db.batch([
      db
        .prepare(
          "UPDATE jobs SET status = 'succeeded', result = 'complete', updated_at = ? WHERE id = ?",
        )
        .bind(now, jobId),
      db
        .prepare(
          "UPDATE capture_runs SET status = 'succeeded', outcome_code = 'complete', updated_at = ? WHERE id = ?",
        )
        .bind(now, runId),
    ]);
  } else if (outcome.kind === "finished") {
    const jobStatus =
      outcome.outcome === "needs_input"
        ? "needs_input"
        : outcome.outcome === "complete" || outcome.outcome === "partial"
          ? "succeeded"
          : "failed";

    await db.batch([
      db
        .prepare("UPDATE jobs SET status = ?, result = ?, updated_at = ? WHERE id = ?")
        .bind(jobStatus, outcome.reasonCode, now, jobId),
      db
        .prepare(
          "UPDATE capture_runs SET status = ?, outcome_code = ?, updated_at = ? WHERE id = ?",
        )
        .bind(jobStatus === "succeeded" ? "succeeded" : jobStatus, outcome.outcome, now, runId),
    ]);
  } else {
    await db.batch([
      db
        .prepare(
          "UPDATE jobs SET status = 'failed', result = 'budget_exceeded', updated_at = ? WHERE id = ?",
        )
        .bind(now, jobId),
      db
        .prepare(
          "UPDATE capture_runs SET status = 'failed', outcome_code = 'budget_exhausted', updated_at = ? WHERE id = ?",
        )
        .bind(now, runId),
    ]);
  }
}
