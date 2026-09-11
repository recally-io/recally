import type { ToolContext } from "@recally/capture";
import type { z } from "zod";

// Site adapters (§5.7): audited, compiled-in modules for URLs that denote a
// structured record rather than a page. No runtime eval, no dynamic install.

export interface AdapterDescriptor {
  id: string;
  domains: string[];
  summary: string;
  produces: string; // e.g. "github repo record with readme"
}

export interface AdapterInput {
  url: string;
  args: Record<string, unknown>;
}

export interface AdapterResult {
  // Structured record the agent can judge and archive.
  record: unknown;
  // Where every field came from — the api endpoint, not just the page.
  provenance: { endpointUrls: string[]; fetchedAt: string };
  // Media/attachment refs the agent may archive via archive_asset.
  mediaRefs: Array<{
    url: string;
    kind: string;
    role: "required" | "optional";
  }>;
  // Canonical url if the record has one distinct from the input.
  canonicalUrl?: string;
}

// Adapter network access goes through this injected fetch — same SSRF policy
// as web_fetch, adapters never get a raw fetch.
export type AdapterFetch = (url: string) => Promise<unknown>;

export interface SiteAdapter {
  readonly descriptor: AdapterDescriptor;
  match(url: URL): boolean;
  inputSchema(): z.ZodType<unknown>;
  run(input: AdapterInput, ctx: ToolContext, fetchJson: AdapterFetch): Promise<AdapterResult>;
}
