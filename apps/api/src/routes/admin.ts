import { createTokenSchema } from "@recally/contracts";
import { AppError, newId, nowIso, sha256Hex } from "@recally/domain";
import { Hono } from "hono";
import { authFor, randomToken } from "../auth";
import type { Env } from "../env";

export const adminRoutes = new Hono<{ Bindings: Env }>()
  .post("/tokens", async (c) => {
    const auth = authFor(c, "tokens:manage");
    const body = createTokenSchema.parse(await c.req.json());
    const token = randomToken("rcl_");
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
    const auth = authFor(c, "tokens:manage");

    const { results } = await c.env.DB.prepare(
      `SELECT id, name, scopes, last_used_at, expires_at, created_at FROM api_tokens
       WHERE library_id = ? AND revoked_at IS NULL`,
    )
      .bind(auth.libraryId)
      .all();

    return c.json({ tokens: results });
  })
  .delete("/tokens/:id", async (c) => {
    const auth = authFor(c, "tokens:manage");
    await c.env.DB.prepare("UPDATE api_tokens SET revoked_at = ? WHERE id = ? AND library_id = ?")
      .bind(nowIso(), c.req.param("id"), auth.libraryId)
      .run();

    return c.json({ ok: true });
  })
  .get("/usage", async (c) => {
    const auth = c.get("auth");

    const { results } = await c.env.DB.prepare(
      `SELECT kind, SUM(amount) AS total FROM usage_events
       WHERE library_id = ? AND created_at > ? GROUP BY kind`,
    )
      .bind(auth.libraryId, new Date(Date.now() - 30 * 86400_000).toISOString())
      .all();

    return c.json({ period_days: 30, usage: results });
  })
  .get("/settings", async (c) => {
    const auth = c.get("auth");

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
    const auth = authFor(c, "settings:write");

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
