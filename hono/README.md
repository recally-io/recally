# Recally Hono Auth Service

This package hosts a lightweight [Hono](https://hono.dev/) application that exposes the `better-auth` REST API under `/api/auth` together with health probes. It is intended to run alongside the main Recally backend and share the same PostgreSQL database.

## Environment variables

The service is configured entirely through environment variables. All variables are prefixed with `BETTER_AUTH_` to avoid collisions with other services.

| Variable | Required | Description |
| --- | --- | --- |
| `BETTER_AUTH_SECRET` | ✅ | Primary secret used by `better-auth` to derive encryption keys. Rotate when compromised. |
| `BETTER_AUTH_DATABASE_URL` | ✅ | PostgreSQL connection string consumed by the Kysely driver. |
| `BETTER_AUTH_DATABASE_SSL` | ❌ | SSL mode for the PostgreSQL driver. Accepts `require`, `prefer`, or `disable` (default: `require`). |
| `BETTER_AUTH_DATABASE_POOL_MAX` | ❌ | Maximum pool size for the shared `pg` connection pool (default: `10`). |
| `BETTER_AUTH_BASE_URL` | ✅ | Absolute URL of this service (e.g. `https://auth.recally.ai`). Used for callback URLs and default JWT issuer/audience. |
| `BETTER_AUTH_BASE_PATH` | ❌ | Mount path for the auth router (default: `/api/auth`). |
| `BETTER_AUTH_TRUSTED_ORIGINS` | ❌ | Comma-separated list of additional origins allowed to call the auth endpoints. The base URL's origin is always appended automatically. |
| `BETTER_AUTH_JWT_ISSUER` | ❌ | Overrides the issuer claim for generated JWTs (defaults to `BETTER_AUTH_BASE_URL`). |
| `BETTER_AUTH_JWT_AUDIENCE` | ❌ | Overrides the audience claim for generated JWTs (defaults to `BETTER_AUTH_BASE_URL`). |
| `BETTER_AUTH_JWT_EXPIRATION` | ❌ | Custom JWT expiration string (e.g. `30m`, `12h`). Defaults to `15m`. |
| `BETTER_AUTH_JWKS_REMOTE_URL` | ❌ | Remote JWKS endpoint. When set, JWKS are fetched remotely instead of exposing `/jwks`. |

## Development

Install dependencies with Bun and start the development server:

```bash
bun install
bun run dev
```

The dev command runs `src/server.ts` with hot reload. Health endpoints are served at `/healthz` and `/readyz`, and all `better-auth` routes are proxied under `/api/auth/*`.

## Build artifacts

`bun run build` emits an ESM bundle into the `dist/` directory:

- `dist/server.js` — bundled Hono application entry
- `dist/server.js.map` — source map for debugging (if enabled by Bun)

The resulting build can be executed with Bun (`bun run dist/server.js`) or imported into another runtime that can execute Bun-compatible ESM output.
