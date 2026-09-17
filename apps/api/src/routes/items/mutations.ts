import { patchItemSchema } from "@recally/contracts";
import { AppError, newId, nowIso } from "@recally/domain";
import { deleteItem, enqueueCaptureRun, requireItem } from "@recally/storage";
import { runRevisions } from "../../capture";
import { Hono } from "hono";
import { authFor } from "../../auth";
import type { Env } from "../../env";

export const mutationRoutes = new Hono<{ Bindings: Env }>()
  .patch("/:id", async (c) => {
    const auth = authFor(c, "items:write");
    const body = patchItemSchema.parse(await c.req.json());
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));
    const sets: string[] = [];
    const binds: unknown[] = [];

    if (body.title !== undefined) {
      sets.push("title = ?");
      binds.push(body.title);
    }

    if (body.read_status !== undefined) {
      sets.push("read_status = ?");
      binds.push(body.read_status);

      if (body.read_status === "read" && !item.first_read_at) {
        sets.push("first_read_at = ?", "last_read_at = ?");
        binds.push(nowIso(), nowIso());
      }
    }

    sets.push("updated_at = ?");
    binds.push(nowIso(), item.id, auth.libraryId);
    await c.env.DB.prepare(`UPDATE items SET ${sets.join(", ")} WHERE id = ? AND library_id = ?`)
      .bind(...(binds as never[]))
      .run();

    return c.json({ ok: true });
  })

  .delete("/:id", async (c) => {
    const auth = authFor(c, "items:write");
    const ok = await deleteItem(c.env.DB, auth.libraryId, c.req.param("id"));

    if (!ok) throw new AppError("not_found", "item not found");
    // Purge job cleans vectors/objects asynchronously; API refuses immediately.
    const now = nowIso();
    const jobId = newId();
    await c.env.DB.batch([
      c.env.DB.prepare(
        `INSERT INTO jobs (id, library_id, kind, item_id, payload, created_at, updated_at)
         VALUES (?, ?, 'purge', ?, ?, ?, ?)`,
      ).bind(jobId, auth.libraryId, c.req.param("id"), "{}", now, now),
      c.env.DB.prepare(
        `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
         VALUES (?, ?, ?, 'purge', ?, ?, ?)`,
      ).bind(newId(), auth.libraryId, jobId, "{}", now, now),
    ]);

    return c.json({ ok: true, purge_job_id: jobId });
  })

  // Explicit re-capture: bumps generation, creates a fresh CaptureRun (§2.5).
  .post("/:id/captures", async (c) => {
    const auth = authFor(c, "items:write");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));
    const generation = item.capture_generation + 1;

    const { jobId } = await enqueueCaptureRun(
      c.env.DB,
      {
        libraryId: auth.libraryId,
        itemId: item.id,
        generation,
        targetUrl: item.original_url,
        kind: "ingest",
        payload: JSON.stringify({ kind: "url", url: item.original_url, recapture: true }),
        revisions: await runRevisions(c.env),
      },
      [
        c.env.DB.prepare(
          "UPDATE items SET capture_generation = ?, updated_at = ? WHERE id = ?",
        ).bind(generation, nowIso(), item.id),
      ],
    );

    return c.json({ job_id: jobId, generation }, 202);
  })

  // Paste content for needs_input items (plan §2.3, M2-T06).
  .post("/:id/content-inputs", async (c) => {
    const auth = authFor(c, "items:write");
    const item = await requireItem(c.env.DB, auth.libraryId, c.req.param("id"));
    const body = (await c.req.json()) as { content?: string; title?: string };

    if (!body.content) throw new AppError("invalid_input", "content required");

    const { jobId } = await enqueueCaptureRun(c.env.DB, {
      libraryId: auth.libraryId,
      itemId: item.id,
      generation: item.capture_generation,
      targetUrl: item.original_url,
      kind: "ingest",
      payload: JSON.stringify({ kind: "manual", content: body.content, title: body.title ?? null }),
      revisions: await runRevisions(c.env),
    });

    return c.json({ job_id: jobId }, 202);
  });
