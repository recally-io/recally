import type { TokenScope } from "@recally/contracts";
import { AppError, sha256Hex } from "@recally/domain";
import { createLibrary, getLibrary, getLibraryByAccessSub } from "@recally/storage";
import { createMiddleware } from "hono/factory";

// Identity resolution order: scoped API token > Cloudflare Access JWT > dev
// bypass. Whatever resolves must still map to a concrete library row — the
// client can never pass library_id in a request body (plan §12.1).

export interface AuthContext {
  libraryId: string;
  actor: string; // token id or access subject
  email: string | null;
  scopes: TokenScope[] | "access"; // "access" = full UI session
  via: "token" | "access" | "dev";
}

declare module "hono" {
  interface ContextVariableMap {
    auth: AuthContext;
  }
}

interface AuthEnv {
  DB: D1Database;
  ENVIRONMENT?: string;
  DEV_LIBRARY_ID?: string;
  ACCESS_TEAM_NAME?: string;
  ACCESS_AUD?: string;
}

// --- API token ---

async function authByToken(env: AuthEnv, raw: string): Promise<AuthContext | null> {
  const hash = await sha256Hex(raw);
  const row = await env.DB.prepare(
    `SELECT t.*, l.id AS lib FROM api_tokens t
     WHERE t.token_hash = ? AND t.revoked_at IS NULL`,
  )
    .bind(hash)
    .first<{
      id: string;
      library_id: string;
      scopes: string;
      expires_at: string | null;
    }>();
  if (!row) return null;
  if (row.expires_at && row.expires_at < new Date().toISOString()) return null;
  const lib = await getLibrary(env.DB, row.library_id);
  if (!lib) return null;
  await env.DB.prepare("UPDATE api_tokens SET last_used_at = ? WHERE id = ?")
    .bind(new Date().toISOString(), row.id)
    .run();
  return {
    libraryId: row.library_id,
    actor: row.id,
    email: null,
    scopes: JSON.parse(row.scopes) as TokenScope[],
    via: "token",
  };
}

// --- Cloudflare Access JWT (RS256, JWKS) ---

const jwksCache = new Map<string, { keys: Map<string, CryptoKey>; at: number }>();

async function accessKeys(team: string): Promise<Map<string, CryptoKey>> {
  const cached = jwksCache.get(team);
  if (cached && Date.now() - cached.at < 5 * 60 * 1000) return cached.keys;
  const res = await fetch(`https://${team}.cloudflareaccess.com/cdn-cgi/access/certs`);
  if (!res.ok)
    throw new AppError("needs_login", "access certs unavailable", {
      retryable: true,
    });
  const jwks = (await res.json()) as {
    keys: Array<{ kid: string } & JsonWebKey>;
  };
  const keys = new Map<string, CryptoKey>();
  for (const jwk of jwks.keys) {
    keys.set(
      jwk.kid,
      await crypto.subtle.importKey(
        "jwk",
        jwk,
        { name: "RSASSA-PKCS1-v1_5", hash: "SHA-256" },
        false,
        ["verify"],
      ),
    );
  }
  jwksCache.set(team, { keys, at: Date.now() });
  return keys;
}

function b64url(s: string): Uint8Array {
  const pad = "=".repeat((4 - (s.length % 4)) % 4);
  const bin = atob(s.replace(/-/g, "+").replace(/_/g, "/") + pad);
  return Uint8Array.from(bin, (c) => c.charCodeAt(0));
}

async function verifyAccessJwt(
  env: AuthEnv,
  jwt: string,
): Promise<{ sub: string; email: string | null } | null> {
  if (!env.ACCESS_TEAM_NAME || !env.ACCESS_AUD) return null;
  const [h, p, s] = jwt.split(".");
  if (!h || !p || !s) return null;
  const header = JSON.parse(new TextDecoder().decode(b64url(h))) as {
    kid?: string;
    alg?: string;
  };
  if (header.alg !== "RS256" || !header.kid) return null;

  const keys = await accessKeys(env.ACCESS_TEAM_NAME);
  const key = keys.get(header.kid);
  if (!key) return null;

  const ok = await crypto.subtle.verify(
    "RSASSA-PKCS1-v1_5",
    key,
    b64url(s) as unknown as ArrayBuffer,
    new TextEncoder().encode(`${h}.${p}`) as unknown as ArrayBuffer,
  );
  if (!ok) return null;

  const claims = JSON.parse(new TextDecoder().decode(b64url(p))) as {
    sub?: string;
    email?: string;
    iss?: string;
    aud?: string | string[];
    exp?: number;
  };
  if (!claims.sub || !claims.exp || claims.exp * 1000 < Date.now()) return null;
  if (claims.iss !== `https://${env.ACCESS_TEAM_NAME}.cloudflareaccess.com`) return null;
  const aud = Array.isArray(claims.aud) ? claims.aud : [claims.aud];
  if (!aud.includes(env.ACCESS_AUD)) return null;
  return { sub: claims.sub, email: claims.email ?? null };
}

async function authByAccess(env: AuthEnv, req: Request): Promise<AuthContext | null> {
  const jwt =
    req.headers.get("cf-access-jwt-assertion") ??
    req.headers.get("cookie")?.match(/CF_Authorization=([^;]+)/)?.[1];
  if (!jwt) return null;
  const verified = await verifyAccessJwt(env, jwt);
  if (!verified) return null;
  // First verified Access login = signup: provision a library on the spot.
  const lib =
    (await getLibraryByAccessSub(env.DB, verified.sub)) ??
    (await createLibrary(env.DB, {
      name: verified.email ?? verified.sub,
      accessSub: verified.sub,
    }));
  return {
    libraryId: lib.id,
    actor: verified.sub,
    email: verified.email,
    scopes: "access",
    via: "access",
  };
}

// --- middleware ---

export const requireAuth = createMiddleware<{ Bindings: AuthEnv }>(async (c, next) => {
  const authz = c.req.header("authorization");
  let auth: AuthContext | null = null;

  if (authz?.startsWith("Bearer ")) {
    auth = await authByToken(c.env, authz.slice(7));
  }
  auth ??= await authByAccess(c.env, c.req.raw);

  // Dev bypass is inert on the production deployment — without it, recally.io
  // would hand the dev library to anyone who hits the API unauthenticated.
  if (!auth && c.env.DEV_LIBRARY_ID && c.env.ENVIRONMENT !== "production") {
    const lib = await getLibrary(c.env.DB, c.env.DEV_LIBRARY_ID);
    if (lib) {
      auth = {
        libraryId: lib.id,
        actor: "dev",
        email: null,
        scopes: "access",
        via: "dev",
      };
    }
  }

  if (!auth) throw new AppError("needs_login", "authentication required");
  c.set("auth", auth);
  await next();
});

export function requireScope(auth: AuthContext, scope: TokenScope): void {
  if (auth.scopes !== "access" && !auth.scopes.includes(scope)) {
    throw new AppError("forbidden", `missing scope ${scope}`);
  }
}
