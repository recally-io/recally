import type { TokenScope } from "@recally/contracts";
import { AppError, sha256Hex } from "@recally/domain";
import { getLibrary } from "@recally/storage";
import { createMiddleware } from "hono/factory";

// Identity resolution order: scoped API token > default library.
// Auth is intentionally minimal for now — this is a single-user instance.
// Whatever resolves must still map to a concrete library row — the client can
// never pass library_id in a request body (plan §12.1).

export interface AuthContext {
  libraryId: string;
  actor: string; // token id or "default"
  email: string | null;
  scopes: TokenScope[] | "access"; // "access" = full UI session
  via: "token" | "default";
}

declare module "hono" {
  interface ContextVariableMap {
    auth: AuthContext;
  }
}

interface AuthEnv {
  DB: D1Database;
  DEV_LIBRARY_ID?: string;
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

// --- middleware ---

export const requireAuth = createMiddleware<{ Bindings: AuthEnv }>(async (c, next) => {
  const authz = c.req.header("authorization");
  let auth: AuthContext | null = null;

  if (authz?.startsWith("Bearer ")) {
    auth = await authByToken(c.env, authz.slice(7));
  }

  // No token (or an unrecognized one): fall back to the default library.
  if (!auth && c.env.DEV_LIBRARY_ID) {
    const lib = await getLibrary(c.env.DB, c.env.DEV_LIBRARY_ID);
    if (lib) {
      auth = {
        libraryId: lib.id,
        actor: "default",
        email: null,
        scopes: "access",
        via: "default",
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
