import { AppError, newId, nowIso } from "@recally/domain";
import type { ItemRow, JobRow } from "../rows";
import { type CaptureRunRevisions, enqueueCaptureRun, getJob } from "./jobs";
import { assertWritable } from "./libraries";

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

export async function requireItem(
  db: D1Database,
  libraryId: string,
  itemId: string,
  message = "item not found",
): Promise<ItemRow> {
  const item = await getItem(db, libraryId, itemId);

  if (!item) throw new AppError("not_found", message);

  return item;
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
  runRevisions: CaptureRunRevisions;
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
    ...(noteStmt ? [noteStmt] : []),
  ];

  const { jobId } = await enqueueCaptureRun(
    db,
    {
      libraryId: input.libraryId,
      itemId,
      generation: 1,
      targetUrl: input.originalUrl,
      kind: input.jobKind,
      payload: input.jobPayload,
      revisions: input.runRevisions,
    },
    stmts,
  );

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
