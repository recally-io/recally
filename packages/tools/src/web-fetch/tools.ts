import type { ToolOutcome, ToolSpec } from "@recally/agent-runtime";
import type { Observation, ToolContext } from "@recally/capture";
import { extractContentInput, webFetchInput } from "@recally/contracts";
import { AppError } from "@recally/domain";
import type { ToolDeps } from "../deps";
import { extractObservation } from "./extract";
import { safeFetch } from "./fetch";

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
      const observation = await extractObservation(html, "wide", fetched);
      observation.source = source;

      return persistObservation(deps, ctx, observation);
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

      const observation = await extractObservation(html, extractor, {
        finalUrl: stored.ref.url,
        status: 200,
        contentType: "text/html",
      });

      observation.source = stored.ref;

      return persistObservation(deps, ctx, observation);
    },
  };
}

async function persistObservation(
  deps: ToolDeps,
  ctx: ToolContext,
  observation: Observation,
): Promise<ToolOutcome> {
  const resultRef = await deps.runStore.saveObservation(ctx, observation);

  return { content: summarizeObservation(observation), details: observation, resultRef };
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
