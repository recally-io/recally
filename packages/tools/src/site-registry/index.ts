import type { ToolSpec } from "@recally/agent-runtime";
import { siteListInput, siteRunInput } from "@recally/contracts";
import { AppError, nowIso } from "@recally/domain";
import { getAdapter, listAdapters } from "@recally/site-adapters";
import type { ToolDeps } from "../deps";

// site_list / site_run (§6.2): discovery exposes only matching, enabled
// adapters; run persists the structured record as evidence. Adapter output is
// untrusted observation like everything else.

export function siteListTool(): ToolSpec {
  return {
    name: "site_list",
    label: "List site adapters",
    description:
      "List enabled site adapters matching a URL. Use for app-like records (repos, posts, items) where a structured record beats scraped HTML.",
    schema: siteListInput,
    async execute(_ctx, args) {
      const { url_ref } = siteListInput.parse(args);
      const matches = listAdapters(new URL(url_ref));

      return {
        content: matches.length
          ? matches.map((a) => `${a.id} — ${a.summary}`).join("\n")
          : "no adapters match this url",
        details: { url: url_ref, adapters: matches },
      };
    },
  };
}

export function siteRunTool(deps: ToolDeps): ToolSpec {
  return {
    name: "site_run",
    label: "Run site adapter",
    description:
      "Execute a registered site adapter for a URL. Returns the structured record with provenance. The adapter cannot publish an archive itself.",
    schema: siteRunInput,
    async execute(ctx, args) {
      const { adapter_id, url, input } = siteRunInput.parse(args);

      if (!deps.adapterFetch) {
        throw new AppError("not_implemented", "adapter fetch not configured");
      }

      const adapter = getAdapter(adapter_id);
      const parsedArgs = adapter.inputSchema().parse(input);

      const result = await adapter.run(
        { url, args: parsedArgs as Record<string, unknown> },
        ctx,
        deps.adapterFetch,
      );

      const source = await deps.runStore.saveSource(ctx, {
        url: result.canonicalUrl ?? url,
        kind: "adapter_record",
        contentType: "application/json",
        body: JSON.stringify(result.record, null, 2),
      });

      const obsRef = await deps.runStore.saveObservation(ctx, {
        source,
        url,
        fetchedAt: nowIso(),
        pageKind: "record",
        provenance: result.provenance,
        mediaRefs: result.mediaRefs,
      });

      const preview = JSON.stringify(result.record).slice(0, 2000);

      return {
        content: `adapter ${adapter_id} record saved as source ${source.sourceId} (${preview.length} preview chars; adapter records archive whole-body — propose with this sourceId and no block ranges):\n${preview}`,
        details: { adapter: adapter_id, mediaRefs: result.mediaRefs },
        resultRef: obsRef,
      };
    },
  };
}
