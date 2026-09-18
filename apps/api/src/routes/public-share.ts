import { AppError, nowIso, sha256Hex } from "@recally/domain";
import { r2Keys } from "@recally/storage";
import { Hono } from "hono";
import type { Env } from "../env";

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
