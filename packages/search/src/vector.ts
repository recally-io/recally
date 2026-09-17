import { sha256Hex } from "@recally/domain";

export const EMBEDDING_VERSION = "bge-m3-v1";

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

export function vectorNamespace(libraryId: string): string {
  return libraryId;
}
