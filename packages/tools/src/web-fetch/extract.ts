import type { BodyCandidate, ContentBlock, Observation } from "@recally/capture";
import { nowIso } from "@recally/domain";
import type { FetchedSource } from "./fetch";
import { parseBlocks, type ParsedSource } from "./parse";

type ExtractorHandler = () => Promise<{ parsed: ParsedSource; candidate: BodyCandidate }>;

// Extractors all emit the same block/provenance contract (§17.1).
async function extractCandidates(
  html: string,
  extractor: string,
): Promise<{ parsed: ParsedSource; candidate: BodyCandidate }> {
  const parsed = await parseBlocks(html);
  const all = parsed.blocks.map((b) => b.id);

  const byId = new Map(parsed.blocks.map((b) => [b.id, b]));

  const makeCandidate = (mode: string, blockIds: string[]): BodyCandidate => ({
    mode,
    blockIds,
    charCount: blockIds.reduce((sum, id) => sum + (byId.get(id)?.text.length ?? 0), 0),
  });

  const EXTRACTORS = {
    async defuddle() {
      const { Defuddle } = await import("defuddle/node");
      const doc = await Defuddle(html);

      // Map defuddle's article html back onto our block list by parsing it.
      if (doc.content) {
        const inner = await parseBlocks(doc.content);
        const ids = inner.blocks.map((b) => b.id);

        return {
          parsed: { ...parsed, blocks: inner.blocks.length ? inner.blocks : parsed.blocks },
          candidate: makeCandidate("defuddle", ids.length ? ids : all),
        };
      }

      return { parsed, candidate: makeCandidate("defuddle", all) };
    },
    async structure() {
      const firstHeading = parsed.blocks.findIndex((b) => b.kind === "heading");
      const ids = firstHeading >= 0 ? all.slice(firstHeading) : all;

      return { parsed, candidate: makeCandidate("structure", ids) };
    },
    async wide() {
      return { parsed, candidate: makeCandidate("wide", all) };
    },
  } satisfies Record<string, ExtractorHandler>;

  // SAFETY: the own-property check excludes unknown names and inherited properties.
  const handler = Object.hasOwn(EXTRACTORS, extractor)
    ? EXTRACTORS[extractor as keyof typeof EXTRACTORS]
    : EXTRACTORS.wide;

  return handler();
}

export async function extractObservation(
  html: string,
  extractor: string,
  fetched: Pick<FetchedSource, "finalUrl" | "status" | "contentType">,
): Promise<Observation> {
  const { parsed, candidate } = await extractCandidates(html, extractor);
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
