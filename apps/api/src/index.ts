import { newId, toErrorResponse } from "@recally/domain";
import { Hono } from "hono";
import { HTTPException } from "hono/http-exception";
import { requireAuth } from "./auth";
import { itemsRoutes } from "./routes/items";
import {
  adminRoutes,
  askRoutes,
  contentRoutes,
  jobsRoutes,
  publicRoutes,
  searchRoutes,
  sharesRoutes,
} from "./routes/misc";

const app = new Hono<{ Bindings: Env }>();

app.use("*", async (c, next) => {
  c.set("requestId" as never, newId() as never);
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

app.route("/api/v1", api);

app.onError((err, c) => {
  const requestId = newId();
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
