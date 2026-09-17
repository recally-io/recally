import type { ContentBlock } from "@recally/capture";

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

const BLOCK_KIND_BY_TAG = {
  PRE: "code",
  TABLE: "table",
  UL: "list",
  OL: "list",
  BLOCKQUOTE: "quote",
  IMG: "image",
  FIGURE: "image",
} satisfies Record<string, ContentBlock["kind"]>;

type DomNode = {
  nodeType: number;
  tagName?: string;
  textContent?: string | null;
  getAttribute?: (name: string) => string | null;
  childNodes?: ArrayLike<DomNode>;
};

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

  const walk = (node: DomNode) => {
    if (node.nodeType !== 1) return;
    const tag = node.tagName ?? "";

    if (SKIP_TAGS.has(tag)) return;
    const text = (node.textContent ?? "").replace(/\s+/g, " ").trim();

    if (BLOCK_TAGS.has(tag) && (text.length <= MAX_ATOMIC_BLOCK_CHARS || text.length === 0)) {
      if (text.length > 0 || tag === "IMG" || tag === "FIGURE") {
        // SAFETY: the own-property check proves tag is a key of this table.
        const kind = Object.hasOwn(BLOCK_KIND_BY_TAG, tag)
          ? BLOCK_KIND_BY_TAG[tag as keyof typeof BLOCK_KIND_BY_TAG]
          : tag.startsWith("H")
            ? "heading"
            : "paragraph";

        blocks.push({ id: `b${++seq}`, kind, text });
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

    for (const c of Array.from(node.childNodes ?? [])) {
      if (c.nodeType === 3 || (c.nodeType === 1 && INLINE_TAGS.has(c.tagName ?? ""))) {
        inline += ` ${c.textContent ?? ""}`;

        if (c.tagName === "A") {
          const href = c.getAttribute?.("href") ?? "";
          const t = (c.textContent ?? "").trim();

          if (href && t) links.push({ id: `l${++linkSeq}`, text: t.slice(0, 200), href });
        }
      } else {
        flush();
        walk(c);
      }
    }

    flush();
  };

  // No document.body fallback: its linkedom getter itself throws when the
  // document has no documentElement (empty/plain-text/JSON bodies).
  const root = document.documentElement;

  if (root) walk(root);

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
