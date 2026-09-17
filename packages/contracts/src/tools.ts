import { z } from "zod";

// Tool input schemas (plan §6.2–6.4). Schemas are enforced by the executor;
// unknown fields and illegal refs are rejected before any side effect.

export const webFetchInput = z.object({
  url_ref: z.string(), // observed url or link id — resolved against policy, never arbitrary
});

export const extractContentInput = z.object({
  source_id: z.string(),
  extractor: z.string(), // extractor id; deterministic, no network
});

export const readSourceInput = z.object({
  source_id: z.string(),
  block_ids: z.array(z.string()).max(200).optional(),
  start_block: z.string().optional(),
  end_block: z.string().optional(),
});

export const siteListInput = z.object({ url_ref: z.string() });

export const siteRunInput = z.object({
  adapter_id: z.string(),
  url: z.string(), // the page url this adapter should process
  input: z.record(z.string(), z.unknown()).default({}),
});

export const browserOpenInput = z.object({ url_ref: z.string() });

export const browserObserveInput = z.object({ scope: z.string().optional() });

export const browserActInput = z.object({
  action: z.enum(["wait_for", "scroll", "expand", "navigate_same_article"]),
  target: z.string().optional(), // observed node/link ref only
});

export const browserCaptureInput = z.object({
  mode: z.enum(["rendered_dom", "screenshot", "text"]),
});

export const browserCloseInput = z.object({});

export const archiveAssetInput = z.object({ asset_ref: z.string() });

export const archiveProposalInput = z.object({
  sourceIds: z.array(z.string()).min(1),
  selectedBlockRanges: z.array(
    z.object({
      sourceId: z.string(),
      startBlockId: z.string(),
      endBlockId: z.string(),
    }),
  ),
  requiredAssetIds: z.array(z.string()).default([]),
  optionalAssetIds: z.array(z.string()).default([]),
  metadataEvidence: z.record(z.string(), z.array(z.string())).default({}),
  missingParts: z.array(z.string()).default([]),
  claimedQuality: z.enum(["complete", "partial"]),
});

export type ArchiveProposal = z.infer<typeof archiveProposalInput>;

export const finishInput = z.object({
  outcome: z.enum(["complete", "partial", "needs_input", "failed", "budget_exhausted"]),
  reason_code: z.string().max(120),
  evidence_refs: z.array(z.string()).default([]),
});

export type FinishInput = z.infer<typeof finishInput>;
