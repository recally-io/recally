import {
  askRequestSchema,
  createNoteSchema,
  createShareSchema,
  createTokenSchema,
  patchNoteSchema,
} from "@recally/contracts";
import { AppError, newId, nowIso, sha256Hex } from "@recally/domain";
import { buildFtsQuery } from "@recally/search";
import { getItem, getJob, getSnapshot, r2Keys, requestJobCancel } from "@recally/storage";
import { Hono } from "hono";
import { type AuthContext, requireScope } from "../auth";

interface Env {
  DB: D1Database;
  ARCHIVE_BUCKET: R2Bucket;
}

export const jobsRoutes = new Hono<{ Bindings: Env }>()
  .get("/", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:read");
    const { results } = await c.env.DB.prepare(
      `SELECT j.id, j.kind, j.status, j.item_id, j.attempt_count, j.error,
              j.created_at, j.updated_at, i.title AS item_title
       FROM jobs j LEFT JOIN items i ON i.id = j.item_id
       WHERE j.library_id = ? ORDER BY j.created_at DESC LIMIT 50`,
    )
      .bind(auth.libraryId)
      .all();
    return c.json({ jobs: results });
  })
  .get("/:id", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:read");
    const job = await getJob(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!job) throw new AppError("not_found", "job not found");
    const events = await c.env.DB.prepare(
      `SELECT e.sequence, e.kind, e.reason_code, e.summary, e.created_at
       FROM agent_events e JOIN capture_runs r ON r.id = e.run_id
       WHERE r.job_id = ? ORDER BY e.sequence`,
    )
      .bind(job.id)
      .all();
    return c.json({
      id: job.id,
      kind: job.kind,
      status: job.status,
      item_id: job.item_id,
      attempt_count: job.attempt_count,
      outcome_code: null,
      error: job.error,
      created_at: job.created_at,
      updated_at: job.updated_at,
      events: events.results,
    });
  })
  .post("/:id/cancel", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const ok = await requestJobCancel(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!ok) throw new AppError("conflict", "job is not cancellable");
    return c.json({ ok: true });
  });

export const searchRoutes = new Hono<{ Bindings: Env }>().get("/", async (c) => {
  const auth = c.get("auth") as AuthContext;
  requireScope(auth, "search:read");
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

export const askRoutes = new Hono<{ Bindings: Env }>().post("/", async (c) => {
  const auth = c.get("auth") as AuthContext;
  requireScope(auth, "search:read");
  askRequestSchema.parse(await c.req.json());
  throw new AppError("not_implemented", "ask lands with M4; search is available");
});

// --- shares + public route ---

export const sharesRoutes = new Hono<{ Bindings: Env }>()
  .post("/", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "shares:write");
    const body = createShareSchema.parse(await c.req.json());
    const item = await getItem(c.env.DB, auth.libraryId, body.item_id);
    if (!item) throw new AppError("not_found", "item not found");
    const snapshot = await getSnapshot(c.env.DB, auth.libraryId, body.snapshot_id);
    if (!snapshot || snapshot.item_id !== item.id) {
      throw new AppError("invalid_input", "snapshot does not belong to item");
    }
    const token = `rsh_${crypto.randomUUID().replaceAll("-", "")}${crypto.randomUUID().replaceAll("-", "")}`;
    const tokenHash = await sha256Hex(token);
    const id = newId();
    const now = nowIso();
    const expiresAt = body.expires_in_days
      ? new Date(Date.now() + body.expires_in_days * 86400_000).toISOString()
      : null;
    await c.env.DB.prepare(
      `INSERT INTO shares (id, library_id, item_id, snapshot_id, content_revision_id, token_hash,
        include_full_text, allowed_artifacts, expires_at, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
    )
      .bind(
        id,
        auth.libraryId,
        item.id,
        body.snapshot_id,
        body.content_revision_id ?? null,
        tokenHash,
        body.include_full_text ? 1 : 0,
        JSON.stringify(body.allowed_artifact_ids),
        expiresAt,
        now,
      )
      .run();
    return c.json({ share_id: id, url: `/s/${token}`, expires_at: expiresAt }, 201);
  })
  .delete("/:id", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "shares:write");
    const r = await c.env.DB.prepare(
      "UPDATE shares SET revoked_at = ? WHERE id = ? AND library_id = ? AND revoked_at IS NULL",
    )
      .bind(nowIso(), c.req.param("id"), auth.libraryId)
      .run();
    if (!(r.meta.changes ?? 0)) throw new AppError("not_found", "share not found");
    return c.json({ ok: true });
  });

// Public share page: token re-validated per request; R2 never public (§14.6).
export const publicRoutes = new Hono<{ Bindings: Env }>().get("/:token", async (c) => {
  const tokenHash = await sha256Hex(c.req.param("token"));
  const share = await c.env.DB.prepare(
    `SELECT s.*, i.title, i.original_url, i.deleted_at FROM shares s
     JOIN items i ON i.id = s.item_id WHERE s.token_hash = ?`,
  )
    .bind(tokenHash)
    .first<{
      id: string;
      library_id: string;
      snapshot_id: string;
      content_revision_id: string | null;
      include_full_text: number;
      expires_at: string | null;
      revoked_at: string | null;
      title: string | null;
      original_url: string;
      deleted_at: string | null;
    }>();
  const now = nowIso();
  if (
    !share ||
    share.revoked_at ||
    share.deleted_at ||
    (share.expires_at && share.expires_at < now)
  ) {
    throw new AppError("not_found", "share not found");
  }
  let body = "";
  if (share.include_full_text && share.content_revision_id) {
    const obj = await c.env.ARCHIVE_BUCKET.get(
      r2Keys.contentArticle(share.library_id, share.content_revision_id),
    );
    body = (await obj?.text()) ?? "";
  }
  const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  return c.html(
    `<!doctype html><meta charset="utf-8"><title>${esc(share.title ?? "shared item")}</title>
<article style="max-width:42rem;margin:3rem auto;font-family:system-ui">
<p><a href="${esc(share.original_url)}">${esc(share.original_url)}</a></p>
<pre style="white-space:pre-wrap">${esc(body || "Full text not included in this share.")}</pre>
</article>`,
    200,
    { "Cache-Control": "private, no-store" },
  );
});

// --- content / notes / artifacts / reading events ---

export const contentRoutes = new Hono<{ Bindings: Env }>()
  .get("/content/:revisionId", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:read");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "notes:write");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "notes:write");
    const { body, version } = patchNoteSchema.parse(await c.req.json());
    const note = await c.env.DB.prepare("SELECT * FROM notes WHERE id = ? AND library_id = ?")
      .bind(c.req.param("id"), auth.libraryId)
      .first<{ id: string; item_id: string; version: number }>();
    if (!note) throw new AppError("not_found", "note not found");
    const item = await getItem(c.env.DB, auth.libraryId, note.item_id);
    if (!item) throw new AppError("not_found", "note not found");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
    const body = (await c.req.json()) as { kind?: string; meta?: unknown };
    const kind = body.kind === "finish" || body.kind === "progress" ? body.kind : "open";
    await c.env.DB.prepare(
      "INSERT INTO reading_events (id, library_id, item_id, kind, occurred_at, meta) VALUES (?, ?, ?, ?, ?, ?)",
    )
      .bind(newId(), auth.libraryId, item.id, kind, nowIso(), JSON.stringify(body.meta ?? {}))
      .run();
    return c.json({ ok: true }, 201);
  });

// --- admin: tokens / usage / settings / collections / tags ---

export const adminRoutes = new Hono<{ Bindings: Env }>()
  .post("/tokens", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "tokens:manage");
    const body = createTokenSchema.parse(await c.req.json());
    const token = `rcl_${crypto.randomUUID().replaceAll("-", "")}${crypto.randomUUID().replaceAll("-", "")}`;
    const id = newId();
    await c.env.DB.prepare(
      `INSERT INTO api_tokens (id, library_id, name, token_hash, scopes, expires_at, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?)`,
    )
      .bind(
        id,
        auth.libraryId,
        body.name,
        await sha256Hex(token),
        JSON.stringify(body.scopes),
        body.expires_in_days
          ? new Date(Date.now() + body.expires_in_days * 86400_000).toISOString()
          : null,
        nowIso(),
      )
      .run();
    return c.json({ token_id: id, token }, 201); // plaintext returned once only
  })
  .get("/tokens", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "tokens:manage");
    const { results } = await c.env.DB.prepare(
      `SELECT id, name, scopes, last_used_at, expires_at, created_at FROM api_tokens
       WHERE library_id = ? AND revoked_at IS NULL`,
    )
      .bind(auth.libraryId)
      .all();
    return c.json({ tokens: results });
  })
  .delete("/tokens/:id", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "tokens:manage");
    await c.env.DB.prepare("UPDATE api_tokens SET revoked_at = ? WHERE id = ? AND library_id = ?")
      .bind(nowIso(), c.req.param("id"), auth.libraryId)
      .run();
    return c.json({ ok: true });
  })
  .get("/usage", async (c) => {
    const auth = c.get("auth") as AuthContext;
    const { results } = await c.env.DB.prepare(
      `SELECT kind, SUM(amount) AS total FROM usage_events
       WHERE library_id = ? AND created_at > ? GROUP BY kind`,
    )
      .bind(auth.libraryId, new Date(Date.now() - 30 * 86400_000).toISOString())
      .all();
    return c.json({ period_days: 30, usage: results });
  })
  .get("/settings", async (c) => {
    const auth = c.get("auth") as AuthContext;
    const lib = await c.env.DB.prepare("SELECT settings, timezone FROM libraries WHERE id = ?")
      .bind(auth.libraryId)
      .first<{ settings: string; timezone: string }>();
    if (!lib) throw new AppError("not_found", "library not found");
    return c.json({
      timezone: lib.timezone,
      settings: JSON.parse(lib.settings),
    });
  })
  .patch("/settings", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "settings:write");
    const body = (await c.req.json()) as {
      timezone?: string;
      settings?: Record<string, unknown>;
    };
    const sets: string[] = [];
    const binds: unknown[] = [];
    if (body.timezone) {
      sets.push("timezone = ?");
      binds.push(body.timezone);
    }
    if (body.settings) {
      sets.push("settings = ?");
      binds.push(JSON.stringify(body.settings));
    }
    if (!sets.length) throw new AppError("invalid_input", "empty patch");
    sets.push("updated_at = ?");
    binds.push(nowIso(), auth.libraryId);
    await c.env.DB.prepare(`UPDATE libraries SET ${sets.join(", ")} WHERE id = ?`)
      .bind(...(binds as never[]))
      .run();
    return c.json({ ok: true });
  });
