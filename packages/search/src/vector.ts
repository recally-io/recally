import { sha256Hex } from "@recally/domain";

// Single source of truth is @recally/domain (plan §9.6): the value lands in
// `chunks.embedding_version` and in the vector id, so a second definition here
// would let a model upgrade collide with an old index.
export { EMBEDDING_VERSION } from "@recally/domain";

// Vectorize ids: globally unique, ≤64 bytes, embed the embedding version so a
// model upgrade can never collide with an old index (plan §11.2). namespace
// scopes queries per library but is not an auth boundary.
export async function vectorId(parts: {
  libraryId: string;
  contentRevisionId: string;
  chunkId: string;
  embeddingVersion: string;
}): Promise<string> {
  const raw = `${parts.libraryId}|${parts.contentRevisionId}|${parts.chunkId}|${parts.embeddingVersion}`;

  return `v-${(await sha256Hex(raw)).slice(0, 56)}`;
}
