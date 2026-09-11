import {
  ARTIFACT_STATUSES,
  ARTIFACT_TYPES,
  CONTENT_QUALITIES,
  JOB_KINDS,
  JOB_STATUSES,
  READ_STATUSES,
  RESOURCE_QUALITIES,
} from "@recally/domain";
import { z } from "zod";

// --- Item submission (plan §12.1) ---

export const itemSourceSchema = z.discriminatedUnion("kind", [
  z.object({ kind: z.literal("url"), url: z.url() }),
  z.object({
    kind: z.literal("content"),
    title: z.string().optional(),
    content: z.string().min(1),
  }),
]);

export const createItemSchema = z.object({
  source: itemSourceSchema,
  capture_intent: z.string().max(2000).optional(),
  collection_ids: z.array(z.string()).default([]),
  note: z.string().max(10_000).optional(),
  enrichment: z.enum(["auto", "defer"]).default("auto"),
});
export type CreateItemInput = z.infer<typeof createItemSchema>;

export const createItemResponseSchema = z.object({
  item_id: z.string(),
  job_id: z.string(),
  status_url: z.string(),
  current_phase: z.string(),
});
export type CreateItemResponse = z.infer<typeof createItemResponseSchema>;

// --- Views ---

export const itemViewSchema = z.object({
  id: z.string(),
  original_url: z.string(),
  title: z.string().nullable(),
  saved_at: z.string(),
  read_status: z.enum(READ_STATUSES),
  current_snapshot_id: z.string().nullable(),
  capture_generation: z.number(),
  content_quality: z.enum(CONTENT_QUALITIES).nullable(),
  resource_quality: z.enum(RESOURCE_QUALITIES).nullable(),
  latest_job: z
    .object({
      id: z.string(),
      kind: z.enum(JOB_KINDS),
      status: z.enum(JOB_STATUSES),
    })
    .nullable(),
});
export type ItemView = z.infer<typeof itemViewSchema>;

export const itemListSchema = z.object({
  items: z.array(itemViewSchema),
  next_cursor: z.string().nullable(),
});
export type ItemList = z.infer<typeof itemListSchema>;

export const jobViewSchema = z.object({
  id: z.string(),
  kind: z.enum(JOB_KINDS),
  status: z.enum(JOB_STATUSES),
  item_id: z.string().nullable(),
  attempt_count: z.number(),
  outcome_code: z.string().nullable(),
  error: z.string().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});
export type JobView = z.infer<typeof jobViewSchema>;

export const patchItemSchema = z
  .object({
    title: z.string().max(1000).optional(),
    read_status: z.enum(READ_STATUSES).optional(),
  })
  .refine((v) => Object.keys(v).length > 0, { message: "empty patch" });

// --- Notes ---

export const createNoteSchema = z.object({
  body: z.string().min(1).max(50_000),
});
export const patchNoteSchema = z.object({
  body: z.string().min(1).max(50_000),
  version: z.number().int().positive(),
});

// --- Search & ask ---

export const searchModeSchema = z.enum(["fts", "vector", "hybrid"]).default("hybrid");

export const askRequestSchema = z.object({
  question: z.string().min(1).max(4000),
  item_id: z.string().optional(),
  include_history: z.boolean().default(false),
});

// --- Shares ---

export const createShareSchema = z.object({
  item_id: z.string(),
  snapshot_id: z.string(),
  content_revision_id: z.string().optional(),
  include_full_text: z.boolean().default(false),
  allowed_artifact_ids: z.array(z.string()).default([]),
  expires_in_days: z.number().int().positive().max(365).nullable().default(7),
});

// --- Tokens ---

export const TOKEN_SCOPES = [
  "items:write",
  "items:read",
  "search:read",
  "notes:write",
  "shares:write",
  "exports:write",
  "settings:write",
  "tokens:manage",
] as const;
export type TokenScope = (typeof TOKEN_SCOPES)[number];

export const createTokenSchema = z.object({
  name: z.string().min(1).max(100),
  scopes: z.array(z.enum(TOKEN_SCOPES)).min(1),
  expires_in_days: z.number().int().positive().max(365).nullable().default(null),
});

// --- Artifacts ---

export const artifactViewSchema = z.object({
  id: z.string(),
  type: z.enum(ARTIFACT_TYPES),
  status: z.enum(ARTIFACT_STATUSES),
  output: z.unknown(),
  created_at: z.string(),
});
