import { defineConfig } from "oxlint";

// anti-slop (vendored at tools/oxlint/anti-slop) owns the opinionated
// TypeScript and Effect rules. Oxfmt is the formatter (`pnpm format`); the
// anti-slop blank-line layout rule is what keeps spacing consistent in lint.
//
// Rules are vendored, not dependency-managed: edit them under
// tools/oxlint/anti-slop/ to match this repo's standards. See UPSTREAM.md
// there for provenance before pulling upstream changes.
export default defineConfig({
  ignorePatterns: [
    // Vendored plugin source, plus installed agent tooling (not our code).
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
    // Build output and generated types.
    ".alchemy/**",
    "**/dist/**",
    "**/worker-configuration.d.ts",
  ],
  jsPlugins: [
    { name: "anti-slop", specifier: "./tools/oxlint/anti-slop/index.ts" },
    {
      name: "anti-slop-effect",
      specifier: "./tools/oxlint/anti-slop/effect/index.ts",
    },
  ],
  rules: {
    // Native companion to anti-slop/no-reduce-accumulator-copy.
    "oxc/no-accumulating-spread": "error",

    "anti-slop/no-array-filter-map": "error",
    "anti-slop/no-reduce-accumulator-copy": "error",
    "anti-slop/no-chained-type-assertions": "error",
    "anti-slop/no-conditional-empty-object-spread": "error",
    "anti-slop/no-known-value-widening": "error",
    "anti-slop/no-module-mocking": "error",
    "anti-slop/no-object-parameters": "error",
    "anti-slop/no-reflect-apply": "error",
    "anti-slop/no-reflect-get": "error",
    "anti-slop/no-runtime-typeof": "error",
    "anti-slop/no-shape-in-symbol-names": "error",
    "anti-slop/no-unknown-parameters": "error",
    "anti-slop/no-unknown-returns": "error",
    "anti-slop/no-unknown-type-aliases": "error",
    "anti-slop/no-unsafe-dictionary-type": "error",
    "anti-slop/no-widen-then-assert": "error",
    "anti-slop/require-readable-spacing": "error",
    "anti-slop/require-safety-comment-for-type-assertion": "error",

    // Effect group: this repo depends on effect directly (root manifest).
    "anti-slop-effect/no-manual-effect-error-tag": "error",
    "anti-slop-effect/no-manual-tag-comparison": "error",
    "anti-slop-effect/no-manual-tagged-construction": "error",
    "anti-slop-effect/no-service-constructor-imports": "error",
    "anti-slop-effect/prefer-effect-match": "error",
  },
});
