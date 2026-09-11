// R2 layout per plan §4.4. Objects under runs/ are immutable evidence;
// GC must consult committed manifests, never delete by path alone.

const root = (libraryId: string) => `archive/libraries/${libraryId}`;

export const r2Keys = {
  observation: (lib: string, run: string, attempt: string, obsId: string) =>
    `${root(lib)}/runs/${run}/attempts/${attempt}/observations/${obsId}.json`,

  source: (lib: string, run: string, attempt: string, sourceId: string, name: string) =>
    `${root(lib)}/runs/${run}/attempts/${attempt}/sources/${sourceId}/${name}`,

  decision: (lib: string, run: string, attempt: string, sequence: number) =>
    `${root(lib)}/runs/${run}/attempts/${attempt}/decisions/${sequence}.json`,

  snapshotManifest: (lib: string, snapshot: string) =>
    `${root(lib)}/snapshots/${snapshot}/manifest.json`,

  snapshotResource: (lib: string, snapshot: string, assetId: string) =>
    `${root(lib)}/snapshots/${snapshot}/resources/${assetId}`,

  contentBlocks: (lib: string, revision: string) => `${root(lib)}/content/${revision}/blocks.jsonl`,

  contentArticle: (lib: string, revision: string) => `${root(lib)}/content/${revision}/article.md`,

  contentProvenance: (lib: string, revision: string) =>
    `${root(lib)}/content/${revision}/provenance.json`,

  artifact: (lib: string, artifactId: string) => `${root(lib)}/artifacts/${artifactId}.json`,

  export: (lib: string, exportId: string, name: string) =>
    `${root(lib)}/exports/${exportId}/${name}`,

  piState: (lib: string, run: string, turn: number) =>
    `${root(lib)}/runs/${run}/state/${String(turn).padStart(4, "0")}.json`,
};
