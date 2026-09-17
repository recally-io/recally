import { createShareSchema } from "@recally/contracts";
import { AppError, newId, nowIso, sha256Hex } from "@recally/domain";
import { getSnapshot, requireItem } from "@recally/storage";
import { Hono } from "hono";
import { authFor, randomToken } from "../auth";
import type { Env } from "../env";

export const sharesRoutes = new Hono<{ Bindings: Env }>()
  .post("/", async (c) => {
    const auth = authFor(c, "shares:write");
    const body = createShareSchema.parse(await c.req.json());
    const item = await requireItem(c.env.DB, auth.libraryId, body.item_id);
    const snapshot = await getSnapshot(c.env.DB, auth.libraryId, body.snapshot_id);

    if (!snapshot || snapshot.item_id !== item.id) {
      throw new AppError("invalid_input", "snapshot does not belong to item");
    }

    const token = randomToken("rsh_");
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
    const auth = authFor(c, "shares:write");

    const r = await c.env.DB.prepare(
      "UPDATE shares SET revoked_at = ? WHERE id = ? AND library_id = ? AND revoked_at IS NULL",
    )
      .bind(nowIso(), c.req.param("id"), auth.libraryId)
      .run();

    if (!(r.meta.changes ?? 0)) throw new AppError("not_found", "share not found");

    return c.json({ ok: true });
  });
