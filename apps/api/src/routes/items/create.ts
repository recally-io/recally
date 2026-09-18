import { createItemSchema } from "@recally/contracts";
import { AppError, newId, normalizedUrlHash, normalizeUrl, sha256Hex } from "@recally/domain";
import {
  type CreateItemWithJobInput,
  createItemWithJob,
  findActiveItemByUrlHash,
  getIdempotencyRecord,
  putIdempotencyRecord,
} from "@recally/storage";
import { Hono } from "hono";
import { authFor } from "../../auth";
import { runRevisions } from "../../capture";
import type { Env } from "../../env";

// POST /api/v1/items (plan §12.1): idempotent, deduped, returns 202.
export const createRoutes = new Hono<{ Bindings: Env }>().post("/", async (c) => {
  const auth = authFor(c, "items:write");
  const body = createItemSchema.parse(await c.req.json());
  const idemKey = c.req.header("idempotency-key");
  const payloadHash = await sha256Hex(JSON.stringify(body));

  if (idemKey) {
    const prior = await getIdempotencyRecord(
      c.env.DB,
      auth.actor,
      "POST /items",
      auth.libraryId,
      idemKey,
    );

    if (prior) {
      if (prior.payload_hash !== payloadHash) {
        throw new AppError("conflict", "idempotency key reused with different input");
      }

      return c.newResponse(prior.response, { status: prior.status_code as 202 });
    }
  }

  let source: Pick<
    CreateItemWithJobInput,
    "originalUrl" | "normalizedUrl" | "normalizedUrlHash" | "title" | "jobPayload"
  >;

  if (body.source.kind === "content") {
    // Pasted content enters the same archive path via a manual capture.
    const normalized = `manual:${await sha256Hex(body.source.content)}`;
    source = {
      originalUrl: `manual:${newId()}`,
      normalizedUrl: normalized,
      normalizedUrlHash: await sha256Hex(normalized),
      title: body.source.title ?? null,
      jobPayload: JSON.stringify({
        kind: "manual",
        content: body.source.content,
        enrichment: body.enrichment,
        capture_intent: body.capture_intent ?? null,
      }),
    };
  } else {
    const normalized = normalizeUrl(body.source.url);
    const hash = await normalizedUrlHash(normalized);
    const existing = await findActiveItemByUrlHash(c.env.DB, auth.libraryId, hash);

    if (existing) {
      // Plain re-save returns the existing item (§2.5); re-capture is explicit.
      return c.json(
        {
          item_id: existing.id,
          job_id: null,
          status_url: `/api/v1/items/${existing.id}`,
          current_phase: "already_saved",
        },
        200,
      );
    }

    source = {
      originalUrl: body.source.url,
      normalizedUrl: normalized,
      normalizedUrlHash: hash,
      jobPayload: JSON.stringify({
        kind: "url",
        url: body.source.url,
        enrichment: body.enrichment,
        capture_intent: body.capture_intent ?? null,
        collection_ids: body.collection_ids,
      }),
    };
  }

  const { item, job } = await createItemWithJob(c.env.DB, {
    ...source,
    libraryId: auth.libraryId,
    ...(body.note !== undefined ? { note: body.note } : {}),
    jobKind: "ingest",
    runRevisions: await runRevisions(c.env),
  });

  const res = {
    item_id: item.id,
    job_id: job.id,
    status_url: `/api/v1/jobs/${job.id}`,
    current_phase: "queued",
  };

  if (idemKey) {
    await putIdempotencyRecord(c.env.DB, {
      actor: auth.actor,
      route: "POST /items",
      libraryId: auth.libraryId,
      key: idemKey,
      payloadHash,
      statusCode: 202,
      response: JSON.stringify(res),
    });
  }

  return c.json(res, 202);
});
