import type { ToolSpec } from "@recally/agent-runtime";
import { summaryOutputSchema } from "@recally/ai";
import { AppError } from "@recally/domain";
import { z } from "zod";
import type { ToolDeps } from "../deps";

// AnalysisAgent toolset (§10.1): read-only over one fixed snapshot. No web,
// no browser, no site tools — a summary can never smuggle live content into
// the archive.

const readOutlineInput = z.object({ content_revision_id: z.string() });

const readBlocksInput = z.object({
  content_revision_id: z.string(),
  start: z.number().int().nonnegative().default(0),
  count: z.number().int().positive().max(200).default(60),
});

const proposeSummaryInput = z.object({
  content_revision_id: z.string(),
  language: z.string(),
  output: summaryOutputSchema,
  coverage: z.object({
    total_blocks: z.number().int().nonnegative(),
    processed_blocks: z.number().int().nonnegative(),
    omitted_ranges: z.array(z.object({ start: z.number(), end: z.number() })).default([]),
  }),
});

interface BlocksDoc {
  blocks: Array<{ id: string; kind: string; text: string }>;
}

async function loadBlocks(deps: ToolDeps, ctx: unknown, revisionId: string): Promise<BlocksDoc> {
  // blocks.jsonl layout: one {id, kind, text} per line (§4.4).
  if (!deps.readRevisionBlocks) throw new AppError("internal", "revision reader not configured");
  const jsonl = await deps.readRevisionBlocks(ctx as never, revisionId);

  if (jsonl === null) throw new AppError("not_found", `no blocks for revision ${revisionId}`);

  return {
    blocks: jsonl
      .split("\n")
      .filter(Boolean)
      .map((l) => JSON.parse(l)),
  };
}

export function analysisTools(deps: ToolDeps): ToolSpec[] {
  return [
    {
      name: "read_snapshot_outline",
      label: "Read outline",
      description: "Block index (id, kind, length) and metadata for a content revision.",
      schema: readOutlineInput,
      async execute(_ctx, args) {
        const { content_revision_id } = readOutlineInput.parse(args);
        const doc = await loadBlocks(deps, _ctx, content_revision_id);

        const index = doc.blocks.map((b) => ({
          id: b.id,
          kind: b.kind,
          chars: b.text.length,
        }));

        return {
          content: `${doc.blocks.length} blocks, ${index.reduce((s, b) => s + b.chars, 0)} chars`,
          details: { index },
        };
      },
    },
    {
      name: "read_blocks",
      label: "Read blocks",
      description: "Read a bounded range of source blocks by ordinal.",
      schema: readBlocksInput,
      async execute(_ctx, args) {
        const { content_revision_id, start, count } = readBlocksInput.parse(args);
        const doc = await loadBlocks(deps, _ctx, content_revision_id);
        const slice = doc.blocks.slice(start, start + count);

        const text = slice
          .map((b) => `[${b.id}] ${b.text}`)
          .join("\n")
          .slice(0, 20_000);

        return {
          content: text || "(empty range)",
          details: { start, count: slice.length, total: doc.blocks.length },
        };
      },
    },
    {
      name: "inspect_metadata",
      label: "Inspect metadata",
      description: "Metadata candidates and quality flags for the revision.",
      schema: readOutlineInput,
      async execute(ctx, args) {
        const { content_revision_id } = readOutlineInput.parse(args);
        const stored = await deps.runStore.getSource(ctx, content_revision_id);

        return {
          content: stored
            ? `revision ${content_revision_id}`
            : `revision ${content_revision_id} (no metadata)`,
          details: stored ?? {},
        };
      },
    },
    {
      name: "propose_summary",
      label: "Propose summary",
      description: "Submit the structured summary artifact with coverage and evidence block ids.",
      schema: proposeSummaryInput,
      async execute(ctx, args) {
        const input = proposeSummaryInput.parse(args);

        if (deps.persistArtifact) {
          await deps.persistArtifact(ctx, {
            type: "summary",
            contentRevisionId: input.content_revision_id,
            output: input.output,
            coverage: input.coverage,
            modelId: ctx.toolsetVersion, // resolved model id lands in usage events
            promptVersion: ctx.skillRevision,
          });
        }

        return {
          content: "summary submitted",
          terminate: true,
          runOutcome: {
            kind: "finished",
            outcome: "complete",
            reasonCode: "summary_ready",
          },
          details: input,
        };
      },
    },
    {
      name: "finish",
      label: "Finish",
      description: "End the analysis run.",
      schema: z.object({ reason_code: z.string().max(120) }),
      async execute(_ctx, args) {
        const { reason_code } = args as { reason_code: string };

        return {
          content: `analysis finished: ${reason_code}`,
          terminate: true,
          runOutcome: {
            kind: "finished",
            outcome: "failed",
            reasonCode: reason_code,
          },
        };
      },
    },
  ];
}
