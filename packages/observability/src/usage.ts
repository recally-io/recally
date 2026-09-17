import { newId, nowIso } from "@recally/domain";

// Budget accounting lives in D1, not worker memory (plan §16.2). Reserve
// before the external call, settle after; a timed-out call with unknown billing
// stays "unknown", never optimistically refunded.

export type UsageKind =
  | "model_input_tokens"
  | "model_output_tokens"
  | "browser_ms"
  | "fetch_bytes"
  | "r2_bytes";

export async function reserve(
  db: D1Database,
  input: {
    libraryId: string;
    operationKey: string;
    kind: UsageKind;
    amount: number;
  },
): Promise<string> {
  const invocationId = newId();
  await db
    .prepare(
      `INSERT INTO usage_reservations (id, library_id, operation_key, invocation_id, kind, amount, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?)`,
    )
    .bind(
      newId(),
      input.libraryId,
      input.operationKey,
      invocationId,
      input.kind,
      input.amount,
      nowIso(),
    )
    .run();

  return invocationId;
}

export async function settle(
  db: D1Database,
  invocationId: string,
  actual: number,
  refs: { libraryId: string; jobId?: string; runId?: string },
): Promise<void> {
  const now = nowIso();
  await db.batch([
    db
      .prepare(
        "UPDATE usage_reservations SET state = 'settled', amount = ?, settled_at = ? WHERE invocation_id = ?",
      )
      .bind(actual, now, invocationId),
    db
      .prepare(
        `INSERT INTO usage_events (id, library_id, invocation_id, kind, amount, job_id, run_id, created_at)
         SELECT ?, r.library_id, r.invocation_id, r.kind, ?, ?, ?, ? FROM usage_reservations r
         WHERE r.invocation_id = ?`,
      )
      .bind(newId(), actual, refs.jobId ?? null, refs.runId ?? null, now, invocationId),
  ]);
}

export async function releaseUnknown(db: D1Database, invocationId: string): Promise<void> {
  await db
    .prepare(
      "UPDATE usage_reservations SET state = 'unknown', settled_at = ? WHERE invocation_id = ?",
    )
    .bind(nowIso(), invocationId)
    .run();
}
