export const JOB_KINDS = ["ingest", "enrich", "index", "digest", "export", "purge"] as const;

export type JobKind = (typeof JOB_KINDS)[number];

export const JOB_STATUSES = [
  "queued",
  "running",
  "retry_wait",
  "succeeded",
  "needs_input",
  "failed",
  "cancelled",
] as const;

export type JobStatus = (typeof JOB_STATUSES)[number];

export const CONTENT_QUALITIES = ["complete", "partial", "unavailable", "unknown"] as const;

export type ContentQuality = (typeof CONTENT_QUALITIES)[number];

export const RESOURCE_QUALITIES = ["complete", "partial", "not_applicable", "unknown"] as const;

export type ResourceQuality = (typeof RESOURCE_QUALITIES)[number];

export const READ_STATUSES = ["unread", "reading", "read", "archived"] as const;

export type ReadStatus = (typeof READ_STATUSES)[number];

const INDEX_STATES = ["pending", "submitted", "queryable", "stale", "removed"] as const;

export type IndexState = (typeof INDEX_STATES)[number];

export const ARTIFACT_TYPES = ["summary", "digest", "imported_answer", "qa"] as const;

export type ArtifactType = (typeof ARTIFACT_TYPES)[number];

export const ARTIFACT_STATUSES = ["pending", "verified", "unverified_import", "failed"] as const;

export type ArtifactStatus = (typeof ARTIFACT_STATUSES)[number];
