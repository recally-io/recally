# Runbook: Alchemy deploy & dev (Cloudflare)

The infra is a single root Stack (`alchemy.run.ts`) powered by alchemy
(Infrastructure-as-Effects, Effect v4 rc). Resources are declared once in
`infra/resources.ts` and bound into both workers. Wrangler configs are gone;
`alchemy` owns provisioning, bundling, migrations, and cron.

## One-time setup

1. Auth alchemy with Cloudflare (OAuth in the browser, or paste an API token):

   ```
   pnpm exec alchemy profile edit --add Cloudflare
   ```

   Credentials land in `~/.alchemy/profiles.json` (profile `default`).
   No `CLOUDFLARE_API_TOKEN` env var and no `wrangler login` needed.

2. Create `.env` from `.env.example` and fill `DEV_LIBRARY_ID` (and model
   overrides if desired). Values present at deploy time are bound to the
   workers as secrets; unset keys fall back to the code defaults in
   `packages/ai/src/models.ts`.

## First deploy (adopts existing resources)

The dev resources created under wrangler (D1 `recally-dev`, R2
`recally-archive-dev` / `recally-backup-dev`, Vectorize `recally-chunks-dev`,
both worker scripts) are adopted — not recreated — by the first deploy:

```
pnpm build        # web SPA → apps/web/dist (worker assets)
pnpm exec alchemy deploy --adopt --yes
```

- D1 migration history applied by `wrangler d1 migrations apply` is copied
  into alchemy's `__alchemy_migrations` bookkeeping automatically.
- `--adopt` is only needed on the first deploy (foreign ownership tags).
- Outputs printed at the end: `apiUrl`, `jobsUrl` (workers.dev).

## Routine deploy / plan / destroy

```
pnpm deploy          # web build + alchemy deploy (dev stage)
pnpm plan            # diff without applying
pnpm exec alchemy destroy   # delete every resource in the stage
```

Stages isolate everything: `alchemy deploy --stage prod` is a physically
separate copy. `dev`/`prod` sharing the same physical resources (the old
wrangler setup) is intentionally NOT reproduced; when prod launches, make the
resource names stage-conditional in `infra/resources.ts`.

## Local development

```
pnpm build   # once, so apps/web/dist exists for the asset layer
pnpm dev     # alchemy dev
```

- Workers run locally in workerd with hot reload; D1/R2 are local simulators
  (fresh and empty — not the adopted remote resources). Pin a resource live in
  dev with `Alchemy.remote()` if you need real data.
- Web: keep `pnpm --filter @recally/web dev` for SPA HMR and point
  `vite.config.ts`'s proxy at the api worker's local URL that `alchemy dev`
  prints.
- Cron triggers do not fire in dev; outbox dispatch is exercised by the
  integration test or on the deployed stage.

## Integration test (live stack)

```
ALCHEMY_INTEG=1 pnpm test tests/integration/stack.test.ts
```

Deploys an isolated `alchemy-integ` stage, asserts /health + share 404 + SPA
fallback, then destroys (keep with `NO_DESTROY=1`).

## Notes / gotchas

- Effect is pinned to `4.0.0-rc.112` via `pnpm-workspace.yaml` overrides:
  alchemy 2.0.0-beta.77 is built against its API (e.g. `Config.string`), which
  newer rcs renamed. Bump both together.
- `.md` files are imported with `?raw` (packages/skills). The worker bundle
  inlines them via rolldown's raw plugin; the alchemy CLI's stack loader never
  evaluates them because `loadSkill` imports them dynamically (runtime only).
- `BACKUP_BUCKET` is managed by the Stack but not bound into any worker yet —
  the export milestone (§15.2) binds it.
- Old wrangler workflow names (`recally-ingest`, …) are superseded by
  alchemy-managed names derived from the host worker; in-flight instances from
  the wrangler era are not migrated (dev-only concern).
