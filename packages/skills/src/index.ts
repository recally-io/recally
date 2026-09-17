import { sha256Hex } from "@recally/domain";

// Skill markdown loads through `?raw` dynamic imports: the worker bundle
// inlines them via rolldown's raw plugin, vitest/vite handle them natively,
// and the alchemy CLI's stack-graph loader never evaluates them (loadSkill
// only runs at runtime). Static `.md` imports would break `alchemy plan` —
// its loader has no text-module support.
const textLoaders: Record<string, Record<string, () => Promise<string>>> = {
  capture: {
    "SKILL.md": () => import("../capture/SKILL.md?raw").then((m) => m.default),
    "references/quality.md": () =>
      import("../capture/references/quality.md?raw").then((m) => m.default),
    "references/source-types.md": () =>
      import("../capture/references/source-types.md?raw").then((m) => m.default),
  },
  analyze: {
    "SKILL.md": () => import("../analyze/SKILL.md?raw").then((m) => m.default),
  },
};

const SKILL_NAMES = Object.keys(textLoaders);

type SkillName = (typeof SKILL_NAMES)[number];

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
      const files = textLoaders[name]!;
      const ordered = Object.keys(files).sort();
      const prompt = (await Promise.all(ordered.map((f) => files[f]!()))).join("\n\n---\n\n");
      const revision = `skill-${(await sha256Hex(prompt)).slice(0, 12)}`;

      return { name, revision, prompt };
    })();
    revisionCache.set(name, cached);
  }

  return cached;
}
