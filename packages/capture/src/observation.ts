// Bounded observation returned to the agent after each tool call (§5.5/§6.2).
// Full content lives in the evidence store; the model sees a summary plus
// stable refs and reads more via read_source. Everything here is UNTRUSTED
// evidence, never instructions (§6.7).

export interface ContentBlock {
  id: string; // stable within a source document, e.g. "b3"
  kind: "heading" | "paragraph" | "code" | "table" | "list" | "image" | "quote" | "link";
  text: string;
}

export interface BodyCandidate {
  mode: string; // extractor id, e.g. "defuddle" | "structure" | "wide"
  blockIds: string[];
  charCount: number;
}

export interface SourceRef {
  sourceId: string;
  url: string;
  kind: "response_body" | "rendered_dom" | "screenshot" | "adapter_record" | "manual";
  sha256: string;
  byteSize: number;
  evidenceRef: string; // object key, never handed to the model as writable
}

export interface Observation {
  source?: SourceRef;
  url: string;
  fetchedAt: string;
  responseStatus?: number;
  contentType?: string;
  titleCandidates: string[];
  authorCandidates: string[];
  pageKind: "article" | "shell" | "login_wall" | "record" | "error" | "unknown";
  blockIndex: Array<{ id: string; kind: string; chars: number }>; // index only, not full text
  candidates: BodyCandidate[];
  headText: string;
  tailText: string;
  links: Array<{ id: string; text: string; href: string }>;
  sourceTruncated: boolean; // the remote content itself was cut off
  modelTruncated: boolean; // we shortened what the model sees
  missingResources: string[];
  securityTags: string[]; // e.g. "prompt_injection_suspect", "challenge_page"
}
