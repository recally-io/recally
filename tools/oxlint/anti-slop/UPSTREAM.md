# Vendored anti-slop plugins

Source: [dmmulroy/anti-slop](https://github.com/dmmulroy/anti-slop), commit
`c44ef22ca116d0ba62a3ff663a0bd13a3f3fa40b` (2026-09-10, merge of PR #36
"contrib/effect-tag-match-rules").

This directory is vendored code, not a dependency. It is ours to read, change,
and extend. Upstream ships no npm package; `oxlint` and `@oxlint/plugins` are
pinned to the same exact version in the root manifest so upgrades move together.

## Installed paths

- `index.ts` — generic rules, registered as plugin `anti-slop`.
- `effect/index.ts` — Effect rules, registered as plugin `anti-slop-effect`.
  Enabled because this repo depends on `effect` directly.
- `vendor/eslint-stylistic/` — the `padding-line-between-statements` port that
  backs `require-readable-spacing`; carries its own `UPSTREAM.md` and `LICENSE`.

Registered in `oxlint.config.ts` at the repo root. `pnpm lint` runs `oxlint`
alongside `biome check`; Biome remains the formatter.

## Verification

The copied files are byte-identical to upstream `src/` at the commit above,
excluding `*.test.ts` (not shipped in the skill assets), and to
`skills/install-anti-slop/assets/anti-slop/` in that same revision. Upstream's
`node scripts/sync-skill-assets.mjs --check` passes at that commit, so the
skill assets and `src/` agree and the recorded revision identifies the actual
copied bytes.

## Intentional deviations

None. Local policy lives in `oxlint.config.ts` (rule severity and enablement),
not in edited rule source. If a rule is changed here, record what changed and
why in this section.
