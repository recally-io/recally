// Reciprocal rank fusion (plan §11.3): fuse ranked lists by position, never
// add raw scores across FTS and vector sources.
export interface RankedItem {
  id: string;
}

export function rrfFuse<T extends RankedItem>(
  lists: Array<Array<T>>,
  k = 60,
): Array<{ item: T; score: number }> {
  const scores = new Map<string, { item: T; score: number }>();

  for (const list of lists) {
    list.forEach((item, rank) => {
      const entry = scores.get(item.id) ?? { item, score: 0 };
      entry.score += 1 / (k + rank + 1);
      scores.set(item.id, entry);
    });
  }

  return [...scores.values()].sort((a, b) => b.score - a.score);
}
