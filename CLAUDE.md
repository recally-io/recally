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
pnpm typecheck        # all packages
pnpm test             # vitest
pnpm lint             # biome
cd apps/api && pnpm dev      # app worker (needs .dev.vars, D1 id)
cd apps/jobs && pnpm dev     # jobs worker
cd apps/web && pnpm dev      # vite dev server, proxies /api → :8787
```

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
