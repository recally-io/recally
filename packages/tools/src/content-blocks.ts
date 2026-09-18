// Ranges are inclusive and never mutate the input. Missing endpoints or a
// reversed range return null; callers own the error message.
import type { ContentBlock } from "@recally/capture";

export function selectBlockRange<T extends { id: string }>(
  blocks: readonly T[],
  startId: string | undefined,
  endId: string | undefined,
): T[] | null {
  const start = blocks.findIndex((block) => block.id === startId);
  const end = blocks.findIndex((block) => block.id === endId);

  return start === -1 || end === -1 || end < start ? null : blocks.slice(start, end + 1);
}

// Committed blocks.jsonl: one ContentBlock per line (§4.4).
export function parseBlocksJsonl(jsonl: string): ContentBlock[] {
  return jsonl
    .split("\n")
    .filter(Boolean)
    .map((line) => JSON.parse(line));
}

export function stringifyBlocksJsonl(blocks: readonly ContentBlock[]): string {
  return blocks.map((block) => JSON.stringify(block)).join("\n");
}
