import type { ToolSpec } from "@recally/agent-runtime";
import { readSourceInput } from "@recally/contracts";
import { AppError } from "@recally/domain";
import type { ToolDeps } from "../deps";
import { parseBlocks } from "../web-fetch";

// read_source: bounded re-read of already-saved evidence (§6.2). Never reads
// outside this run's sources, never touches the network.
export function readSourceTool(deps: ToolDeps): ToolSpec {
  const MAX_CHARS = 20_000;
  return {
    name: "read_source",
    label: "Read source blocks",
    description:
      "Read a bounded range of blocks from a source you already fetched this run. Block ids come from earlier observations.",
    schema: readSourceInput,
    async execute(ctx, args) {
      const { source_id, block_ids, start_block, end_block } = readSourceInput.parse(args);
      const stored = await deps.runStore.getSource(ctx, source_id);
      if (!stored) {
        const known = (await deps.runStore.listSources(ctx)).map((s) => s.sourceId);
        throw new AppError(
          "invalid_input",
          `unknown source ${source_id}; sources saved this run: ${known.join(", ") || "none yet"}`,
        );
      }
      const html = await deps.evidence.getText(stored.bodyKey);
      if (html === null) throw new AppError("internal", `missing source body ${stored.bodyKey}`);

      const { blocks } = await parseBlocks(html);
      // Non-HTML bodies (adapter JSON, plain text) parse to zero blocks —
      // return the raw body instead of an empty read.
      if (blocks.length === 0) {
        return {
          content:
            html.slice(0, MAX_CHARS) +
            (html.length > MAX_CHARS ? "\n[truncated — body shown raw, no blocks]" : ""),
          details: { source_id, blocks: [], truncated: html.length > MAX_CHARS },
        };
      }
      const rangeHint = `valid blocks: ${blocks[0]?.id}..${blocks[blocks.length - 1]?.id} (${blocks.length} total)`;
      let selected = blocks;
      if (block_ids?.length) {
        const wanted = new Set(block_ids);
        selected = blocks.filter((b) => wanted.has(b.id));
        if (selected.length === 0) {
          throw new AppError("invalid_input", `no matching block ids; ${rangeHint}`);
        }
      } else if (start_block || end_block) {
        const start = blocks.findIndex((b) => b.id === start_block);
        const end = blocks.findIndex((b) => b.id === end_block);
        if (start === -1 || end === -1 || end < start) {
          throw new AppError("invalid_input", `block range not found in source; ${rangeHint}`);
        }
        selected = blocks.slice(start, end + 1);
      }

      let out = "";
      let truncated = false;
      for (const b of selected) {
        const line = `[${b.id}] ${b.text}\n`;
        if (out.length + line.length > MAX_CHARS) {
          truncated = true;
          break;
        }
        out += line;
      }
      return {
        content: out + (truncated ? "\n[truncated — narrow the range]" : ""),
        details: { source_id, blocks: selected.map((b) => b.id), truncated },
      };
    },
  };
}
