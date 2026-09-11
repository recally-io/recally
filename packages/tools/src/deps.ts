import type { SourceRef, ToolContext } from "@recally/capture";

// Platform-neutral seams. Cloudflare implementations live in
// packages/platform-cloudflare; the agent core and tools never import
// bindings directly (plan §3.5, ADR-06).

export interface EvidenceStore {
  put(
    key: string,
    body: string | Uint8Array,
    contentType: string,
  ): Promise<{
    key: string;
    sha256: string;
    byteSize: number;
  }>;
  getText(key: string): Promise<string | null>;
  getBytes(key: string): Promise<Uint8Array | null>;
}

export interface StoredSource {
  ref: SourceRef;
  bodyKey: string;
  blocksKey: string | null;
}

// Run-scoped persistence for sources/observations/events. The jobs worker
// implements this over D1 + R2.
export interface RunStore {
  saveSource(
    ctx: ToolContext,
    input: {
      url: string;
      kind: SourceRef["kind"];
      contentType: string;
      body: Uint8Array | string;
    },
  ): Promise<SourceRef>;
  getSource(ctx: ToolContext, sourceId: string): Promise<StoredSource | null>;
  saveObservation(ctx: ToolContext, observation: unknown): Promise<string>;
}

export interface BrowserSessionHandle {
  observe(): Promise<{
    url: string;
    title: string;
    text: string;
    nodeRefs: unknown[];
  }>;
  scroll(amount?: number): Promise<void>;
  click(selector: string): Promise<void>;
  waitFor(selector: string): Promise<void>;
  navigate(url: string): Promise<void>;
  renderedDom(): Promise<string>;
  screenshot(): Promise<Uint8Array>;
  close(): Promise<void>;
}

export interface BrowserFactory {
  open(url: string): Promise<BrowserSessionHandle>;
}

export type AdapterFetch = (url: string) => Promise<unknown>;

// Final commit boundary (§5.10, §8.2). The agent can only propose; the
// committer validates, persists the manifest, and publishes — implemented
// platform-side over D1 + R2.
export interface CommitResult {
  status: "accepted" | "needs_revision" | "rejected";
  snapshotId?: string;
  contentRevisionId?: string;
  problems?: string[];
}

export interface Committer {
  commit(
    ctx: ToolContext,
    proposal: import("@recally/contracts").ArchiveProposal,
  ): Promise<CommitResult>;
}

// Analysis artifacts persist through this hook (AIArtifact rows + R2 object),
// never through the agent writing product tables itself.
export interface PersistedArtifact {
  type: "summary" | "imported_answer" | "digest";
  contentRevisionId: string;
  output: unknown;
  coverage?: unknown;
  modelId: string;
  promptVersion: string;
}

export interface ToolDeps {
  evidence: EvidenceStore;
  runStore: RunStore;
  committer?: Committer;
  persistArtifact?: (ctx: ToolContext, artifact: PersistedArtifact) => Promise<string>;
  // Analysis tools read committed content by revision id; the platform layer
  // resolves the object key so tools never learn the bucket layout.
  readRevisionBlocks?: (ctx: ToolContext, revisionId: string) => Promise<string | null>;
  browser?: BrowserFactory;
  adapterFetch?: AdapterFetch;
  fetchFn?: typeof fetch; // defaults to global fetch; injectable for tests
}
