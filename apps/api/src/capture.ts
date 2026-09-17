import { resolveModels } from "@recally/ai";
import { modelConfigVersion, PIPELINE_VERSION, POLICY_VERSION } from "@recally/domain";
import { loadSkill } from "@recally/skills";
import type { CaptureRunRevisions } from "@recally/storage";
import type { Env } from "./env";

const TOOLSET_VERSION = "toolset-0001";

// Revision fencing (§9.6): pin the runtime/skill/toolset/model config at run
// creation so a retry replays exactly what it started with.
export async function runRevisions(env: Env): Promise<CaptureRunRevisions> {
  const skill = await loadSkill("capture");
  const mcfg = await modelConfigVersion({ ...resolveModels(env) });

  return {
    agentRuntimeVersion: "pi-0.85.1",
    skillRevision: skill.revision,
    toolsetVersion: TOOLSET_VERSION,
    policyVersion: POLICY_VERSION,
    pipelineVersion: PIPELINE_VERSION,
    modelConfigVersion: mcfg,
  };
}
