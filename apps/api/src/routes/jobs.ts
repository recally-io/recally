import { AppError } from "@recally/domain";
import { getJob, requestJobCancel } from "@recally/storage";
import { Hono } from "hono";
import { authFor } from "../auth";
import type { Env } from "../env";

export const jobsRoutes = new Hono<{ Bindings: Env }>()
  .get("/", async (c) => {
    const auth = authFor(c, "items:read");

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
    const auth = authFor(c, "items:read");
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
    const auth = authFor(c, "items:write");
    const ok = await requestJobCancel(c.env.DB, auth.libraryId, c.req.param("id"));

    if (!ok) throw new AppError("conflict", "job is not cancellable");

    return c.json({ ok: true });
  });
