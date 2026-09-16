// Runtime environment handed to the platform-neutral code. Under alchemy the
// typed binding clients are closed over in each worker/workflow Construction
// phase; `.raw` exposes the native objects this interface describes.
import type { Ai, D1Database, Fetcher, R2Bucket } from "@cloudflare/workers-types";

export interface JobsEnv {
  DB: D1Database;
  ARCHIVE_BUCKET: R2Bucket;
  AI: Ai;
  CAPTURE_MODEL?: string | undefined;
  VERIFY_MODEL?: string | undefined;
  SUMMARY_MODEL?: string | undefined;
  ANSWER_MODEL?: string | undefined;
  EMBEDDING_MODEL?: string | undefined;
  AI_GATEWAY_NAME?: string | undefined;
  CLOUDFLARE_ACCOUNT_ID?: string | undefined;
}

// Only the ingest (capture agent) runtime can open browser episodes.
export interface IngestEnv extends JobsEnv {
  BROWSER: Fetcher;
}
