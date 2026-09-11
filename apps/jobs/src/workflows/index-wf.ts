import { WorkflowEntrypoint, type WorkflowEvent, type WorkflowStep } from "cloudflare:workers";
import { resolveModels, WorkersAIProvider } from "@recally/ai";
import { newId, nowIso } from "@recally/domain";
import { R2EvidenceStore, VectorIndex } from "@recally/platform-cloudflare";
import { chunkBlocks, EMBEDDING_VERSION, tokenizeForFts, vectorId } from "@recally/search";
import { r2Keys } from "@recally/storage";
import type { Env } from "../env";

interface IndexParams {
  jobId: string;
  libraryId: string;
}

// IndexWorkflow is deterministic: chunk → FTS → embed → vectorize, no agent.
// FTS and chunk rows commit in one D1 batch; Vectorize tracks
// pending → submitted → queryable separately (§11.4).
export class IndexWorkflow extends WorkflowEntrypoint<Env, IndexParams> {
  async run(event: WorkflowEvent<IndexParams>, step: WorkflowStep) {
    const { jobId } = event.payload;
    const env = this.env;

    try {
      return await this.runInner(event, step);
    } catch (err) {
      // Same guard as IngestWorkflow: an exhausted step must not leave the
      // job 'running' forever.
      await env.DB.prepare(
        "UPDATE jobs SET status = 'failed', result = 'workflow_error', error = ?, updated_at = ? WHERE id = ?",
      )
        .bind(err instanceof Error ? err.message.slice(0, 500) : String(err), nowIso(), jobId)
        .run();
      throw err;
    }
  }

  private async runInner(event: WorkflowEvent<IndexParams>, step: WorkflowStep) {
    const { jobId, libraryId } = event.payload;
    const env = this.env;

    const job = await step.do("load", async () => {
      const j = await env.DB.prepare("SELECT * FROM jobs WHERE id = ? AND library_id = ?")
        .bind(jobId, libraryId)
        .first<{ id: string; item_id: string; payload: string }>();
      if (!j) throw new Error(`job ${jobId} not found`);
      await env.DB.prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?")
        .bind(nowIso(), jobId)
        .run();
      return j;
    });
    const { content_revision_id } = JSON.parse(job.payload) as {
      content_revision_id: string;
    };
    if (!job.item_id) throw new Error("index job missing item_id");

    const blocks = await step.do("load-blocks", async () => {
      const evidence = new R2EvidenceStore(env.ARCHIVE_BUCKET);
      const jsonl = await evidence.getText(r2Keys.contentBlocks(libraryId, content_revision_id));
      if (!jsonl) throw new Error(`no blocks for ${content_revision_id}`);
      return jsonl
        .split("\n")
        .filter(Boolean)
        .map((l) => JSON.parse(l) as { id: string; kind: string; text: string });
    });

    const chunks = await step.do("chunk-fts", async () => {
      const chunked = chunkBlocks(blocks);
      const now = nowIso();
      // Idempotent re-index: a prior attempt (or the same step retried after
      // partial batch) leaves chunks behind — clear them or the ordinal
      // UNIQUE constraint fails forever.
      const stmts: D1PreparedStatement[] = [
        env.DB.prepare(
          "DELETE FROM search_index WHERE library_id = ? AND entity_type = 'chunk' AND entity_id IN (SELECT id FROM chunks WHERE content_revision_id = ? AND library_id = ?)",
        ).bind(libraryId, content_revision_id, libraryId),
        env.DB.prepare("DELETE FROM chunks WHERE content_revision_id = ? AND library_id = ?").bind(
          content_revision_id,
          libraryId,
        ),
      ];
      const rows: Array<{ id: string; text: string }> = [];
      for (let i = 0; i < chunked.length; i++) {
        const c = chunked[i]!;
        const id = newId();
        rows.push({ id, text: c.text });
        stmts.push(
          env.DB.prepare(
            `INSERT INTO chunks (id, library_id, content_revision_id, item_id, ordinal, block_range, text,
              text_hash, tokenizer_version, embedding_version, index_state, created_at)
             VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, ?, 'pending', ?)`,
          ).bind(
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
          env.DB.prepare(
            "INSERT INTO search_index (library_id, entity_type, entity_id, text) VALUES (?, 'chunk', ?, ?)",
          ).bind(libraryId, id, tokenizeForFts(c.text)),
        );
      }
      await env.DB.batch(stmts);
      return rows.map((r, i) => ({ ...r, ordinal: i }));
    });

    await step.do("embed-upsert", async () => {
      const models = resolveModels(env);
      const ai = new WorkersAIProvider(env.AI as never);
      const vectors = new VectorIndex(env.VECTOR_INDEX);
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
      await env.DB.batch(
        pairs.map((p) =>
          env.DB.prepare(
            "UPDATE chunks SET index_state = 'submitted', vector_id = ? WHERE id = ?",
          ).bind(p.vectorId, p.chunkId),
        ),
      );
      await env.DB.prepare(
        "UPDATE jobs SET status = 'succeeded', result = 'indexed', updated_at = ? WHERE id = ?",
      )
        .bind(now, jobId)
        .run();
    });
  }
}
