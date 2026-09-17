import type { ToolSpec } from "@recally/agent-runtime";
import { archiveTools } from "./archive";
import { browserTools } from "./browser";
import type { ToolDeps } from "./deps";
import { readSourceTool } from "./read-source";
import { siteListTool, siteRunTool } from "./site-registry";
import { extractContentTool, webFetchTool } from "./web-fetch";

export * from "./analysis";

export * from "./archive";

export * from "./browser";

export * from "./deps";

export * from "./read-source";

export * from "./site-registry";

export * from "./web-fetch";

// The capture toolset (§5.4). Which tools exist is a runtime choice; which
// strategy to use is the agent's, guided by the skill.
export function captureTools(deps: ToolDeps): ToolSpec[] {
  const tools: ToolSpec[] = [
    webFetchTool(deps),
    extractContentTool(deps),
    readSourceTool(deps),
    siteListTool(deps),
    siteRunTool(deps),
    ...archiveTools(deps),
  ];

  if (deps.browser) tools.push(...browserTools(deps));

  return tools;
}
