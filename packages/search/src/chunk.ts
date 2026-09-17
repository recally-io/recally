import { DEFAULT_LIMITS } from "@recally/domain";
import { tokenizeForFts } from "./tokenize";

export interface SourceBlock {
  id: string;
  kind: string;
  text: string;
}

export interface Chunk {
  blockRange: { start: string; end: string };
  text: string;
}

// Title/paragraph chunking (§11.2): ~600 token target, ~100 token overlap,
// code/table blocks kept whole where possible.
export function chunkBlocks(blocks: SourceBlock[]): Chunk[] {
  const target = DEFAULT_LIMITS.search.chunkTokens;
  const overlap = DEFAULT_LIMITS.search.chunkOverlap;

  const approxTokens = (s: string) => {
    const toks = tokenizeForFts(s).split(" ").filter(Boolean).length;

    return toks || Math.ceil(s.length / 4);
  };

  const chunks: Chunk[] = [];
  let buf: SourceBlock[] = [];
  let tokens = 0;

  const flush = () => {
    if (!buf.length) return;
    chunks.push({
      blockRange: { start: buf[0]!.id, end: buf[buf.length - 1]!.id },
      text: buf.map((b) => b.text).join("\n\n"),
    });
    // Keep a tail of blocks as the next chunk's overlap.
    const kept: SourceBlock[] = [];
    let keptTokens = 0;

    for (let i = buf.length - 1; i >= 0 && keptTokens < overlap; i--) {
      kept.unshift(buf[i]!);
      keptTokens += approxTokens(buf[i]!.text);
    }

    buf = kept;
    tokens = keptTokens;
  };

  for (const b of blocks) {
    const t = approxTokens(b.text);
    const atomic = b.kind === "code" || b.kind === "table";

    // A single oversized block (adapter records archive whole-body, layout
    // tables) must never become one oversized chunk — embedding models have
    // hard context limits. Split on whitespace near the char budget.
    if (t > target * 1.5) {
      flush();
      const maxChars = Math.ceil(target * 1.5 * 4);

      for (let i = 0; i < b.text.length;) {
        let end = Math.min(i + maxChars, b.text.length);

        if (end < b.text.length) {
          const ws = b.text.lastIndexOf(" ", end);

          if (ws > i + maxChars / 2) end = ws;
        }

        chunks.push({ blockRange: { start: b.id, end: b.id }, text: b.text.slice(i, end) });
        i = end;
      }

      continue;
    }

    if (!atomic && tokens + t > target) flush();
    buf.push(b);
    tokens += t;

    if (tokens > target * 1.5) flush();
  }

  flush();

  return chunks;
}
