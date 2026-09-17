import { AppError } from "@recally/domain";
import { githubRepoAdapter } from "./github";
import { hackernewsItemAdapter } from "./hackernews";
import type { AdapterDescriptor, SiteAdapter } from "./types";

// Explicit registry — adapters are compiled into the bundle, never loaded
// dynamically (§5.7). Adding a source = adding a file + one line here.
const ADAPTERS: SiteAdapter[] = [githubRepoAdapter, hackernewsItemAdapter];

export function listAdapters(url: URL): AdapterDescriptor[] {
  return ADAPTERS.filter((a) => a.match(url)).map((a) => a.descriptor);
}

export function getAdapter(id: string): SiteAdapter {
  const adapter = ADAPTERS.find((a) => a.descriptor.id === id);

  if (!adapter) throw new AppError("invalid_input", `unknown adapter ${id}`);

  return adapter;
}
