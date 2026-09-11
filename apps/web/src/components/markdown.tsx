// article.md is produced deterministically by the committer: `## ` headings,
// ``` fenced code, and plain paragraphs separated by blank lines (§7.1 — the
// archived original is rendered from source blocks, never model prose).
export function Article({ md }: { md: string }) {
  const blocks = md.split(/\n{2,}/).filter((b) => b.trim());
  return (
    <div className="font-serif text-[16.5px] leading-[1.75] text-ink">
      {blocks.map((b, i) => {
        const key = `${i}-${b.length}`;
        if (b.startsWith("## ")) {
          return (
            <h2 key={key} className="mt-8 mb-3 font-serif text-xl leading-snug font-medium">
              {b.slice(3)}
            </h2>
          );
        }
        if (b.startsWith("```")) {
          const code = b.replace(/^```\w*\n?/, "").replace(/\n?```$/, "");
          return (
            <pre
              key={key}
              className="my-5 overflow-x-auto rounded-lg border border-line bg-paper p-4 font-mono text-[13px] leading-relaxed whitespace-pre-wrap"
            >
              {code}
            </pre>
          );
        }
        if (b.startsWith("> ")) {
          // Nested `> >` markers (adapter comment trees) become indentation.
          const depth = (b.match(/^>+/g)?.[0] ?? ">").length;
          return (
            <blockquote
              key={key}
              className="mt-3 border-l-2 border-line pl-3 text-[15px] leading-relaxed text-ink-2"
              style={{ marginLeft: `${(depth - 1) * 20}px` }}
            >
              {b.replace(/^>+\s?/gm, "")}
            </blockquote>
          );
        }
        return (
          <p key={key} className="mb-5">
            {b}
          </p>
        );
      })}
    </div>
  );
}
