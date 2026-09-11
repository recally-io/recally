import type { WorkflowBinding } from "@recally/platform-cloudflare";

export interface Env {
  DB: D1Database;
  ARCHIVE_BUCKET: R2Bucket;
  BACKUP_BUCKET: R2Bucket;
  VECTOR_INDEX: VectorizeIndex;
  AI: Ai;
  BROWSER: Fetcher;
  INGEST_WORKFLOW: WorkflowBinding;
  ENRICH_WORKFLOW: WorkflowBinding;
  INDEX_WORKFLOW: WorkflowBinding;
  DIGEST_WORKFLOW: WorkflowBinding;
  EXPORT_WORKFLOW: WorkflowBinding;
  PURGE_WORKFLOW: WorkflowBinding;
  ENVIRONMENT?: string;
  CAPTURE_MODEL?: string;
  VERIFY_MODEL?: string;
  SUMMARY_MODEL?: string;
  ANSWER_MODEL?: string;
  EMBEDDING_MODEL?: string;
  AI_GATEWAY_NAME?: string;
  CLOUDFLARE_ACCOUNT_ID?: string;
}
