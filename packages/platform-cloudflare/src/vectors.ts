// Vectorize adapter (plan §11.2/§11.4). Upsert success ≠ queryable — the
// caller tracks pending → submitted → queryable in D1 and probes before
// marking queryable. namespace scopes queries; it is not an auth boundary.

export class VectorIndex {
  constructor(private readonly index: VectorizeIndex) {}

  async upsert(
    vectors: Array<{
      id: string;
      values: number[];
      namespace: string;
      metadata?: Record<string, string>;
    }>,
  ): Promise<void> {
    await this.index.upsert(
      vectors.map((v) => ({
        id: v.id,
        values: v.values,
        namespace: v.namespace,
        metadata: v.metadata ?? {},
      })),
    );
  }

  async query(
    namespace: string,
    vector: number[],
    topK: number,
  ): Promise<Array<{ id: string; score: number }>> {
    const res = await this.index.query(vector, {
      topK,
      namespace,
      returnMetadata: "none",
    });

    return res.matches.map((m) => ({ id: m.id, score: m.score }));
  }

  async deleteByIds(ids: string[]): Promise<void> {
    if (ids.length) await this.index.deleteByIds(ids);
  }
}
