import { describe, expect, it } from "vitest";
import { chunkBlocks } from "./chunk";
import { rrfFuse } from "./rrf";
import { buildFtsQuery, tokenizeForFts } from "./tokenize";

describe("tokenizeForFts", () => {
  it("segments CJK into bigrams and keeps latin words", () => {
    const t = tokenizeForFts("你好世界 hello");
    expect(t).toContain("hello");
    expect(t).toContain("你好");
    expect(t).toContain("世界");
  });
});

describe("buildFtsQuery", () => {
  it("quotes every token so FTS5 syntax can't be injected", () => {
    const q = buildFtsQuery('foo " OR bar');
    expect(q).not.toContain(" OR ");
    expect(q.split(" ").every((t) => t.startsWith('"') && t.endsWith('"'))).toBe(true);
  });

  it("returns empty for empty input", () => {
    expect(buildFtsQuery("   ")).toBe("");
  });
});

describe("rrfFuse", () => {
  it("merges two lists by reciprocal rank", () => {
    const a = { id: "a" };
    const b = { id: "b" };
    const c = { id: "c" };
    const d = { id: "d" };

    const merged = rrfFuse([
      [a, b, c],
      [b, d, a],
    ]);

    expect(merged[0]!.item.id).toBe("b"); // b ranks 2nd and 1st
  });
});

describe("chunkBlocks", () => {
  it("produces chunks with block ranges", () => {
    const blocks = Array.from({ length: 50 }, (_, i) => ({
      id: `b${i}`,
      kind: "paragraph",
      text: "word ".repeat(50),
    }));

    const chunks = chunkBlocks(blocks);
    expect(chunks.length).toBeGreaterThan(1);
    expect(chunks[0]!.blockRange.start).toBe("b0");
  });
});
