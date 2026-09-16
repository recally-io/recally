import { nowIso } from "@recally/domain";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import { Database } from "../../../../infra/resources";

interface JobParams {
  jobId: string;
  libraryId: string;
}

// PurgeWorkflow (§15.1): tombstone-first delete. M5 scope — this skeleton
// clears search index rows and marks the job; R2/Vectorize sweep lands there.
export class PurgeWorkflow extends Cloudflare.Workflow<PurgeWorkflow>()(
  "PurgeWorkflow",
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);

    return Effect.fn(function* (input: JobParams) {
      const db = yield* dbClient.raw;
      const { jobId, libraryId } = input;
      const loaded = yield* Cloudflare.Workflows.task(
        "load",
        Effect.tryPromise(async () => {
          const j = await db
            .prepare("SELECT item_id, payload FROM jobs WHERE id = ?")
            .bind(jobId)
            .first<{ item_id: string; payload: string }>();
          return {
            item_id: j?.item_id ?? (JSON.parse(j?.payload ?? "{}").item_id as string),
          };
        }).pipe(Effect.orDie),
      );
      yield* Cloudflare.Workflows.task(
        "purge-fts",
        Effect.tryPromise(async () => {
          await db
            .prepare(
              "DELETE FROM search_index WHERE library_id = ? AND object_id IN (SELECT id FROM items WHERE id = ?)",
            )
            .bind(libraryId, loaded.item_id)
            .run();
        }).pipe(Effect.orDie),
      );
      yield* Cloudflare.Workflows.task(
        "done",
        Effect.tryPromise(async () => {
          await db
            .prepare("UPDATE jobs SET status = 'succeeded', updated_at = ? WHERE id = ?")
            .bind(nowIso(), jobId)
            .run();
        }).pipe(Effect.orDie),
      );
    });
  }).pipe(Effect.provide(Cloudflare.D1.QueryDatabaseBinding)),
) {}

// ExportWorkflow (§15.2): versioned manifest + JSONL + assets, written to R2
// in bounded batches. M5 scope.
export class ExportWorkflow extends Cloudflare.Workflow<ExportWorkflow>()(
  "ExportWorkflow",
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);

    return Effect.fn(function* (input: JobParams) {
      const db = yield* dbClient.raw;
      const { jobId } = input;
      yield* Cloudflare.Workflows.task(
        "todo",
        Effect.tryPromise(async () => {
          await db
            .prepare(
              "UPDATE jobs SET status = 'failed', result = 'not_implemented', updated_at = ? WHERE id = ?",
            )
            .bind(nowIso(), jobId)
            .run();
        }).pipe(Effect.orDie),
      );
    });
  }).pipe(Effect.provide(Cloudflare.D1.QueryDatabaseBinding)),
) {}

// DigestWorkflow: recurring library回顾 (M3+). Skeleton only.
export class DigestWorkflow extends Cloudflare.Workflow<DigestWorkflow>()(
  "DigestWorkflow",
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);

    return Effect.fn(function* (input: JobParams) {
      const db = yield* dbClient.raw;
      const { jobId } = input;
      yield* Cloudflare.Workflows.task(
        "todo",
        Effect.tryPromise(async () => {
          await db
            .prepare(
              "UPDATE jobs SET status = 'failed', result = 'not_implemented', updated_at = ? WHERE id = ?",
            )
            .bind(nowIso(), jobId)
            .run();
        }).pipe(Effect.orDie),
      );
    });
  }).pipe(Effect.provide(Cloudflare.D1.QueryDatabaseBinding)),
) {}
