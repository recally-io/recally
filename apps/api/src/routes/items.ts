import { resolveModels } from "@recally/ai";
import { createItemSchema, patchItemSchema } from "@recally/contracts";
import {
  AppError,
  modelConfigVersion,
  newId,
  normalizedUrlHash,
  normalizeUrl,
  nowIso,
  PIPELINE_VERSION,
  POLICY_VERSION,
  sha256Hex,
} from "@recally/domain";
import { loadSkill } from "@recally/skills";
import {
  createItemWithJob,
  deleteItem,
  findActiveItemByUrlHash,
  getItem,
  listItems,
} from "@recally/storage";
import { Hono } from "hono";
import { type AuthContext, requireScope } from "../auth";

interface Env {
  DB: D1Database;
  CAPTURE_MODEL?: string;
}

const TOOLSET_VERSION = "toolset-0001";

// Revision fencing (§9.6): pin the runtime/skill/toolset/model config at run
// creation so a retry replays exactly what it started with.
async function runRevisions(env: Env) {
  const skill = await loadSkill("capture");
  const mcfg = await modelConfigVersion({ ...resolveModels(env) });
  return {
    agentRuntimeVersion: "pi-0.85.1",
    skillRevision: skill.revision,
    toolsetVersion: TOOLSET_VERSION,
    policyVersion: POLICY_VERSION,
    pipelineVersion: PIPELINE_VERSION,
    modelConfigVersion: mcfg,
  };
}

const json = (data: unknown) => JSON.parse(JSON.stringify(data));

export const itemsRoutes = new Hono<{ Bindings: Env }>()

  // POST /api/v1/items (plan §12.1): idempotent, deduped, returns 202.
  .post("/", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const body = createItemSchema.parse(await c.req.json());
    const idemKey = c.req.header("idempotency-key");
    const payloadHash = await sha256Hex(JSON.stringify(body));

    if (idemKey) {
      const prior = await c.env.DB.prepare(
        `SELECT status_code, response, payload_hash FROM idempotency_records
         WHERE actor = ? AND route = ? AND library_id = ? AND idem_key = ?`,
      )
        .bind(auth.actor, "POST /items", auth.libraryId, idemKey)
        .first<{
          status_code: number;
          response: string;
          payload_hash: string;
        }>();
      if (prior) {
        if (prior.payload_hash !== payloadHash) {
          throw new AppError("conflict", "idempotency key reused with different input");
        }
        return c.newResponse(prior.response, {
          status: prior.status_code as 202,
        });
      }
    }

    if (body.source.kind === "content") {
      // Pasted content enters the same archive path via a manual capture.
      const now = nowIso();
      const { item, job } = await createItemWithJob(c.env.DB, {
        libraryId: auth.libraryId,
        originalUrl: `manual:${newId()}`,
        normalizedUrl: `manual:${await sha256Hex(body.source.content)}`,
        normalizedUrlHash: await sha256Hex(`manual:${await sha256Hex(body.source.content)}`),
        title: body.source.title ?? null,
        ...(body.note !== undefined ? { note: body.note } : {}),
        jobKind: "ingest",
        runRevisions: await runRevisions(c.env),
        jobPayload: JSON.stringify({
          kind: "manual",
          content: body.source.content,
          enrichment: body.enrichment,
          capture_intent: body.capture_intent ?? null,
        }),
      });
      void now;
      const res = json({
        item_id: item.id,
        job_id: job.id,
        status_url: `/api/v1/jobs/${job.id}`,
        current_phase: "queued",
      });
      if (idemKey) {
        await c.env.DB.prepare(
          `INSERT OR IGNORE INTO idempotency_records
           (id, actor, library_id, route, idem_key, payload_hash, status_code, response, created_at)
           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        )
          .bind(
            newId(),
            auth.actor,
            auth.libraryId,
            "POST /items",
            idemKey,
            payloadHash,
            202,
            JSON.stringify(res),
            nowIso(),
          )
          .run();
      }
      return c.json(res, 202);
    }

    const normalized = normalizeUrl(body.source.url);
    const hash = await normalizedUrlHash(normalized);

    const existing = await findActiveItemByUrlHash(c.env.DB, auth.libraryId, hash);
    if (existing) {
      // Plain re-save returns the existing item (plan §2.5); re-capture is an
      // explicit POST /items/{id}/captures.
      const res = json({
        item_id: existing.id,
        job_id: null,
        status_url: `/api/v1/items/${existing.id}`,
        current_phase: "already_saved",
      });
      return c.json(res, 200);
    }

    const { item, job } = await createItemWithJob(c.env.DB, {
      libraryId: auth.libraryId,
      originalUrl: body.source.url,
      normalizedUrl: normalized,
      normalizedUrlHash: hash,
      ...(body.note !== undefined ? { note: body.note } : {}),
      jobKind: "ingest",
      runRevisions: await runRevisions(c.env),
      jobPayload: JSON.stringify({
        kind: "url",
        url: body.source.url,
        enrichment: body.enrichment,
        capture_intent: body.capture_intent ?? null,
        collection_ids: body.collection_ids,
      }),
    });

    const res = json({
      item_id: item.id,
      job_id: job.id,
      status_url: `/api/v1/jobs/${job.id}`,
      current_phase: "queued",
    });
    if (idemKey) {
      await c.env.DB.prepare(
        `INSERT OR IGNORE INTO idempotency_records
         (id, actor, library_id, route, idem_key, payload_hash, status_code, response, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      )
        .bind(
          newId(),
          auth.actor,
          auth.libraryId,
          "POST /items",
          idemKey,
          payloadHash,
          202,
          JSON.stringify(res),
          nowIso(),
        )
        .run();
    }
    return c.json(res, 202);
  })

  .get("/", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:read");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:read");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
    const snapshots = await c.env.DB.prepare(
      `SELECT id, captured_at, capture_method, content_quality, resource_quality, final_url
       FROM snapshots WHERE item_id = ? AND library_id = ? ORDER BY captured_at DESC`,
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
  })

  .patch("/:id", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const body = patchItemSchema.parse(await c.req.json());
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
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
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
    const now = nowIso();
    const jobId = newId();
    const runId = newId();
    const generation = item.capture_generation + 1;
    const rev = await runRevisions(c.env);
    await c.env.DB.batch([
      c.env.DB.prepare("UPDATE items SET capture_generation = ?, updated_at = ? WHERE id = ?").bind(
        generation,
        now,
        item.id,
      ),
      // Job before run: capture_runs.job_id references jobs(id) and D1
      // enforces FKs per statement inside a batch.
      c.env.DB.prepare(
        `INSERT INTO jobs (id, library_id, kind, item_id, generation, payload, created_at, updated_at)
         VALUES (?, ?, 'ingest', ?, ?, ?, ?, ?)`,
      ).bind(
        jobId,
        auth.libraryId,
        item.id,
        generation,
        JSON.stringify({
          kind: "url",
          url: item.original_url,
          recapture: true,
        }),
        now,
        now,
      ),
      c.env.DB.prepare(
        `INSERT INTO capture_runs
         (id, library_id, item_id, job_id, generation, agent_runtime_version, skill_revision,
          toolset_version, policy_version, pipeline_version, model_config_version,
          requested_at, target_url, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      ).bind(
        runId,
        auth.libraryId,
        item.id,
        jobId,
        generation,
        rev.agentRuntimeVersion,
        rev.skillRevision,
        rev.toolsetVersion,
        rev.policyVersion,
        rev.pipelineVersion,
        rev.modelConfigVersion,
        now,
        item.original_url,
        now,
        now,
      ),
      c.env.DB.prepare(
        `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
         VALUES (?, ?, ?, 'ingest', ?, ?, ?)`,
      ).bind(
        newId(),
        auth.libraryId,
        jobId,
        JSON.stringify({
          kind: "url",
          url: item.original_url,
          recapture: true,
        }),
        now,
        now,
      ),
    ]);
    return c.json({ job_id: jobId, generation }, 202);
  })

  // Paste content for needs_input items (plan §2.3, M2-T06).
  .post("/:id/content-inputs", async (c) => {
    const auth = c.get("auth") as AuthContext;
    requireScope(auth, "items:write");
    const item = await getItem(c.env.DB, auth.libraryId, c.req.param("id"));
    if (!item) throw new AppError("not_found", "item not found");
    const body = (await c.req.json()) as { content?: string; title?: string };
    if (!body.content) throw new AppError("invalid_input", "content required");
    const now = nowIso();
    const jobId = newId();
    const runId = newId();
    const rev = await runRevisions(c.env);
    await c.env.DB.batch([
      c.env.DB.prepare(
        `INSERT INTO jobs (id, library_id, kind, item_id, payload, created_at, updated_at)
         VALUES (?, ?, 'ingest', ?, ?, ?, ?)`,
      ).bind(
        jobId,
        auth.libraryId,
        item.id,
        JSON.stringify({
          kind: "manual",
          content: body.content,
          title: body.title ?? null,
        }),
        now,
        now,
      ),
      c.env.DB.prepare(
        `INSERT INTO capture_runs
         (id, library_id, item_id, job_id, generation, agent_runtime_version, skill_revision,
          toolset_version, policy_version, pipeline_version, model_config_version,
          requested_at, target_url, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      ).bind(
        runId,
        auth.libraryId,
        item.id,
        jobId,
        item.capture_generation,
        rev.agentRuntimeVersion,
        rev.skillRevision,
        rev.toolsetVersion,
        rev.policyVersion,
        rev.pipelineVersion,
        rev.modelConfigVersion,
        now,
        item.original_url,
        now,
        now,
      ),
      c.env.DB.prepare(
        `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
         VALUES (?, ?, ?, 'ingest', ?, ?, ?)`,
      ).bind(
        newId(),
        auth.libraryId,
        jobId,
        JSON.stringify({
          kind: "manual",
          content: body.content,
          title: body.title ?? null,
        }),
        now,
        now,
      ),
    ]);
    return c.json({ job_id: jobId }, 202);
  });
