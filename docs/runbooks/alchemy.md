# Runbook: Alchemy deploy & dev (Cloudflare)

The infra is a single root Stack (`alchemy.run.ts`) powered by alchemy
(Infrastructure-as-Effects, Effect v4 rc). Resources are declared once in
`infra/resources.ts` and bound into both workers. Wrangler is gone from the
repo and from the account: `alchemy` owns provisioning, bundling, migrations,
and cron.

## One-time setup

1. Auth alchemy with Cloudflare (OAuth in the browser, or paste an API token):

   ```
   pnpm exec alchemy profile edit --add Cloudflare
   ```

   Credentials land under `~/.alchemy/` (profile `default`). No
   `CLOUDFLARE_API_TOKEN` env var and no `wrangler login` needed.

2. Create `.env` from `.env.example` and fill `DEV_LIBRARY_ID` (and model
   overrides if desired). Values present at deploy time are bound to the
   workers as secrets; unset keys fall back to the code defaults in
   `packages/ai/src/models.ts`.

   `DEV_LIBRARY_ID` only opens the anonymous auth path if a `libraries` row
   with that id exists — on a fresh database every `/api/v1` request answers
   `needs_login` until one is provisioned (Access login, or an insert).

## Deploy

```
pnpm deploy                  # web build + alchemy deploy
pnpm exec alchemy deploy --yes
```

Outputs printed at the end: `apiUrl`, `jobsUrl` (workers.dev).

Every resource name — the two worker scripts included — is derived from the
app, the logical id, and the stage. The deploy stage defaults to `live_$USER`
and `alchemy dev` to `dev_$USER`; override with `--stage <name>` or
`$ALCHEMY_STAGE`. Two stages therefore never share a D1, bucket, index, or
script, and no `--adopt` is ever needed.

D1 migrations in `migrations/` are applied on create and recorded in alchemy's
`__alchemy_migrations` table. There is no separate apply or status command.

## Routine plan / destroy

```
pnpm plan                    # diff without applying
pnpm destroy                 # delete every resource in the stage
pnpm exec alchemy plan --stage prod
```

## Local development

```
pnpm build   # once, so apps/web/dist exists for the asset layer
pnpm dev     # alchemy dev
```

- Workers run locally in workerd with hot reload; D1/R2 are local simulators
  (fresh and empty — not the deployed resources). Pin a resource live in dev
  with `Alchemy.remote()` if you need real data.
- Web: keep `pnpm --filter @recally/web dev` for SPA HMR and point
  `vite.config.ts`'s proxy at the api worker's local URL that `alchemy dev`
  prints.
- Cron triggers do not fire in dev; outbox dispatch is exercised by the
  integration test or on the deployed stage.

## Integration test (live stack)

```
ALCHEMY_INTEG=1 pnpm test tests/integration/stack.test.ts
```

Deploys an isolated `alchemy-integ` stage, asserts SPA fallback + share 404 +
the auth gate, then destroys (keep with `NO_DESTROY=1`).

## Notes / gotchas

- Effect is pinned to `4.0.0-rc.112` via `pnpm-workspace.yaml` overrides:
  alchemy 2.0.0-beta.77 is built against its API (e.g. `Config.string`), which
  newer rcs renamed. Bump both together.
- `.md` files are imported with `?raw` (packages/skills). The worker bundle
  inlines them via rolldown's raw plugin; the alchemy CLI's stack loader never
  evaluates them because `loadSkill` imports them dynamically (runtime only).
- `BACKUP_BUCKET` is managed by the Stack but not bound into any worker yet —
  the export milestone (§15.2) binds it.
- The wrangler-era resources (`recally-api`/`recally-jobs`/`recally-api-prod`
  scripts, D1 `recally-dev`, both `recally-*-dev` buckets, Vectorize
  `recally-chunks-dev`) were deleted on 2026-09-17; their data was not carried
  over. In particular `recally.io` and `www.recally.io` no longer have a route
  bound, so the public hostname is dark until a stage binds it (see below).
- Binding `recally.io` needs a decision first: alchemy stages do not share a
  data plane, so the route belongs to whichever stage serves real data. Add it
  via the Worker's `routes` prop (zone routes, like the wrangler `env.prod`
  config) once that is settled.
