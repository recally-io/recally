import { buildFtsQuery } from "@recally/search";
import { Hono } from "hono";
import { authFor } from "../auth";
import type { Env } from "../env";

export const searchRoutes = new Hono<{ Bindings: Env }>().get("/", async (c) => {
  const auth = authFor(c, "search:read");
  const q = c.req.query("q") ?? "";
  const ftsQuery = buildFtsQuery(q);

  if (!ftsQuery) return c.json({ results: [] });

  // FTS hits always pass through D1 ownership/deletion/version checks before
  // reaching the user (plan §11.4).
  const { results } = await c.env.DB.prepare(
    `SELECT s.entity_type, s.entity_id, snippet(search_index, 3, '<b>', '</b>', '…', 32) AS snippet
       FROM search_index s
       WHERE search_index MATCH ? AND s.library_id = ?
       LIMIT 40`,
  )
    .bind(ftsQuery, auth.libraryId)
    .all<{ entity_type: string; entity_id: string; snippet: string }>();

  const resolved = [];

  for (const r of results) {
    if (r.entity_type === "chunk") {
      const row = await c.env.DB.prepare(
        `SELECT i.id AS item_id, i.title, i.original_url, c.id AS chunk_id
           FROM chunks c JOIN items i ON i.id = c.item_id
           WHERE c.id = ? AND c.library_id = ? AND i.deleted_at IS NULL
             AND i.current_snapshot_id IS NOT NULL`,
      )
        .bind(r.entity_id, auth.libraryId)
        .first();

      if (row) resolved.push({ ...row, snippet: r.snippet });
    } else if (r.entity_type === "item_title") {
      const row = await c.env.DB.prepare(
        "SELECT id AS item_id, title, original_url FROM items WHERE id = ? AND library_id = ? AND deleted_at IS NULL",
      )
        .bind(r.entity_id, auth.libraryId)
        .first();

      if (row) resolved.push({ ...row, snippet: r.snippet });
    }
  }

  return c.json({ results: resolved, mode: "fts" });
});
