// Shared Cloudflare resource declarations (single Stack, plan §3.2).
//
// No physical name is pinned. Alchemy derives each name from the app, the
// logical id, and the stage, so `--stage prod` is a physically separate copy
// of the whole data plane. Pinning a name here would silently collapse every
// stage onto one resource — the failure mode that made the wrangler era's
// stage isolation nominal.
import * as Cloudflare from "alchemy/Cloudflare";

// The shared D1 database. `migrations` applies ./migrations on create and
// records applied files in alchemy's `__alchemy_migrations` table.
export const Database = Cloudflare.D1.Database("Database", {
  migrations: "./migrations",
});

// Evidence + archive objects (plan §4.4). Read via typed clients in the
// workers; the native R2Bucket flows to the platform-neutral stores.
export const ArchiveBucket = Cloudflare.R2.Bucket("ArchiveBucket");

// Vector index for chunk embeddings (bge-m3 → 1024 dims, cosine). The index
// is immutable: changing dimensions/metric would replace it.
export const VectorIndex = Cloudflare.Vectorize.Index("VectorIndex", {
  dimensions: 1024,
  metric: "cosine",
});
