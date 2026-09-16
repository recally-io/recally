import type { D1Database, Fetcher, R2Bucket } from "@cloudflare/workers-types";

// Runtime environment of the app worker. Under alchemy these are produced by
// the typed binding clients in worker.ts (`.raw` native handles + Config
// values bound during Construction), not by a wrangler-generated Env.
export interface Env {
  DB: D1Database;
  ARCHIVE_BUCKET: R2Bucket;
  ASSETS: Fetcher;
  DEV_LIBRARY_ID?: string;
  CAPTURE_MODEL?: string;
}
