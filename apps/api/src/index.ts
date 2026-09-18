import { newId, toErrorResponse } from "@recally/domain";
import { getLibrary } from "@recally/storage";
import { Hono } from "hono";
import { HTTPException } from "hono/http-exception";
import { requireAuth } from "./auth";
import type { Env } from "./env";
import { itemsRoutes } from "./routes/items";
import { adminRoutes } from "./routes/admin";
import { askRoutes } from "./routes/ask";
import { contentRoutes } from "./routes/content";
import { jobsRoutes } from "./routes/jobs";
import { publicRoutes } from "./routes/public-share";
import { searchRoutes } from "./routes/search";
import { sharesRoutes } from "./routes/shares";

const app = new Hono<{ Bindings: Env }>();

app.use("*", async (c, next) => {
  c.set("requestId", newId());
  await next();
});

// Public share surface is matched before auth; it does not expose private API.
app.route("/s", publicRoutes);

const api = new Hono<{ Bindings: Env }>();

api.use("*", requireAuth);

api.route("/items", itemsRoutes);

api.route("/jobs", jobsRoutes);

api.route("/search", searchRoutes);

api.route("/ask", askRoutes);

api.route("/shares", sharesRoutes);

api.route("/", contentRoutes);

api.route("/", adminRoutes);

api.get("/health", (c) => c.json({ ok: true }));

api.get("/me", async (c) => {
  const auth = c.get("auth");

  const lib = await getLibrary(c.env.DB, auth.libraryId);

  return c.json({
    library_id: auth.libraryId,
    library_name: lib?.name ?? "",
    actor: auth.actor,
    email: auth.email,
    via: auth.via,
  });
});

app.route("/api/v1", api);

app.onError((err, c) => {
  const requestId = c.get("requestId");

  if (err instanceof HTTPException) {
    const { status, body } = toErrorResponse(
      Object.assign(new Error(err.message), { name: "HTTPException" }),
      requestId,
    );

    return c.json(body, status as 500);
  }

  const { status, body } = toErrorResponse(err, requestId);

  return c.json(body, status as 500);
});

// Everything else falls through to static assets (the React app).
app.all("*", async (c) => c.env.ASSETS.fetch(c.req.raw));

export default app;
