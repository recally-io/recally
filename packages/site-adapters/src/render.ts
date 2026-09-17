// Deterministic record → readable-block renderers (§7.1). The raw JSON record
// stays in the source document as evidence; these blocks become the archived
// article. Unknown record kinds return null and the committer falls back to
// archiving the record body verbatim.

export interface RenderedBlock {
  kind: "heading" | "paragraph" | "quote" | "code";
  text: string;
}

export function renderAdapterRecord(bodyJson: string): RenderedBlock[] | null {
  let record: { kind?: string };

  try {
    record = JSON.parse(bodyJson);
  } catch {
    return null;
  }

  if (record.kind === "hn_item") return renderHnItem(record as never);

  if (record.kind === "github_repo") return renderGithubRepo(record as never);

  return null;
}

function formatMeta(parts: Array<string | null | undefined>): string {
  return parts.filter(Boolean).join(" · ");
}

// --- HN ---

interface HnNode {
  author?: string | null;
  text?: string | null;
  children?: HnNode[];
}

interface HnItem extends HnNode {
  id: number;
  type?: string | null;
  title?: string | null;
  url?: string | null;
  points?: number | null;
  created_at?: string | null;
  comment_count_fetched?: number;
  comments_truncated?: boolean;
}

function renderHnItem(r: HnItem): RenderedBlock[] {
  const blocks: RenderedBlock[] = [];

  if (r.title) blocks.push({ kind: "heading", text: r.title });

  const meta = formatMeta([
    r.points != null ? `${r.points} points` : null,
    r.author ? `by ${r.author}` : null,
    r.created_at ? new Date(r.created_at).toDateString() : null,
    r.comment_count_fetched != null
      ? `${r.comment_count_fetched} comments${r.comments_truncated ? " (truncated)" : ""}`
      : null,
  ]);

  if (meta) blocks.push({ kind: "paragraph", text: meta });

  if (r.url) blocks.push({ kind: "paragraph", text: r.url });

  if (r.text) blocks.push({ kind: "paragraph", text: htmlToText(r.text) });

  if (r.children?.length) {
    blocks.push({ kind: "heading", text: "Comments" });

    const walk = (nodes: HnNode[], depth: number) => {
      for (const n of nodes) {
        const body = htmlToText(n.text ?? "");

        if (!body) {
          walk(n.children ?? [], depth);
          continue;
        }

        // Depth is carried as extra `>` markers inside the quote text, so the
        // committed markdown is a proper nested blockquote.
        blocks.push({
          kind: "quote",
          text: `${"> ".repeat(depth)}${n.author ?? "[deleted]"} — ${body}`,
        });
        walk(n.children ?? [], depth + 1);
      }
    };

    walk(r.children, 0);
  }

  return blocks;
}

// HN comment bodies are small HTML fragments (<p>, <a>, <i>, <pre>). Flatten
// to one readable line per comment; thread structure carries the nesting.
function htmlToText(html: string): string {
  return decodeEntities(
    html
      .replace(/<a\s+href="([^"]*)"[^>]*>([\s\S]*?)<\/a>/gi, (_m, href: string, t: string) =>
        t.trim() === href ? href : `${t} (${href})`,
      )
      .replace(/<(p|div|br|pre|li|blockquote)[^>]*>/gi, " ")
      .replace(/<[^>]+>/g, "")
      .replace(/\s+/g, " ")
      .trim(),
  );
}

const ENTITIES: Record<string, string> = {
  amp: "&",
  lt: "<",
  gt: ">",
  quot: '"',
  apos: "'",
  nbsp: " ",
};

function decodeEntities(s: string): string {
  return s
    .replace(/&#x([0-9a-f]+);/gi, (_m, h) => String.fromCodePoint(Number.parseInt(h, 16)))
    .replace(/&#(\d+);/g, (_m, d) => String.fromCodePoint(Number.parseInt(d, 10)))
    .replace(/&([a-z]+);/gi, (m, name) => ENTITIES[name.toLowerCase()] ?? m);
}

// --- GitHub ---

interface GhRepo {
  full_name?: string;
  description?: string | null;
  homepage?: string | null;
  language?: string | null;
  stars?: number | null;
  topics?: string[];
  license?: string | null;
  pushed_at?: string | null;
  readme_markdown?: string | null;
}

function renderGithubRepo(r: GhRepo): RenderedBlock[] {
  const blocks: RenderedBlock[] = [];

  if (r.full_name) blocks.push({ kind: "heading", text: r.full_name });

  if (r.description) blocks.push({ kind: "paragraph", text: r.description });

  const meta = formatMeta([
    r.stars != null ? `${r.stars} stars` : null,
    r.language,
    r.license,
    r.pushed_at ? `updated ${new Date(r.pushed_at).toDateString()}` : null,
    r.topics?.length ? r.topics.join(", ") : null,
    r.homepage,
  ]);

  if (meta) blocks.push({ kind: "paragraph", text: meta });

  if (r.readme_markdown) {
    blocks.push({ kind: "heading", text: "README" });
    blocks.push(...markdownToBlocks(r.readme_markdown));
  }

  return blocks;
}

// README arrives as markdown; split it back into semantic blocks so the
// committed article keeps headings and code fences instead of raw markup.
function markdownToBlocks(md: string): RenderedBlock[] {
  const blocks: RenderedBlock[] = [];

  for (const chunk of md.split(/\n{2,}/)) {
    const t = chunk.trim();

    if (!t) continue;

    if (/^#{1,6}\s/.test(t)) {
      const line = t.split("\n")[0] ?? t;
      blocks.push({ kind: "heading", text: line.replace(/^#{1,6}\s+/, "") });
    } else if (t.startsWith("```")) {
      blocks.push({ kind: "code", text: t.replace(/^```\w*\n?/, "").replace(/\n?```$/, "") });
    } else {
      blocks.push({ kind: "paragraph", text: t });
    }
  }

  return blocks;
}
