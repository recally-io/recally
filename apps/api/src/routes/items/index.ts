import { Hono } from "hono";
import type { Env } from "../../env";
import { createRoutes } from "./create";
import { mutationRoutes } from "./mutations";
import { readRoutes } from "./read";

export const itemsRoutes = new Hono<{ Bindings: Env }>()
  .route("/", createRoutes)
  .route("/", readRoutes)
  .route("/", mutationRoutes);
