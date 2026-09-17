import { listItems, requireItem } from "@recally/storage";
import { Hono } from "hono";
import { authFor } from "../../auth";
import type { Env } from "../../env";

export const readRoutes = new Hono<{ Bindings: Env }>()
  .get("/", async (c) => {
    const auth = authFor(c, "items:read");

    const { items, nextCursor } = await listItems(c.env.DB, auth.libraryId, {
      cursor: c.req.query("cursor") ?? undefined,
      limit: Number(c.req.query("limit") ?? 50),
    });

    const latestJobs = items.length
      ? await c.env.DB.prepare(
          `SELECT item_id, id, kind, status FROM jobs j1
           WHERE j1.item_id IN (${items.map(() => "?").join(",")})
             AND j1.created_at = (SELECT MAX(created_at) FROM jobs WHERE item_id = j1.item_id)`,
        )
          .bind(...items.map((i) => i.id))
          .all<{ item_id: string; id: string; kind: string; status: string }>()
      : { results: [] };

    const jobByItem = new Map(latestJobs.results.map((j) => [j.item_id, j]));

    const snapIds = items.map((i) => i.current_snapshot_id).filter(Boolean) as string[];

    const snaps = snapIds.length
      ? await c.env.DB.prepare(
          `SELECT id, content_quality, resource_quality FROM snapshots WHERE id IN (${snapIds.map(() => "?").join(",")})`,
        )
          .bind(...snapIds)
          .all<{
            id: string;
            content_quality: string;
            resource_quality: string;
          }>()
      : { results: [] };

    const snapById = new Map(snaps.results.map((s) => [s.id, s]));

    return c.json({
      items: items.map((i) => ({
        id: i.id,
        original_url: i.original_url,
        title: i.title,
        saved_at: i.saved_at,
        read_status: i.read_status,
        current_snapshot_id: i.current_snapshot_id,
        capture_generation: i.capture_generation,
        content_quality: i.current_snapshot_id
          ? (snapById.get(i.current_snapshot_id)?.content_quality ?? null)
          : null,
        resource_quality: i.current_snapshot_id
          ? (snapById.get(i.current_snapshot_id)?.resource_quality ?? null)
          : null,
        latest_job: jobByItem.get(i.id)
          ? {
              id: jobByItem.get(i.id)!.id,
              kind: jobByItem.get(i.id)!.kind,
              status: jobByItem.get(i.id)!.status,
            }
          : null,
      })),
      next_cursor: nextCursor,
    });
  })

  .get("/:id", async (c) => {
    const auth = authFor(c, "items:read");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));

    const snapshots = await c.env.DB.prepare(
      `SELECT s.id, s.captured_at, s.capture_method, s.content_quality, s.resource_quality,
              s.final_url, cr.id AS content_revision_id
       FROM snapshots s
       LEFT JOIN content_revisions cr ON cr.snapshot_id = s.id
       WHERE s.item_id = ? AND s.library_id = ? ORDER BY s.captured_at DESC`,
    )
      .bind(item.id, auth.libraryId)
      .all();

    const notes = await c.env.DB.prepare(
      "SELECT id, body, version, updated_at FROM notes WHERE item_id = ? AND library_id = ?",
    )
      .bind(item.id, auth.libraryId)
      .all();

    const artifacts = await c.env.DB.prepare(
      `SELECT id, type, status, output, created_at FROM ai_artifacts
       WHERE item_id = ? AND library_id = ? ORDER BY created_at DESC`,
    )
      .bind(item.id, auth.libraryId)
      .all();

    return c.json({
      item,
      snapshots: snapshots.results,
      notes: notes.results,
      artifacts: artifacts.results,
    });
  });
