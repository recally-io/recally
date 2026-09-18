import { askRequestSchema } from "@recally/contracts";
import { AppError } from "@recally/domain";
import { Hono } from "hono";
import { authFor } from "../auth";
import type { Env } from "../env";

export const askRoutes = new Hono<{ Bindings: Env }>().post("/", async (c) => {
  authFor(c, "search:read");
  askRequestSchema.parse(await c.req.json());
  throw new AppError("not_implemented", "ask lands with M4; search is available");
});
