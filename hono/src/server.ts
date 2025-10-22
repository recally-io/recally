import { Hono } from "hono";
import { sql } from "kysely";

import { auth, database } from "./auth";

const forwardToAuth = async (request: Request) => auth.handler(request);

export const app = new Hono();

app.get("/healthz", (c) => c.json({ status: "ok" }));

app.get("/readyz", async (c) => {
  try {
    await sql`select 1`.execute(database);
    return c.json({ status: "ready" });
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unknown error";
    return c.json({ status: "error", message }, 503);
  }
});

app.all("/api/auth", (c) => forwardToAuth(c.req.raw));
app.all("/api/auth/*", (c) => forwardToAuth(c.req.raw));

export type AppType = typeof app;

export default app;
