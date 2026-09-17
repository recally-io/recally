import type { SnapshotRow } from "../rows";

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
