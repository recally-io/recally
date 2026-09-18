import { newId, nowIso } from "@recally/domain";

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
