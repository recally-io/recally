// Shared Cloudflare resource declarations (single Stack, plan §3.2).
//
// Physical names are pinned to the pre-existing wrangler-provisioned
// resources so the first `alchemy deploy --adopt` takes them over instead of
// creating replacements. When prod gets its own resources, make these
// stage-conditional here — the decision belongs next to the declaration.
import * as Cloudflare from "alchemy/Cloudflare";

// The single shared D1 database. Migration history previously applied with
// `wrangler d1 migrations apply` is adopted automatically into alchemy's
// bookkeeping on first apply (docs: cloudflare/data/d1).
export const Database = Cloudflare.D1.Database("Database", {
  name: "recally-dev",
  migrations: "./migrations",
});

// Evidence + archive objects (plan §4.4). Read via typed clients in the
// workers; the native R2Bucket flows to the platform-neutral stores.
export const ArchiveBucket = Cloudflare.R2.Bucket("ArchiveBucket", {
  name: "recally-archive-dev",
});

// Declared now so the bucket is managed by the Stack; the export milestone
// (§15.2) binds it into the jobs worker.
export const BackupBucket = Cloudflare.R2.Bucket("BackupBucket", {
  name: "recally-backup-dev",
});

// Vector index for chunk embeddings (bge-m3 → 1024 dims, cosine). The index
// is immutable: changing dimensions/metric would replace it.
export const VectorIndex = Cloudflare.Vectorize.Index("VectorIndex", {
  name: "recally-chunks-dev",
  dimensions: 1024,
  metric: "cosine",
});
