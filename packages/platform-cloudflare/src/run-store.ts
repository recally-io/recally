import type { SourceRef, ToolContext } from "@recally/capture";
import { newId, nowIso } from "@recally/domain";
import { r2Keys } from "@recally/storage";
import type { RunStore, StoredSource } from "@recally/tools";
import type { R2EvidenceStore } from "./evidence";

// Run-scoped evidence persistence over D1 + R2 (§6.6, §8.1). Sources land in
// R2 under runs/{run}/attempts/{attempt}/sources/ and are indexed in
// source_documents; observations likewise. GC consults committed manifests —
// never delete by path alone.
export class D1R2RunStore implements RunStore {
  constructor(
    private readonly db: D1Database,
    private readonly evidence: R2EvidenceStore,
  ) {}

  async saveSource(
    ctx: ToolContext,
    input: {
      url: string;
      kind: SourceRef["kind"];
      contentType: string;
      body: Uint8Array | string;
    },
  ): Promise<SourceRef> {
    const sourceId = newId();
    const file =
      input.kind === "screenshot"
        ? "screenshot.webp"
        : input.kind === "rendered_dom"
          ? "rendered.html"
          : input.kind === "adapter_record"
            ? "record.json"
            : "response-body.html";
    const key = r2Keys.source(ctx.libraryId, ctx.runId, ctx.attemptId, sourceId, file);
    const put = await this.evidence.put(key, input.body, input.contentType);
    await this.db
      .prepare(
        `INSERT INTO source_documents
         (id, library_id, run_id, attempt_id, url, kind, content_type, r2_key, sha256, byte_size, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
      )
      .bind(
        sourceId,
        ctx.libraryId,
        ctx.runId,
        ctx.attemptId,
        input.url,
        input.kind,
        input.contentType,
        key,
        put.sha256,
        put.byteSize,
        nowIso(),
      )
      .run();
    return {
      sourceId,
      url: input.url,
      kind: input.kind,
      sha256: put.sha256,
      byteSize: put.byteSize,
      evidenceRef: key,
    };
  }

  async getSource(ctx: ToolContext, sourceId: string): Promise<StoredSource | null> {
    const row = await this.db
      .prepare(
        `SELECT id, url, kind, r2_key, sha256, byte_size FROM source_documents
         WHERE id = ? AND library_id = ? AND run_id = ?`,
      )
      .bind(sourceId, ctx.libraryId, ctx.runId)
      .first<{
        id: string;
        url: string;
        kind: SourceRef["kind"];
        r2_key: string;
        sha256: string;
        byte_size: number;
      }>();
    if (!row) return null;
    return {
      ref: {
        sourceId: row.id,
        url: row.url,
        kind: row.kind,
        sha256: row.sha256,
        byteSize: row.byte_size,
        evidenceRef: row.r2_key,
      },
      bodyKey: row.r2_key,
      blocksKey: null,
    };
  }

  async saveObservation(ctx: ToolContext, observation: unknown): Promise<string> {
    const obsId = newId();
    const key = r2Keys.observation(ctx.libraryId, ctx.runId, ctx.attemptId, obsId);
    await this.evidence.put(key, JSON.stringify(observation), "application/json");
    // Next sequence for this attempt: count existing events. Sequence conflicts
    // surface via the UNIQUE(run_id, attempt_id, sequence) constraint.
    const seq = await this.db
      .prepare(
        "SELECT COALESCE(MAX(sequence), 0) + 1 AS seq FROM agent_events WHERE run_id = ? AND attempt_id = ?",
      )
      .bind(ctx.runId, ctx.attemptId)
      .first<{ seq: number }>();
    await this.db
      .prepare(
        `INSERT INTO agent_events (id, library_id, run_id, attempt_id, sequence, kind, evidence_ref, created_at)
         VALUES (?, ?, ?, ?, ?, 'observation', ?, ?)`,
      )
      .bind(newId(), ctx.libraryId, ctx.runId, ctx.attemptId, seq?.seq ?? 1, key, nowIso())
      .run();
    return key;
  }
}
