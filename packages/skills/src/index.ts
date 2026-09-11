import { sha256Hex } from "@recally/domain";
import analyzeSkill from "../analyze/SKILL.md";
import captureQuality from "../capture/references/quality.md";
import captureSourceTypes from "../capture/references/source-types.md";
import captureSkill from "../capture/SKILL.md";

// Skill revisions are content hashes: a run pins the exact text it saw, a new
// revision only affects new runs (§5.3, §9.6, ADR-07). Revisions are computed
// once per isolate and cached.

const SKILL_FILES: Record<string, Record<string, string>> = {
  capture: {
    "SKILL.md": captureSkill,
    "references/quality.md": captureQuality,
    "references/source-types.md": captureSourceTypes,
  },
  analyze: {
    "SKILL.md": analyzeSkill,
  },
};

export type SkillName = keyof typeof SKILL_FILES;

export interface LoadedSkill {
  name: SkillName;
  revision: string;
  // Skill body plus references, concatenated for the system prompt.
  prompt: string;
}

const revisionCache = new Map<SkillName, Promise<LoadedSkill>>();

export function loadSkill(name: SkillName): Promise<LoadedSkill> {
  let cached = revisionCache.get(name);
  if (!cached) {
    cached = (async () => {
      const files: Record<string, string> = SKILL_FILES[name]!;
      const ordered = Object.keys(files).sort();
      const prompt = ordered.map((f) => files[f]!).join("\n\n---\n\n");
      const revision = `skill-${(await sha256Hex(prompt)).slice(0, 12)}`;
      return { name, revision, prompt };
    })();
    revisionCache.set(name, cached);
  }
  return cached;
}
