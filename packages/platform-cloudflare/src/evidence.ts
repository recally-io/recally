import { sha256Hex } from "@recally/domain";

// R2-backed evidence store (plan §8.1). Every tool output lands here BEFORE the
// agent sees a summary — a model failure can never orphan fetched content.
// Keys come only from r2Keys; callers (and the model) cannot pick arbitrary
// buckets or paths.

export interface PutResult {
  key: string;
  sha256: string;
  byteSize: number;
}

export class R2EvidenceStore {
  constructor(private readonly bucket: R2Bucket) {}

  async put(key: string, body: string | Uint8Array, contentType: string): Promise<PutResult> {
    const bytes = typeof body === "string" ? new TextEncoder().encode(body) : body;
    const sha256 = await sha256Hex(bytes);
    await this.bucket.put(key, bytes as unknown as ArrayBuffer, {
      httpMetadata: { contentType },
      sha256,
    });

    return { key, sha256, byteSize: bytes.byteLength };
  }

  // Conditional put: never overwrite committed evidence.
  async putIfAbsent(
    key: string,
    body: string | Uint8Array,
    contentType: string,
  ): Promise<PutResult | null> {
    const bytes = typeof body === "string" ? new TextEncoder().encode(body) : body;
    const sha256 = await sha256Hex(bytes);

    const res = await this.bucket.put(key, bytes as unknown as ArrayBuffer, {
      httpMetadata: { contentType },
      sha256,
      onlyIf: { etagDoesNotMatch: "*" },
    });

    if (!res) return null;

    return { key, sha256, byteSize: bytes.byteLength };
  }

  async getText(key: string): Promise<string | null> {
    const obj = await this.bucket.get(key);

    return obj ? obj.text() : null;
  }

  async getBytes(key: string): Promise<Uint8Array | null> {
    const obj = await this.bucket.get(key);

    return obj ? new Uint8Array(await obj.arrayBuffer()) : null;
  }
}
