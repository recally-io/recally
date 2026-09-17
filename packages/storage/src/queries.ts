import { AppError, newId, nowIso } from "@recally/domain";
import type { ItemRow, JobRow, LibraryRow, OutboxRow, SnapshotRow } from "./rows";

// Single write layer: every business write goes through helpers here so the
// maintenance barrier and deletion tombstones are checked in exactly one place
// (plan §15.4 — a hidden button is not a barrier).

export async function assertWritable(db: D1Database, libraryId: string): Promise<void> {
  const row = await db
    .prepare("SELECT mode FROM library_maintenance WHERE library_id = ?")
    .bind(libraryId)
    .first<{ mode: string }>();

  if (row?.mode === "write_barrier") {
    throw new AppError("write_barrier", "library is in backup write barrier", {
      retryable: true,
      nextAction: "retry_later",
    });
  }
}

// --- Libraries ---

export async function getLibrary(db: D1Database, id: string): Promise<LibraryRow | null> {
  return db.prepare("SELECT * FROM libraries WHERE id = ?").bind(id).first<LibraryRow>();
}

export async function getLibraryByAccessSub(
  db: D1Database,
  sub: string,
): Promise<LibraryRow | null> {
  return db.prepare("SELECT * FROM libraries WHERE access_sub = ?").bind(sub).first<LibraryRow>();
}

export async function createLibrary(
  db: D1Database,
  input: { id?: string; name: string; timezone?: string; accessSub?: string },
): Promise<LibraryRow> {
  const id = input.id ?? newId();
  const now = nowIso();
  await db
    .prepare(
      "INSERT INTO libraries (id, name, timezone, access_sub, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
    )
    .bind(id, input.name, input.timezone ?? "UTC", input.accessSub ?? null, now, now)
    .run();
  const row = await getLibrary(db, id);

  if (!row) throw new AppError("internal", "library insert failed");

  return row;
}

// --- Items ---

export async function getItem(
  db: D1Database,
  libraryId: string,
  itemId: string,
): Promise<ItemRow | null> {
  return db
    .prepare("SELECT * FROM items WHERE id = ? AND library_id = ? AND deleted_at IS NULL")
    .bind(itemId, libraryId)
    .first<ItemRow>();
}

export async function findActiveItemByUrlHash(
  db: D1Database,
  libraryId: string,
  normalizedUrlHash: string,
): Promise<ItemRow | null> {
  return db
    .prepare(
      `SELECT i.* FROM item_urls u JOIN items i ON i.id = u.item_id
       WHERE u.library_id = ? AND u.normalized_url_hash = ? AND u.is_active = 1
         AND i.deleted_at IS NULL`,
    )
    .bind(libraryId, normalizedUrlHash)
    .first<ItemRow>();
}

// Keyset pagination on (saved_at, id); cursor is opaque base64 of "saved_at|id".
export async function listItems(
  db: D1Database,
  libraryId: string,
  opts: { cursor?: string | undefined; limit?: number | undefined } = {},
): Promise<{ items: ItemRow[]; nextCursor: string | null }> {
  const limit = Math.min(opts.limit ?? 50, 100);
  let sql = "SELECT * FROM items WHERE library_id = ? AND deleted_at IS NULL";
  const binds: unknown[] = [libraryId];

  if (opts.cursor) {
    const [savedAt, id] = atob(opts.cursor).split("|");
    sql += " AND (saved_at < ? OR (saved_at = ? AND id < ?))";
    binds.push(savedAt, savedAt, id);
  }

  sql += " ORDER BY saved_at DESC, id DESC LIMIT ?";
  binds.push(limit + 1);

  const { results } = await db
    .prepare(sql)
    .bind(...(binds as never[]))
    .all<ItemRow>();

  const hasMore = results.length > limit;
  const items = hasMore ? results.slice(0, limit) : results;
  const last = items[items.length - 1];
  const nextCursor = hasMore && last ? btoa(`${last.saved_at}|${last.id}`) : null;

  return { items, nextCursor };
}

export interface CreateItemWithJobInput {
  libraryId: string;
  originalUrl: string;
  normalizedUrl: string;
  normalizedUrlHash: string;
  title?: string | null;
  note?: string;
  jobKind: string;
  jobPayload: string;
  // Revision fencing pinned at run creation (§9.6).
  runRevisions?: {
    agentRuntimeVersion: string;
    skillRevision: string;
    toolsetVersion: string;
    policyVersion: string;
    pipelineVersion: string;
    modelConfigVersion: string;
  };
}

// Item + url mapping + capture run + job + outbox in one D1 batch so a
// dispatch failure can never lose the save (plan §9.3).
export async function createItemWithJob(
  db: D1Database,
  input: CreateItemWithJobInput,
): Promise<{ item: ItemRow; job: JobRow }> {
  await assertWritable(db, input.libraryId);
  const now = nowIso();
  const itemId = newId();
  const jobId = newId();
  const runId = newId();
  const outboxId = newId();

  const runStmt = input.runRevisions
    ? db
        .prepare(
          `INSERT INTO capture_runs
           (id, library_id, item_id, job_id, generation, agent_runtime_version, skill_revision,
            toolset_version, policy_version, pipeline_version, model_config_version,
            requested_at, target_url, created_at, updated_at)
           VALUES (?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        )
        .bind(
          runId,
          input.libraryId,
          itemId,
          jobId,
          input.runRevisions.agentRuntimeVersion,
          input.runRevisions.skillRevision,
          input.runRevisions.toolsetVersion,
          input.runRevisions.policyVersion,
          input.runRevisions.pipelineVersion,
          input.runRevisions.modelConfigVersion,
          now,
          input.originalUrl,
          now,
          now,
        )
    : null;

  const noteStmt = input.note
    ? db
        .prepare(
          "INSERT INTO notes (id, library_id, item_id, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
        )
        .bind(newId(), input.libraryId, itemId, input.note, now, now)
    : null;

  const stmts = [
    db
      .prepare(
        `INSERT INTO items (id, library_id, original_url, normalized_url, title, saved_at, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(
        itemId,
        input.libraryId,
        input.originalUrl,
        input.normalizedUrl,
        input.title ?? null,
        now,
        now,
        now,
      ),
    db
      .prepare(
        `INSERT INTO item_urls (id, library_id, item_id, normalized_url_hash, normalized_url, created_at)
         VALUES (?, ?, ?, ?, ?, ?)`,
      )
      .bind(newId(), input.libraryId, itemId, input.normalizedUrlHash, input.normalizedUrl, now),
    db
      .prepare(
        `INSERT INTO jobs (id, library_id, kind, item_id, payload, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(jobId, input.libraryId, input.jobKind, itemId, input.jobPayload, now, now),
    db
      .prepare(
        `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(outboxId, input.libraryId, jobId, input.jobKind, input.jobPayload, now, now),
    ...(runStmt ? [runStmt] : []),
    ...(noteStmt ? [noteStmt] : []),
  ];

  await db.batch(stmts);

  const item = await getItem(db, input.libraryId, itemId);
  const job = await getJob(db, input.libraryId, jobId);

  if (!item || !job) throw new AppError("internal", "item creation failed");

  return { item, job };
}

// Soft delete: tombstone + deactivate url mappings + revoke shares. Purge is a
// separate workflow; reads refuse immediately (plan §15.1).
export async function deleteItem(
  db: D1Database,
  libraryId: string,
  itemId: string,
): Promise<boolean> {
  await assertWritable(db, libraryId);
  const now = nowIso();

  const result = await db.batch([
    db
      .prepare(
        "UPDATE items SET deleted_at = ?, capture_generation = capture_generation + 1, updated_at = ? WHERE id = ? AND library_id = ? AND deleted_at IS NULL",
      )
      .bind(now, now, itemId, libraryId),
    db
      .prepare("UPDATE item_urls SET is_active = 0 WHERE item_id = ? AND library_id = ?")
      .bind(itemId, libraryId),
    db
      .prepare(
        "UPDATE shares SET revoked_at = ? WHERE item_id = ? AND library_id = ? AND revoked_at IS NULL",
      )
      .bind(now, itemId, libraryId),
    db
      .prepare(
        "INSERT INTO deletion_tombstones (id, library_id, object_type, object_id, deleted_at) VALUES (?, ?, 'item', ?, ?)",
      )
      .bind(newId(), libraryId, itemId, now),
  ]);

  return (result[0]?.meta.changes ?? 0) > 0;
}

// --- Jobs ---

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

// --- Outbox ---

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

// --- Idempotency ---

export interface IdemRecord {
  id: string;
  status_code: number | null;
  response: string | null;
  payload_hash: string;
}

export async function getIdempotencyRecord(
  db: D1Database,
  actor: string,
  route: string,
  libraryId: string,
  key: string,
): Promise<IdemRecord | null> {
  return db
    .prepare(
      `SELECT id, status_code, response, payload_hash FROM idempotency_records
       WHERE actor = ? AND route = ? AND library_id = ? AND idem_key = ?`,
    )
    .bind(actor, route, libraryId, key)
    .first<IdemRecord>();
}

export async function putIdempotencyRecord(
  db: D1Database,
  input: {
    actor: string;
    route: string;
    libraryId: string;
    key: string;
    payloadHash: string;
    statusCode: number;
    response: string;
  },
): Promise<void> {
  await db
    .prepare(
      `INSERT OR IGNORE INTO idempotency_records
       (id, actor, library_id, route, idem_key, payload_hash, status_code, response, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    )
    .bind(
      newId(),
      input.actor,
      input.libraryId,
      input.route,
      input.key,
      input.payloadHash,
      input.statusCode,
      input.response,
      nowIso(),
    )
    .run();
}

// --- Snapshots ---

export async function getSnapshot(
  db: D1Database,
  libraryId: string,
  snapshotId: string,
): Promise<SnapshotRow | null> {
  return db
    .prepare("SELECT * FROM snapshots WHERE id = ? AND library_id = ?")
    .bind(snapshotId, libraryId)
    .first<SnapshotRow>();
}
