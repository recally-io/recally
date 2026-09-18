import type { Verification } from "@recally/ai";
import type { ContentBlock, ToolContext } from "@recally/capture";
import type { ArchiveProposal } from "@recally/contracts";
import { EXTRACTOR_VERSION, newId, nowIso, PIPELINE_VERSION } from "@recally/domain";
import { renderAdapterRecord } from "@recally/site-adapters";
import { r2Keys } from "@recally/storage";
import type { CommitResult, Committer } from "@recally/tools";
import { parseBlocks } from "@recally/tools";
import { selectBlockRange, stringifyBlocksJsonl } from "@recally/tools/content-blocks";
import type { R2EvidenceStore } from "./evidence";
import type { D1R2RunStore } from "./run-store";

// ArchiveService (§5.10, §8.2): the only path from agent proposal to published
// snapshot. Deterministic checks here; optional independent model verifier via
// the `verify` hook — it judges content evidence, not the agent's self-report.

interface CommitterDeps {
  db: D1Database;
  evidence: R2EvidenceStore;
  runStore: D1R2RunStore;
  verify?: (input: {
    title: string | null;
    headText: string;
    tailText: string;
    missingParts: string[];
    claimedQuality: "complete" | "partial";
  }) => Promise<Verification>;
}

export class ArchiveService implements Committer {
  constructor(private readonly deps: CommitterDeps) {}

  async commit(ctx: ToolContext, proposal: ArchiveProposal): Promise<CommitResult> {
    const { db, evidence, runStore } = this.deps;

    // Fencing: item must be live and this run must still be the current
    // generation before anything publishes (§9.7).
    const run = await db
      .prepare("SELECT * FROM capture_runs WHERE id = ? AND library_id = ?")
      .bind(ctx.runId, ctx.libraryId)
      .first<{ generation: number; item_id: string; target_url: string }>();

    const item = run
      ? await db
          .prepare(
            "SELECT capture_generation, deleted_at, current_snapshot_id FROM items WHERE id = ? AND library_id = ?",
          )
          .bind(run.item_id, ctx.libraryId)
          .first<{
            capture_generation: number;
            deleted_at: string | null;
            current_snapshot_id: string | null;
          }>()
      : null;

    if (!run || !item) return { status: "rejected", problems: ["run or item not found"] };

    if (item.deleted_at) return { status: "rejected", problems: ["item deleted"] };

    if (run.generation !== item.capture_generation || run.generation !== ctx.generation) {
      return { status: "rejected", problems: ["stale_generation"] };
    }

    // Validate sources + block ranges against stored evidence.
    const problems: string[] = [];
    const selectedBlocks: ContentBlock[] = [];
    const knownSources = (await runStore.listSources(ctx)).map((s) => s.sourceId);
    let finalUrl = "";

    for (const sourceId of proposal.sourceIds) {
      const stored = await runStore.getSource(ctx, sourceId);

      if (!stored) {
        problems.push(
          `unknown source ${sourceId}; sources saved this run: ${knownSources.join(", ") || "none"}`,
        );
        continue;
      }

      if (!finalUrl) finalUrl = stored.ref.url;
      const ranges = proposal.selectedBlockRanges.filter((r) => r.sourceId === sourceId);
      const body = await evidence.getText(stored.bodyKey);

      if (body === null) {
        problems.push(`missing source body ${sourceId}`);
        continue;
      }

      if (stored.ref.kind === "adapter_record") {
        // Structured records get a deterministic readable rendering when one
        // exists (§7.1); the JSON body remains the evidence either way.
        const rendered = renderAdapterRecord(body);

        if (rendered?.length) {
          selectedBlocks.push(
            ...rendered.map((b, i) => ({ id: `${sourceId}:r${i}`, kind: b.kind, text: b.text })),
          );
        } else {
          selectedBlocks.push({ id: `${sourceId}:record`, kind: "paragraph", text: body });
        }

        continue;
      }

      if (stored.ref.kind === "manual") {
        // User-submitted content: whole body is one block.
        selectedBlocks.push({
          id: `${sourceId}:record`,
          kind: "paragraph",
          text: body,
        });
        continue;
      }

      const { blocks } = await parseBlocks(body);

      for (const range of ranges) {
        const selected = selectBlockRange(blocks, range.startBlockId, range.endBlockId);

        if (!selected) {
          problems.push(
            `invalid block range ${range.startBlockId}..${range.endBlockId} on ${sourceId} (valid: ${blocks[0]?.id}..${blocks[blocks.length - 1]?.id}, ${blocks.length} blocks)`,
          );
          continue;
        }

        selectedBlocks.push(...selected);
      }
    }

    if (selectedBlocks.length === 0) {
      problems.push(
        "no content selected — pass block ranges for fetched/rendered sources; adapter_record and manual sources archive whole-body and need no ranges",
      );
    }

    if (problems.length) return { status: "rejected", problems };

    // Optional independent model check on real content evidence (§7.3). A
    // verifier *verdict* can reject; a verifier infra failure degrades to
    // partial — unverifiable is not the same as verified (one retry first,
    // schema drift on this model has been observed on dev).
    let contentQuality: "complete" | "partial" | "unavailable" = proposal.claimedQuality;

    if (this.deps.verify) {
      const head = selectedBlocks
        .slice(0, 8)
        .map((b) => b.text)
        .join("\n")
        .slice(0, 3000);

      const tail = selectedBlocks
        .slice(-8)
        .map((b) => b.text)
        .join("\n")
        .slice(-3000);

      const verifyInput = {
        title: null,
        headText: head,
        tailText: tail,
        missingParts: proposal.missingParts,
        claimedQuality: proposal.claimedQuality,
      };

      let verdict: Verification | null = null;
      let verifyError: string | null = null;

      for (let attempt = 0; attempt < 2 && !verdict; attempt++) {
        try {
          verdict = await this.deps.verify(verifyInput);
        } catch (err) {
          verifyError = err instanceof Error ? err.message : String(err);
        }
      }

      if (verdict?.verdict === "wrong") {
        return {
          status: "rejected",
          problems: ["verifier: content does not match the page", ...verdict.problems],
        };
      }

      if (verifyError && !verdict) {
        console.warn(`verifier failed twice for run ${ctx.runId}: ${verifyError.slice(0, 400)}`);
        proposal.missingParts.push("verifier unavailable — quality unverified");
        contentQuality = "partial";
      } else if (verdict?.verdict === "partial" || proposal.missingParts.length > 0) {
        contentQuality = "partial";
      }
    } else if (proposal.missingParts.length > 0) {
      contentQuality = "partial";
    }

    // Render deterministically from selected source blocks — model text never
    // becomes the archived original (§7.1).
    const articleMd = selectedBlocks
      .map((b) =>
        b.kind === "heading"
          ? `## ${b.text}`
          : b.kind === "code"
            ? `\`\`\`\n${b.text}\n\`\`\``
            : b.kind === "quote"
              ? `> ${b.text}`
              : b.text,
      )
      .join("\n\n");

    const blocksJsonl = stringifyBlocksJsonl(selectedBlocks);
    const now = nowIso();
    const snapshotId = newId();
    const revisionId = newId();

    const articleKey = r2Keys.contentArticle(ctx.libraryId, revisionId);
    const blocksKey = r2Keys.contentBlocks(ctx.libraryId, revisionId);
    const provenanceKey = r2Keys.contentProvenance(ctx.libraryId, revisionId);
    const manifestKey = r2Keys.snapshotManifest(ctx.libraryId, snapshotId);

    const [artPut, blkPut] = await Promise.all([
      evidence.put(articleKey, articleMd, "text/markdown"),
      evidence.put(blocksKey, blocksJsonl, "application/jsonl"),
    ]);

    await evidence.put(
      provenanceKey,
      JSON.stringify({
        sourceIds: proposal.sourceIds,
        ranges: proposal.selectedBlockRanges,
        extractorVersion: EXTRACTOR_VERSION,
      }),
      "application/json",
    );

    const manifest = {
      snapshot_id: snapshotId,
      item_id: run.item_id,
      run_id: ctx.runId,
      captured_at: now,
      sources: proposal.sourceIds,
      content_revision_id: revisionId,
      objects: [
        { key: articleKey, sha256: artPut.sha256, bytes: artPut.byteSize },
        { key: blocksKey, sha256: blkPut.sha256, bytes: blkPut.byteSize },
        { key: provenanceKey },
      ],
      missing_parts: proposal.missingParts,
      content_quality: contentQuality,
      pipeline_version: PIPELINE_VERSION,
    };

    const manifestPut = await evidence.put(
      manifestKey,
      JSON.stringify(manifest, null, 2),
      "application/json",
    );

    // Publish in one D1 batch; a worse new snapshot never demotes a better
    // current one (§8.3).
    const current = item.current_snapshot_id
      ? await db
          .prepare("SELECT content_quality FROM snapshots WHERE id = ?")
          .bind(item.current_snapshot_id)
          .first<{ content_quality: string }>()
      : null;

    const keepCurrent = current?.content_quality === "complete" && contentQuality === "partial";

    const enrichJobId = newId();
    const indexJobId = newId();
    await db.batch([
      db
        .prepare(
          `INSERT INTO snapshots (id, library_id, item_id, capture_run_id, captured_at, final_url,
          capture_method, source_manifest_key, manifest_hash, content_quality, resource_quality, committed_at, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'unknown', ?, ?)`,
        )
        .bind(
          snapshotId,
          ctx.libraryId,
          run.item_id,
          ctx.runId,
          now,
          finalUrl.startsWith("browser:") ? run.target_url : finalUrl || run.target_url,
          "fetch",
          manifestKey,
          manifestPut.sha256,
          contentQuality,
          now,
          now,
        ),
      db
        .prepare(
          `INSERT INTO content_revisions (id, library_id, snapshot_id, item_id, extractor_version,
          selected_source_blocks, body_hash, article_key, blocks_key, token_count, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        )
        .bind(
          revisionId,
          ctx.libraryId,
          snapshotId,
          run.item_id,
          EXTRACTOR_VERSION,
          JSON.stringify(proposal.selectedBlockRanges),
          artPut.sha256,
          articleKey,
          blocksKey,
          Math.ceil(articleMd.length / 4),
          now,
        ),
      db
        .prepare("UPDATE source_documents SET snapshot_id = ? WHERE run_id = ? AND library_id = ?")
        .bind(snapshotId, ctx.runId, ctx.libraryId),
      ...(keepCurrent
        ? []
        : [
            db
              .prepare("UPDATE items SET current_snapshot_id = ?, updated_at = ? WHERE id = ?")
              .bind(snapshotId, now, run.item_id),
          ]),
      // Enrich + index fan-out, registered in the same transaction (§8.2).
      ...[enrichJobId, indexJobId].flatMap((jobId, i) => [
        db
          .prepare(
            `INSERT INTO jobs (id, library_id, kind, item_id, payload, created_at, updated_at)
           VALUES (?, ?, ?, ?, ?, ?, ?)`,
          )
          .bind(
            jobId,
            ctx.libraryId,
            i === 0 ? "enrich" : "index",
            run.item_id,
            JSON.stringify({
              snapshot_id: snapshotId,
              content_revision_id: revisionId,
            }),
            now,
            now,
          ),
        db
          .prepare(
            `INSERT INTO outbox (id, library_id, job_id, kind, payload, available_at, created_at)
           VALUES (?, ?, ?, ?, ?, ?, ?)`,
          )
          .bind(
            newId(),
            ctx.libraryId,
            jobId,
            i === 0 ? "enrich" : "index",
            JSON.stringify({
              snapshot_id: snapshotId,
              content_revision_id: revisionId,
            }),
            now,
            now,
          ),
      ]),
      db
        .prepare(
          "UPDATE capture_runs SET status = 'succeeded', outcome_code = ?, updated_at = ? WHERE id = ?",
        )
        .bind(contentQuality === "complete" ? "complete" : "partial", now, ctx.runId),
    ]);

    return { status: "accepted", snapshotId, contentRevisionId: revisionId };
  }
}
