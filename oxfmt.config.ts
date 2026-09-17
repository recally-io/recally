import { defineConfig } from "oxfmt";

// Oxfmt is the formatter; Oxlint is the linter (see oxlint.config.ts).
// The style options below are pinned explicitly rather than left to defaults so
// a future Oxfmt default change cannot silently reformat the repo.
export default defineConfig({
  ignorePatterns: [
    // Vendored plugin and installed agent tooling: upstream-owned files.
    "tools/oxlint/anti-slop/**",
    ".agent/**",
    ".agents/**",
    ".claude/**",
    ".codex/**",
    ".continue/**",
    ".cursor/**",
    ".gemini/**",
    ".opencode/**",
    ".pi/**",
    ".roo/**",
    ".windsurf/**",
    // Build output, generated types, and local state.
    ".alchemy/**",
    "**/dist/**",
    "**/worker-configuration.d.ts",
    // Long-form design docs and the Tailwind entry stylesheet are kept as
    // written: reformatting them churns review diffs without benefit.
    "docs/plan.md",
    "apps/web/src/index.css",
    // Design mockups are artifacts, not source.
    "docs/design/**",
    // Skill prompts are hashed into skill_revision, so a formatting-only edit
    // would invalidate in-flight runs and every recorded revision.
    "**/SKILL.md",
  ],
  printWidth: 100,
  tabWidth: 2,
  useTabs: false,
  semi: true,
  singleQuote: false,
  trailingComma: "all",
  bracketSpacing: true,
  insertFinalNewline: true,
  endOfLine: "lf",
});
