import { createNoteSchema, patchNoteSchema } from "@recally/contracts";
import { AppError, newId, nowIso } from "@recally/domain";
import { requireItem } from "@recally/storage";
import { Hono } from "hono";
import { authFor } from "../auth";
import type { Env } from "../env";

export const contentRoutes = new Hono<{ Bindings: Env }>()
  .get("/content/:revisionId", async (c) => {
    const auth = authFor(c, "items:read");

    const rev = await c.env.DB.prepare(
      "SELECT * FROM content_revisions WHERE id = ? AND library_id = ?",
    )
      .bind(c.req.param("revisionId"), auth.libraryId)
      .first<{ article_key: string; blocks_key: string }>();

    if (!rev) throw new AppError("not_found", "revision not found");
    const article = await c.env.ARCHIVE_BUCKET.get(rev.article_key);
    const blocks = await c.env.ARCHIVE_BUCKET.get(rev.blocks_key);

    return c.json({
      article_md: (await article?.text()) ?? null,
      blocks: blocks
        ? (await blocks.text())
            .split("\n")
            .filter(Boolean)
            .map((l) => JSON.parse(l))
        : [],
    });
  })
  .post("/items/:id/notes", async (c) => {
    const auth = authFor(c, "notes:write");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));
    const { body } = createNoteSchema.parse(await c.req.json());
    const id = newId();
    const now = nowIso();
    await c.env.DB.batch([
      c.env.DB.prepare(
        "INSERT INTO notes (id, library_id, item_id, body, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
      ).bind(id, auth.libraryId, item.id, body, now, now),
      c.env.DB.prepare(
        "INSERT INTO note_revisions (id, note_id, version, body, created_at) VALUES (?, ?, 1, ?, ?)",
      ).bind(newId(), id, body, now),
    ]);

    return c.json({ note_id: id, version: 1 }, 201);
  })
  .patch("/notes/:id", async (c) => {
    const auth = authFor(c, "notes:write");
    const { body, version } = patchNoteSchema.parse(await c.req.json());

    const note = await c.env.DB.prepare("SELECT * FROM notes WHERE id = ? AND library_id = ?")
      .bind(c.req.param("id"), auth.libraryId)
      .first<{ id: string; item_id: string; version: number }>();

    if (!note) throw new AppError("not_found", "note not found");
    await requireItem(c.env.DB, auth.libraryId, note.item_id, "note not found");

    if (note.version !== version) {
      throw new AppError("conflict", "note version mismatch", {
        nextAction: "reload",
      });
    }

    const now = nowIso();
    await c.env.DB.batch([
      c.env.DB.prepare("UPDATE notes SET body = ?, version = ?, updated_at = ? WHERE id = ?").bind(
        body,
        version + 1,
        now,
        note.id,
      ),
      c.env.DB.prepare(
        "INSERT INTO note_revisions (id, note_id, version, body, created_at) VALUES (?, ?, ?, ?, ?)",
      ).bind(newId(), note.id, version + 1, body, now),
    ]);

    return c.json({ ok: true, version: version + 1 });
  })
  .post("/items/:id/artifacts/import", async (c) => {
    const auth = authFor(c, "items:write");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));

    const body = (await c.req.json()) as {
      content?: string;
      source_label?: string;
    };

    if (!body.content) throw new AppError("invalid_input", "content required");
    const id = newId();
    await c.env.DB.prepare(
      `INSERT INTO ai_artifacts (id, library_id, item_id, type, status, output, created_at)
       VALUES (?, ?, ?, 'imported_answer', 'unverified_import', ?, ?)`,
    )
      .bind(
        id,
        auth.libraryId,
        item.id,
        JSON.stringify({
          content: body.content,
          source_label: body.source_label ?? null,
        }),
        nowIso(),
      )
      .run();

    return c.json({ artifact_id: id, status: "unverified_import" }, 201);
  })
  .post("/items/:id/reading-events", async (c) => {
    const auth = authFor(c, "items:write");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));
    const body = (await c.req.json()) as { kind?: string; meta?: unknown };
    const kind = body.kind === "finish" || body.kind === "progress" ? body.kind : "open";
    await c.env.DB.prepare(
      "INSERT INTO reading_events (id, library_id, item_id, kind, occurred_at, meta) VALUES (?, ?, ?, ?, ?, ?)",
    )
      .bind(newId(), auth.libraryId, item.id, kind, nowIso(), JSON.stringify(body.meta ?? {}))
      .run();

    return c.json({ ok: true }, 201);
  });
