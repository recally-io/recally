import { betterAuth } from "better-auth";
import { jwt } from "better-auth/plugins/jwt";
import { Kysely, PostgresDialect } from "kysely";
import { Pool } from "pg";

interface Env {
  authSecret: string;
  databaseUrl: string;
  databaseSsl: "require" | "disable" | "prefer";
  databasePoolMax: number;
  baseUrl: string;
  basePath: string;
  trustedOrigins: string[];
  jwtIssuer?: string;
  jwtAudience?: string;
  jwtExpiration?: string;
  jwksRemoteUrl?: string;
}

const env: Env = (() => {
  const requireEnv = (key: string): string => {
    const value = process.env[key];
    if (!value) {
      throw new Error(`Missing required environment variable: ${key}`);
    }
    return value;
  };

  const optionalEnv = (key: string): string | undefined => {
    const value = process.env[key];
    return value && value.trim().length > 0 ? value.trim() : undefined;
  };

  const parseNumber = (value: string | undefined, fallback: number): number => {
    if (!value) return fallback;
    const parsed = Number.parseInt(value, 10);
    return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
  };

  const sslMode = (optionalEnv("BETTER_AUTH_DATABASE_SSL" ) ?? "require").toLowerCase();
  const normalizedSsl = sslMode === "disable" ? "disable" : sslMode === "prefer" ? "prefer" : "require";

  const trustedOriginsRaw = optionalEnv("BETTER_AUTH_TRUSTED_ORIGINS");
  const trustedOrigins = trustedOriginsRaw
    ? trustedOriginsRaw.split(",").map((origin) => origin.trim()).filter((origin) => origin.length > 0)
    : [];

  return {
    authSecret: requireEnv("BETTER_AUTH_SECRET"),
    databaseUrl: requireEnv("BETTER_AUTH_DATABASE_URL"),
    databaseSsl: normalizedSsl,
    databasePoolMax: parseNumber(optionalEnv("BETTER_AUTH_DATABASE_POOL_MAX"), 10),
    baseUrl: requireEnv("BETTER_AUTH_BASE_URL"),
    basePath: optionalEnv("BETTER_AUTH_BASE_PATH") ?? "/api/auth",
    trustedOrigins,
    jwtIssuer: optionalEnv("BETTER_AUTH_JWT_ISSUER"),
    jwtAudience: optionalEnv("BETTER_AUTH_JWT_AUDIENCE"),
    jwtExpiration: optionalEnv("BETTER_AUTH_JWT_EXPIRATION"),
    jwksRemoteUrl: optionalEnv("BETTER_AUTH_JWKS_REMOTE_URL"),
  } satisfies Env;
})();

const baseOrigin = (() => {
  try {
    return new URL(env.baseUrl).origin;
  } catch {
    return undefined;
  }
})();

if (baseOrigin && !env.trustedOrigins.includes(baseOrigin)) {
  env.trustedOrigins.push(baseOrigin);
}

type AuthDatabase = Record<string, never>;

const createPool = () =>
  new Pool({
    connectionString: env.databaseUrl,
    max: env.databasePoolMax,
    ssl:
      env.databaseSsl === "disable"
        ? undefined
        : env.databaseSsl === "prefer"
          ? { rejectUnauthorized: false }
          : { rejectUnauthorized: true },
  });

const createDatabase = (pool: Pool) =>
  new Kysely<AuthDatabase>({
    dialect: new PostgresDialect({ pool }),
  });

declare global {
  // eslint-disable-next-line no-var
  var __recallyAuthPool: Pool | undefined;
  // eslint-disable-next-line no-var
  var __recallyAuthDb: Kysely<AuthDatabase> | undefined;
  // eslint-disable-next-line no-var
  var __recallyAuthInstance: ReturnType<typeof betterAuth> | undefined;
}

const pool = globalThis.__recallyAuthPool ?? createPool();
if (!globalThis.__recallyAuthPool) {
  globalThis.__recallyAuthPool = pool;
}

export const database = globalThis.__recallyAuthDb ?? createDatabase(pool);
if (!globalThis.__recallyAuthDb) {
  globalThis.__recallyAuthDb = database;
}

const plugins = [
  jwt({
    jwks: env.jwksRemoteUrl ? { remoteUrl: env.jwksRemoteUrl } : undefined,
    jwt: {
      issuer: env.jwtIssuer ?? env.baseUrl,
      audience: env.jwtAudience ?? env.baseUrl,
      expirationTime: env.jwtExpiration ?? "15m",
    },
  }),
];

export const auth =
  globalThis.__recallyAuthInstance ??
  betterAuth({
    secret: env.authSecret,
    baseURL: env.baseUrl,
    basePath: env.basePath,
    trustedOrigins: env.trustedOrigins,
    database,
    plugins,
    emailAndPassword: {
      enabled: true,
    },
  });

if (!globalThis.__recallyAuthInstance) {
  globalThis.__recallyAuthInstance = auth;
}

export type AuthInstance = typeof auth;
