import type { ToolSpec } from "@recally/agent-runtime";
import {
  type BodyCandidate,
  type ContentBlock,
  checkUrlTarget,
  type Observation,
} from "@recally/capture";
import { extractContentInput, webFetchInput } from "@recally/contracts";
import { AppError, DEFAULT_LIMITS, nowIso, sha256Hex } from "@recally/domain";
import type { ToolDeps } from "../deps";

// web_fetch (§6.2): ONE controlled HTTP GET + deterministic extraction +
// evidence persistence. It never starts a browser, never calls a secondary
// fetch service, never declares capture success — the agent decides.

const HEADER_ALLOWLIST = new Set([
  "content-type",
  "content-language",
  "last-modified",
  "etag",
  "content-length",
  "x-robots-tag",
]);

export interface FetchedSource {
  finalUrl: string;
  status: number;
  contentType: string;
  headers: Record<string, string>;
  body: Uint8Array;
  sha256: string;
}

export async function safeFetch(
  rawUrl: string,
  fetchFn: typeof fetch,
  opts: { maxBytes?: number; timeoutMs?: number; maxRedirects?: number } = {},
): Promise<FetchedSource> {
  const maxBytes = opts.maxBytes ?? DEFAULT_LIMITS.fetch.maxHtmlBytes;
  const timeoutMs = opts.timeoutMs ?? DEFAULT_LIMITS.fetch.timeoutMs;
  const maxRedirects = opts.maxRedirects ?? DEFAULT_LIMITS.fetch.maxRedirects;

  let url = checkUrlTarget(rawUrl).toString();

  for (let hop = 0; hop <= maxRedirects; hop++) {
    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort("timeout"), timeoutMs);
    let res: Response;
    try {
      res = await fetchFn(url, {
        signal: controller.signal,
        redirect: "manual",
        headers: {
          "user-agent": "recally-archive/0.1 (+personal reading archive)",
          accept: "text/html,application/xhtml+xml,application/pdf;q=0.9,*/*;q=0.5",
          // Only advertise encodings we can decode ourselves (workerd's
          // DecompressionStream has no br); otherwise brotli bytes reach the
          // parser as garbage.
          "accept-encoding": "gzip, deflate",
        },
      });
    } catch (err) {
      if (controller.signal.aborted) {
        throw new AppError("source_unavailable", `fetch timeout after ${timeoutMs}ms`, {
          retryable: true,
        });
      }
      throw new AppError(
        "source_unavailable",
        `fetch failed: ${err instanceof Error ? err.message : String(err)}`,
        { retryable: true },
      );
    } finally {
      clearTimeout(timer);
    }

    if (res.status >= 300 && res.status < 400) {
      const location = res.headers.get("location");
      if (!location) {
        throw new AppError("source_unavailable", `redirect ${res.status} without location`);
      }
      url = checkUrlTarget(new URL(location, url).toString()).toString();
      continue;
    }

    const contentType = res.headers.get("content-type") ?? "";
    let body = await readLimited(res, maxBytes);
    // A still-set content-encoding means the runtime did not decompress for
    // us — do it ourselves so downstream parsers see real bytes.
    const encoding = res.headers.get("content-encoding")?.toLowerCase();
    if (encoding === "gzip" || encoding === "deflate") {
      const ds = new DecompressionStream(encoding);
      body = new Uint8Array(
        await new Response(new Blob([new Uint8Array(body)]).stream().pipeThrough(ds)).arrayBuffer(),
      );
    }
    // brotli has no DecompressionStream in workerd — leave it as evidence and
    // let extraction report unavailable instead of producing garbage.
    const headers: Record<string, string> = {};
    res.headers.forEach((v, k) => {
      if (HEADER_ALLOWLIST.has(k.toLowerCase())) headers[k.toLowerCase()] = v;
    });
    return {
      finalUrl: url,
      status: res.status,
      contentType,
      headers,
      body,
      sha256: await sha256Hex(body),
    };
  }
  throw new AppError("policy_denied", `more than ${maxRedirects} redirects`);
}

async function readLimited(res: Response, maxBytes: number): Promise<Uint8Array> {
  const reader = res.body?.getReader();
  if (!reader) return new Uint8Array(await res.arrayBuffer());
  const parts: Uint8Array[] = [];
  let total = 0;
  for (;;) {
    const { done, value } = await reader.read();
    if (done) break;
    total += value.byteLength;
    if (total > maxBytes) {
      await reader.cancel();
      throw new AppError("incomplete_content", `body exceeded ${maxBytes} bytes`, {
        nextAction: "archive_partial",
      });
    }
    parts.push(value);
  }
  const out = new Uint8Array(total);
  let offset = 0;
  for (const c of parts) {
    out.set(c, offset);
    offset += c.byteLength;
  }
  return out;
}

// --- deterministic extraction ---

const BLOCK_TAGS = new Set([
  "P",
  "H1",
  "H2",
  "H3",
  "H4",
  "H5",
  "H6",
  "PRE",
  "TABLE",
  "UL",
  "OL",
  "BLOCKQUOTE",
  "FIGURE",
  "IMG",
]);
const SKIP_TAGS = new Set([
  "SCRIPT",
  "STYLE",
  "NOSCRIPT",
  "TEMPLATE",
  "NAV",
  "HEADER",
  "FOOTER",
  "ASIDE",
  "FORM",
  "IFRAME",
]);
const INLINE_TAGS = new Set([
  "A",
  "ABBR",
  "B",
  "BDI",
  "BDO",
  "BR",
  "CITE",
  "CODE",
  "DATA",
  "DFN",
  "EM",
  "I",
  "KBD",
  "MARK",
  "Q",
  "S",
  "SAMP",
  "SMALL",
  "SPAN",
  "STRONG",
  "SUB",
  "SUP",
  "TIME",
  "U",
  "VAR",
  "WBR",
]);
// Past this size a "block" element is really a layout container — HN nests
// whole comment trees in tables, and one atomic block leaves the agent no
// usable ranges to select.
const MAX_ATOMIC_BLOCK_CHARS = 8_000;

function blockKind(tag: string): ContentBlock["kind"] {
  switch (tag) {
    case "PRE":
      return "code";
    case "TABLE":
      return "table";
    case "UL":
    case "OL":
      return "list";
    case "BLOCKQUOTE":
      return "quote";
    case "IMG":
    case "FIGURE":
      return "image";
    default:
      return tag.startsWith("H") ? "heading" : "paragraph";
  }
}

export interface ParsedSource {
  blocks: ContentBlock[];
  links: Array<{ id: string; text: string; href: string }>;
  titles: string[];
  authors: string[];
}

export async function parseBlocks(html: string): Promise<ParsedSource> {
  const { parseHTML } = await import("linkedom");
  const { document } = parseHTML(html);
  const blocks: ContentBlock[] = [];
  const links: ParsedSource["links"] = [];
  let seq = 0;
  let linkSeq = 0;

  const walk = (node: {
    nodeType: number;
    tagName: string;
    textContent?: string | null;
    getAttribute?: (n: string) => string | null;
    childNodes?: ArrayLike<unknown>;
  }) => {
    if (node.nodeType !== 1) return;
    const tag = node.tagName;
    if (SKIP_TAGS.has(tag)) return;
    const text = (node.textContent ?? "").replace(/\s+/g, " ").trim();
    if (BLOCK_TAGS.has(tag) && (text.length <= MAX_ATOMIC_BLOCK_CHARS || text.length === 0)) {
      if (text.length > 0 || tag === "IMG" || tag === "FIGURE") {
        blocks.push({ id: `b${++seq}`, kind: blockKind(tag), text });
      }
      return;
    }
    // Container or oversized block element: flush direct inline text as a
    // paragraph block (e.g. text sitting bare in div/td like HN commtext),
    // then descend into non-inline children in document order.
    let inline = "";
    const flush = () => {
      const t = inline.replace(/\s+/g, " ").trim();
      if (t) blocks.push({ id: `b${++seq}`, kind: "paragraph", text: t });
      inline = "";
    };
    for (const child of Array.from(node.childNodes ?? [])) {
      const c = child as {
        nodeType: number;
        tagName?: string;
        textContent?: string | null;
        getAttribute?: (n: string) => string | null;
      };
      if (c.nodeType === 3 || (c.nodeType === 1 && INLINE_TAGS.has(c.tagName ?? ""))) {
        inline += ` ${c.textContent ?? ""}`;
        if (c.tagName === "A") {
          const href = c.getAttribute?.("href") ?? "";
          const t = (c.textContent ?? "").trim();
          if (href && t) links.push({ id: `l${++linkSeq}`, text: t.slice(0, 200), href });
        }
      } else {
        flush();
        walk(c as never);
      }
    }
    flush();
  };
  // No document.body fallback: its linkedom getter itself throws when the
  // document has no documentElement (empty/plain-text/JSON bodies).
  const root = document.documentElement;
  if (root) walk(root as never);

  const meta = (sel: string, attr: string) =>
    document.querySelector(sel)?.getAttribute(attr)?.trim();
  const titles = [
    meta("meta[property='og:title']", "content"),
    document.querySelector("title")?.textContent?.trim(),
    document.querySelector("h1")?.textContent?.trim(),
  ].filter((t): t is string => Boolean(t));
  const authors = [
    meta("meta[name='author']", "content"),
    meta("meta[property='article:author']", "content"),
  ].filter((t): t is string => Boolean(t));

  return { blocks, links, titles, authors };
}

// Extractors all emit the same block/provenance contract (§17.1).
export async function extractCandidates(
  html: string,
  extractor: string,
): Promise<{ parsed: ParsedSource; candidate: BodyCandidate }> {
  const parsed = await parseBlocks(html);
  const all = parsed.blocks.map((b) => b.id);
  const chars = (ids: string[]) => {
    const byId = new Map(parsed.blocks.map((b) => [b.id, b]));
    return ids.reduce((sum, id) => sum + (byId.get(id)?.text.length ?? 0), 0);
  };

  if (extractor === "defuddle") {
    const { Defuddle } = await import("defuddle/node");
    const doc = await Defuddle(html);
    // Map defuddle's article html back onto our block list by parsing it.
    if (doc.content) {
      const inner = await parseBlocks(doc.content);
      const ids = inner.blocks.map((b) => b.id);
      return {
        parsed: {
          ...parsed,
          blocks: inner.blocks.length ? inner.blocks : parsed.blocks,
        },
        candidate: {
          mode: "defuddle",
          blockIds: ids.length ? ids : all,
          charCount: chars(ids.length ? ids : all),
        },
      };
    }
    return {
      parsed,
      candidate: { mode: "defuddle", blockIds: all, charCount: chars(all) },
    };
  }
  if (extractor === "structure") {
    // First heading to last content block.
    const firstHeading = parsed.blocks.findIndex((b) => b.kind === "heading");
    const ids = firstHeading >= 0 ? all.slice(firstHeading) : all;
    return {
      parsed,
      candidate: { mode: "structure", blockIds: ids, charCount: chars(ids) },
    };
  }
  return {
    parsed,
    candidate: { mode: "wide", blockIds: all, charCount: chars(all) },
  };
}

export function buildObservation(
  fetched: FetchedSource,
  parsed: ParsedSource,
  candidate: BodyCandidate,
): Observation {
  const byId = new Map(parsed.blocks.map((b) => [b.id, b]));
  const selected = candidate.blockIds
    .map((id) => byId.get(id))
    .filter((b): b is ContentBlock => !!b);
  const bodyText = selected.map((b) => b.text).join("\n");
  const navShare = parsed.blocks.length ? selected.length / parsed.blocks.length : 0;
  const pageKind: Observation["pageKind"] =
    fetched.status === 401 || fetched.status === 403
      ? "login_wall"
      : parsed.blocks.length < 3 && bodyText.length < 200
        ? "shell"
        : navShare > 0.3
          ? "article"
          : "unknown";

  return {
    url: fetched.finalUrl,
    fetchedAt: nowIso(),
    responseStatus: fetched.status,
    contentType: fetched.contentType,
    titleCandidates: parsed.titles,
    authorCandidates: parsed.authors,
    pageKind,
    blockIndex: parsed.blocks.map((b) => ({
      id: b.id,
      kind: b.kind,
      chars: b.text.length,
    })),
    candidates: [candidate],
    headText: bodyText.slice(0, 800),
    tailText: bodyText.slice(-800),
    links: parsed.links.slice(0, 100),
    sourceTruncated: false,
    modelTruncated: false,
    missingResources: [],
    securityTags: [],
  };
}

// --- tool specs ---

export function webFetchTool(deps: ToolDeps): ToolSpec {
  return {
    name: "web_fetch",
    label: "Fetch page",
    description:
      "One controlled HTTP GET of a URL plus deterministic body extraction. Returns an observation with source id, title/body candidates, head/tail text, links and truncation flags. Evidence is persisted before you see it. Never chains to other acquisition strategies.",
    schema: webFetchInput,
    async execute(ctx, args) {
      const { url_ref } = webFetchInput.parse(args);
      const fetched = await safeFetch(url_ref, deps.fetchFn ?? fetch);
      const source = await deps.runStore.saveSource(ctx, {
        url: fetched.finalUrl,
        kind: "response_body",
        contentType: fetched.contentType,
        body: fetched.body,
      });
      const html = new TextDecoder().decode(fetched.body);
      const { parsed, candidate } = await extractCandidates(html, "wide");
      const observation = buildObservation(fetched, parsed, candidate);
      observation.source = source;
      const obsRef = await deps.runStore.saveObservation(ctx, observation);
      return {
        content: summarizeObservation(observation),
        details: observation,
        resultRef: obsRef,
      };
    },
  };
}

export function extractContentTool(deps: ToolDeps): ToolSpec {
  return {
    name: "extract_content",
    label: "Re-extract source",
    description:
      "Run a different deterministic extractor (defuddle | structure | wide) over an already-saved source. No network access.",
    schema: extractContentInput,
    async execute(ctx, args) {
      const { source_id, extractor } = extractContentInput.parse(args);
      const stored = await deps.runStore.getSource(ctx, source_id);
      if (!stored) throw new AppError("invalid_input", `unknown source ${source_id}`);
      const html = await deps.evidence.getText(stored.bodyKey);
      if (!html) throw new AppError("internal", `missing source body ${stored.bodyKey}`);
      const { parsed, candidate } = await extractCandidates(html, extractor);
      const observation = buildObservation(
        {
          finalUrl: stored.ref.url,
          status: 200,
          contentType: "text/html",
          headers: {},
          body: new Uint8Array(),
          sha256: stored.ref.sha256,
        },
        parsed,
        candidate,
      );
      observation.source = stored.ref;
      const obsRef = await deps.runStore.saveObservation(ctx, observation);
      return {
        content: summarizeObservation(observation),
        details: observation,
        resultRef: obsRef,
      };
    },
  };
}

function summarizeObservation(o: Observation): string {
  const total = o.blockIndex.reduce((s, b) => s + b.chars, 0);
  const first = o.blockIndex[0]?.id;
  const last = o.blockIndex[o.blockIndex.length - 1]?.id;
  return [
    `pageKind=${o.pageKind} status=${o.responseStatus ?? "?"} blocks=${o.blockIndex.length} chars=${total}`,
    o.blockIndex.length ? `block ids: ${first}..${last}` : "",
    o.source ? `source: ${o.source.sourceId}` : "",
    `titles: ${o.titleCandidates.slice(0, 3).join(" | ") || "-"}`,
    `head: ${o.headText.slice(0, 300)}`,
    `tail: ${o.tailText.slice(-300)}`,
    o.missingResources.length ? `missing: ${o.missingResources.join(", ")}` : "",
    o.securityTags.length ? `security: ${o.securityTags.join(", ")}` : "",
  ]
    .filter(Boolean)
    .join("\n");
}
