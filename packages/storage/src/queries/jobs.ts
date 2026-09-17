import { newId, nowIso } from "@recally/domain";
import type { JobRow, OutboxRow } from "../rows";

export interface CaptureRunRevisions {
  agentRuntimeVersion: string;
  skillRevision: string;
  toolsetVersion: string;
  policyVersion: string;
  pipelineVersion: string;
  modelConfigVersion: string;
}

// Preceding item/url/note writes share this batch, so a dispatch failure cannot
// lose a save. Every run pins revisions at creation (§9.3, §9.6).
export async function enqueueCaptureRun(
  db: D1Database,
  input: {
    libraryId: string;
    itemId: string;
    generation: number;
    targetUrl: string;
    kind: string;
    payload: string;
    revisions: CaptureRunRevisions;
  },
  precedingStatements: D1PreparedStatement[] = [],
): Promise<{ jobId: string; runId: string }> {
  const now = nowIso();
  const jobId = newId();
  const runId = newId();
  const { libraryId, itemId, generation, targetUrl, kind, payload, revisions } = input;

  await db.batch([
    ...precedingStatements,
    // Job before run: capture_runs.job_id references jobs(id), and D1 enforces
    // foreign keys per statement inside a batch.
    db
      .prepare(
        `INSERT INTO jobs (id, library_id, kind, item_id, generation, payload, created_at, updated_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(jobId, libraryId, kind, itemId, generation, payload, now, now),
    db
      .prepare(
        `INSERT INTO capture_runs
       (id, library_id, item_id, job_id, generation, agent_runtime_version, skill_revision,
        toolset_version, policy_version, pipeline_version, model_config_version,
        requested_at, target_url, created_at, updated_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(
        runId,
        libraryId,
        itemId,
        jobId,
        generation,
        revisions.agentRuntimeVersion,
        revisions.skillRevision,
        revisions.toolsetVersion,
        revisions.policyVersion,
        revisions.pipelineVersion,
        revisions.modelConfigVersion,
        now,
        targetUrl,
        now,
        now,
      ),
    db
      .prepare(
        `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(newId(), libraryId, jobId, kind, payload, now, now),
  ]);

  return { jobId, runId };
}

export async function getJob(
  db: D1Database,
  libraryId: string,
  jobId: string,
): Promise<JobRow | null> {
  return db
    .prepare("SELECT * FROM jobs WHERE id = ? AND library_id = ?")
    .bind(jobId, libraryId)
    .first<JobRow>();
}

export async function updateJobStatus(
  db: D1Database,
  jobId: string,
  status: string,
  extra: { error?: string; result?: string; nextAttemptAt?: string } = {},
): Promise<void> {
  await db
    .prepare(
      `UPDATE jobs SET status = ?, error = COALESCE(?, error), result = COALESCE(?, result),
       next_attempt_at = COALESCE(?, next_attempt_at), updated_at = ?
       WHERE id = ?`,
    )
    .bind(
      status,
      extra.error ?? null,
      extra.result ?? null,
      extra.nextAttemptAt ?? null,
      nowIso(),
      jobId,
    )
    .run();
}

export async function requestJobCancel(
  db: D1Database,
  libraryId: string,
  jobId: string,
): Promise<boolean> {
  const r = await db
    .prepare(
      `UPDATE jobs SET cancel_requested = 1, updated_at = ? WHERE id = ? AND library_id = ?
       AND status IN ('queued', 'running', 'retry_wait', 'needs_input')`,
    )
    .bind(nowIso(), jobId, libraryId)
    .run();

  return (r.meta.changes ?? 0) > 0;
}

export async function dueOutbox(db: D1Database, now: string, limit: number): Promise<OutboxRow[]> {
  const { results } = await db
    .prepare(
      `SELECT o.* FROM outbox o
       WHERE o.dispatched_at IS NULL AND o.available_at <= ?
       ORDER BY o.available_at, o.id LIMIT ?`,
    )
    .bind(now, limit)
    .all<OutboxRow>();

  return results;
}

export async function markOutboxDispatched(db: D1Database, outboxId: string): Promise<void> {
  await db
    .prepare("UPDATE outbox SET dispatched_at = ? WHERE id = ?")
    .bind(nowIso(), outboxId)
    .run();
}

export async function claimJobInstance(
  db: D1Database,
  jobId: string,
  instanceId: string,
): Promise<void> {
  await db
    .prepare("UPDATE jobs SET workflow_instance_id = ?, updated_at = ? WHERE id = ?")
    .bind(instanceId, nowIso(), jobId)
    .run();
}
