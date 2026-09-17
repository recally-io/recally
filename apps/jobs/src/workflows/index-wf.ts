import { resolveModels, WorkersAIProvider } from "@recally/ai";
import { newId, nowIso } from "@recally/domain";
import { R2EvidenceStore, VectorIndex as VectorIndexClient } from "@recally/platform-cloudflare";
import { chunkBlocks, EMBEDDING_VERSION, tokenizeForFts, vectorId } from "@recally/search";
import { r2Keys } from "@recally/storage";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import { definedModelVars, ModelVarsConfig } from "../../../../infra/model-env";
import { ArchiveBucket, Database, VectorIndex } from "../../../../infra/resources";
import type { JobsEnv } from "../env";

interface IndexParams {
  jobId: string;
  libraryId: string;
}

// IndexWorkflow is deterministic: chunk → FTS → embed → vectorize, no agent.
// FTS and chunk rows commit in one D1 batch; Vectorize tracks
// pending → submitted → queryable separately (§11.4).
export class IndexWorkflow extends Cloudflare.Workflow<IndexWorkflow>()(
  "IndexWorkflow",
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);
    const bucketClient = yield* Cloudflare.R2.ReadWriteBucket(ArchiveBucket);
    const aiClient = yield* Cloudflare.Workers.AI();
    const vecClient = yield* Cloudflare.Vectorize.SearchIndex(VectorIndex);
    const modelVars = yield* ModelVarsConfig;

    return Effect.fn(function* (input: IndexParams) {
      const [db, r2, aiRaw, vec] = yield* Effect.all([
        dbClient.raw,
        bucketClient.raw,
        aiClient.raw,
        vecClient.raw,
      ]);

      const env: JobsEnv = {
        DB: db,
        ARCHIVE_BUCKET: r2,
        AI: aiRaw,
        ...definedModelVars(modelVars),
      };

      const evidence = new R2EvidenceStore(r2);
      const { jobId, libraryId } = input;

      const pipeline = Effect.gen(function* () {
        const job = yield* Cloudflare.Workflows.task(
          "load",
          Effect.tryPromise(async () => {
            const j = await db
              .prepare("SELECT * FROM jobs WHERE id = ? AND library_id = ?")
              .bind(jobId, libraryId)
              .first<{ id: string; item_id: string; payload: string }>();

            if (!j) throw new Error(`job ${jobId} not found`);
            await db
              .prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?")
              .bind(nowIso(), jobId)
              .run();

            return j;
          }).pipe(Effect.orDie),
        );

        const { content_revision_id } = JSON.parse(job.payload) as {
          content_revision_id: string;
        };

        if (!job.item_id) throw new Error("index job missing item_id");

        const blocks = yield* Cloudflare.Workflows.task(
          "load-blocks",
          Effect.tryPromise(async () => {
            const jsonl = await evidence.getText(
              r2Keys.contentBlocks(libraryId, content_revision_id),
            );

            if (!jsonl) throw new Error(`no blocks for ${content_revision_id}`);

            return jsonl
              .split("\n")
              .filter(Boolean)
              .map((l) => JSON.parse(l) as { id: string; kind: string; text: string });
          }).pipe(Effect.orDie),
        );

        const chunks = yield* Cloudflare.Workflows.task(
          "chunk-fts",
          Effect.tryPromise(async () => {
            const chunked = chunkBlocks(blocks);
            const now = nowIso();

            // Idempotent re-index: a prior attempt (or the same step retried after
            // partial batch) leaves chunks behind — clear them or the ordinal
            // UNIQUE constraint fails forever.
            const stmts: D1PreparedStatement[] = [
              db
                .prepare(
                  "DELETE FROM search_index WHERE library_id = ? AND entity_type = 'chunk' AND entity_id IN (SELECT id FROM chunks WHERE content_revision_id = ? AND library_id = ?)",
                )
                .bind(libraryId, content_revision_id, libraryId),
              db
                .prepare("DELETE FROM chunks WHERE content_revision_id = ? AND library_id = ?")
                .bind(content_revision_id, libraryId),
            ];

            const rows: Array<{ id: string; text: string }> = [];

            for (let i = 0; i < chunked.length; i++) {
              const c = chunked[i]!;
              const id = newId();
              rows.push({ id, text: c.text });
              stmts.push(
                db
                  .prepare(
                    `INSERT INTO chunks (id, library_id, content_revision_id, item_id, ordinal, block_range, text,
              text_hash, tokenizer_version, embedding_version, index_state, created_at)
             VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, ?, 'pending', ?)`,
                  )
                  .bind(
                    id,
                    libraryId,
                    content_revision_id,
                    job.item_id,
                    i,
                    JSON.stringify(c.blockRange),
                    c.text,
                    "tok-v1",
                    EMBEDDING_VERSION,
                    now,
                  ),
                db
                  .prepare(
                    "INSERT INTO search_index (library_id, entity_type, entity_id, text) VALUES (?, 'chunk', ?, ?)",
                  )
                  .bind(libraryId, id, tokenizeForFts(c.text)),
              );
            }

            await db.batch(stmts);

            return rows.map((r, i) => ({ ...r, ordinal: i }));
          }).pipe(Effect.orDie),
        );

        yield* Cloudflare.Workflows.task(
          "embed-upsert",
          Effect.tryPromise(async () => {
            const models = resolveModels(env);
            const ai = new WorkersAIProvider(env.AI as never);
            const vectors = new VectorIndexClient(vec);
            const pairs: Array<{ chunkId: string; vectorId: string }> = [];
            // Batched embed+upsert: one AI.run per ~32 chunks. Sequential per-chunk
            // calls overran the step's wall-clock on a 73-chunk revision (runtime
            // killed the isolate as hung).
            const BATCH = 32;

            for (let i = 0; i < chunks.length; i += BATCH) {
              const batch = chunks.slice(i, i + BATCH);

              const { vectors: vecs } = await ai.embed({
                model: models.embedding,
                texts: batch.map((c) => c.text),
              });

              const upserts = [];

              for (const [j, c] of batch.entries()) {
                const vid = await vectorId({
                  libraryId,
                  contentRevisionId: content_revision_id,
                  chunkId: c.id,
                  embeddingVersion: EMBEDDING_VERSION,
                });

                upserts.push({
                  id: vid,
                  values: vecs[j]!,
                  namespace: libraryId,
                  metadata: { chunk_id: c.id },
                });
                pairs.push({ chunkId: c.id, vectorId: vid });
              }

              await vectors.upsert(upserts);
            }

            const now = nowIso();
            await db.batch(
              pairs.map((p) =>
                db
                  .prepare(
                    "UPDATE chunks SET index_state = 'submitted', vector_id = ? WHERE id = ?",
                  )
                  .bind(p.vectorId, p.chunkId),
              ),
            );
            await db
              .prepare(
                "UPDATE jobs SET status = 'succeeded', result = 'indexed', updated_at = ? WHERE id = ?",
              )
              .bind(now, jobId)
              .run();
          }).pipe(Effect.orDie),
        );
      });

      // Same guard as IngestWorkflow: an exhausted step must not leave the
      // job 'running' forever.
      yield* pipeline.pipe(
        Effect.catchDefect((err) =>
          Effect.gen(function* () {
            yield* Effect.tryPromise(async () => {
              await db
                .prepare(
                  "UPDATE jobs SET status = 'failed', result = 'workflow_error', error = ?, updated_at = ? WHERE id = ?",
                )
                .bind(
                  err instanceof Error ? err.message.slice(0, 500) : String(err),
                  nowIso(),
                  jobId,
                )
                .run();
            }).pipe(Effect.ignore);

            return yield* Effect.die(err);
          }),
        ),
      );
    });
  }).pipe(
    Effect.provide(
      Layer.mergeAll(
        Cloudflare.D1.QueryDatabaseBinding,
        Cloudflare.R2.ReadWriteBucketBinding,
        Cloudflare.Workers.AIBinding,
        Cloudflare.Vectorize.SearchIndexBinding,
      ),
    ),
  ),
) {}
