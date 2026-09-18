import { AppError } from "@recally/domain";
import type { LibraryRow } from "../rows";

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

export async function getLibrary(db: D1Database, id: string): Promise<LibraryRow | null> {
  return db.prepare("SELECT * FROM libraries WHERE id = ?").bind(id).first<LibraryRow>();
}
