# CLAUDE.md

Guidance for agents working in this repository. The authoritative design is `docs/plan.md` (§ numbers below refer to it).

## What this is

Recally — a personal reading archive running entirely on Cloudflare. Users save URLs; a **Pi CaptureAgent** (skills + tools, not a fixed pipeline) decides per-URL whether to use `web_fetch`, a site adapter, or a Browser Run session; an **AnalysisAgent** then summarizes committed snapshots. See plan §1.2 for the governing principles — they are requirements, not suggestions:

- Skills encode strategy; tools expose capability; the agent owns orchestration.
- Workflow is the durable supervisor, never the planner.
- Model-generated text can never become the archived original (§7.1).
- Every tool result is `UNTRUSTED_SOURCE` (§6.7, §14.3).

## Layout

```
apps/
  api/    App Worker (Hono). Auth, items/jobs/search/shares APIs, static web assets. No BROWSER binding.
  jobs/   Jobs Worker. Pi runtime host, 6 Workflow classes, cron outbox dispatch. Owns BROWSER binding.
  web/    React 19 + Vite + TanStack Router + Tailwind 4.
packages/
  domain/             ids, hashing, URL normalization, statuses, errors, budgets, versions
  contracts/          zod schemas: API payloads + tool inputs
  agent-runtime/      Pi wrapper (ToolSpec → AgentTool), per-turn durable loop contract
  skills/             capture/SKILL.md + references, analyze/SKILL.md, revision loader
  tools/              web-fetch/, site-registry/, browser/, archive/, read-source/, analysis/
  site-adapters/      github, hackernews (+registry). Compiled-in only, never remote eval.
  platform-cloudflare/  bindings: R2EvidenceStore, D1R2RunStore, ArchiveService, VectorIndex,
                        BrowserSession, workflow dispatch, Pi/Workers-AI model wiring
  storage/            D1 row types + queries + R2 key layout
  ai/                 model roles, WorkersAIProvider (generate/embed), prompts, schemas
  search/             CJK tokenizer, chunking, RRF, vector ids
  observability/      usage reservation/settlement, audit events
migrations/           D1 SQL migrations (manual, sequential numbering)
evals/                capture fixtures/golden/strategy cases, analysis, search
tests/                integration, fault-injection, security
docs/                 plan.md, decisions/, runbooks/
```

## Commands

```bash
pnpm install
pnpm typecheck        # root stack program + all packages
pnpm test             # vitest (ALCHEMY_INTEG=1 adds live-stack tests)
pnpm lint             # oxlint, including the vendored anti-slop rules
pnpm lint:fix         # oxlint autofix
pnpm format           # oxfmt
pnpm format:check     # oxfmt without writing
pnpm build            # build web SPA (required before deploy/dev)
pnpm deploy           # web build + alchemy deploy (stage defaults to live_$USER)
pnpm dev              # alchemy dev — workers in workerd, web via Vite, hot reload
pnpm plan             # preview stack changes without applying
alchemy profile edit --add Cloudflare   # one-time Cloudflare auth (OAuth)
```

Infra lives in a single root Stack (`alchemy.run.ts`) + `infra/`; both workers
are Effectful Constructors (`apps/*/src/worker.ts`). There is no wrangler
config; worker types come from the worker files, not `wrangler types`. Every
resource name (worker scripts included) is stage-derived, so `--stage prod` is
a physically separate stack — never pin a physical name in `infra/`.

## Code quality gate

`pnpm lint` is [Oxlint](https://oxc.rs/docs/guide/usage/linter), configured in
`oxlint.config.ts`. `pnpm format` is Oxfmt (`oxfmt.config.ts`). Biome is gone.
[anti-slop](https://github.com/dmmulroy/anti-slop) is **vendored** at
`tools/oxlint/anti-slop/` (see its `UPSTREAM.md`) and owns the opinionated
TypeScript + Effect rules: no `unknown` at boundaries, no type-assertion
laundering, no chained `as`, a `SAFETY:` comment on every remaining non-const
assertion, and a blank-line layout rule.

Treat a new anti-slop finding as a design signal, not lint noise. Fix the
underlying type or boundary; a `SAFETY:` comment is only correct when the
invariant really is checked nearby. The rules are ours to edit — change them in
`tools/oxlint/anti-slop/` and record the deviation in its `UPSTREAM.md` rather
than weakening severity in `oxlint.config.ts`.

## Hard rules (from the plan, enforced in review)

- Tools never chain fetch→browser internally; `web_fetch` is one GET + extraction only (§6.2).
- `propose_archive` is the only publish path; ArchiveService validates + commits (§5.10, §8.2).
- Revision fencing: a run replays its pinned `agent_runtime_version + skill_revision + toolset_version` (§9.6). Bumping a skill never changes an in-flight run.
- Item/`generation`/delete checks re-run before every external action and before commit (§9.7).
- Vectorize hits must pass D1 permission/deletion/version checks before reaching the user (§11.4).
- No tool ever gets: SQL, arbitrary JS eval, shell, cookie import, share/delete/budget authority (§6.8).
- Adapters and skills are platform-neutral; only `platform-cloudflare` imports bindings (§3.5, ADR-06).

## Known scaffold boundaries

- Pi on Workers is wired against real `@earendil-works/pi-agent-core` 0.85.1 APIs but M0-T02/T06/T07 verification (bundle, state restore, browser episode) is still open.
- `source_documents`/`assets` asset pipeline, needs_input child runs, usage settlement, backup/export — stubbed or absent until their milestone.
