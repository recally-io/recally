import * as Cloudflare from "alchemy/Cloudflare";
import * as Config from "effect/Config";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import * as Option from "effect/Option";
import * as HttpServerRequest from "effect/unstable/http/HttpServerRequest";
import * as HttpServerResponse from "effect/unstable/http/HttpServerResponse";
import { ArchiveBucket, Database } from "../../../infra/resources";
import type { Env } from "./env";
import app from "./index";

// App Worker (plan §3.2): Hono API + the SPA on one origin. The Construction
// phase registers the bindings and resolves the Config values (each yield*
// binds them onto the worker); the returned fetch resolves the native handles
// (Runtime phase) and hands them to the Hono app via its Env, so the route
// layer stays platform-neutral.
export default class Api extends Cloudflare.Worker<Api>()(
  "Api",
  {
    name: "recally-api",
    main: import.meta.url,
    compatibility: { date: "2026-09-11", flags: ["nodejs_compat"] },
    // Assets are served first; requests matching no asset invoke the Worker,
    // whose catch-all forwards to the ASSETS binding (SPA fallback for deep
    // links, Hono routes for /api/v1 and /s).
    assets: {
      directory: "apps/web/dist",
      notFoundHandling: "single-page-application",
    },
  },
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);
    const bucketClient = yield* Cloudflare.R2.ReadWriteBucket(ArchiveBucket);
    const devLibraryId = yield* Config.string("DEV_LIBRARY_ID").pipe(
      Config.option,
      Effect.map((o) => Option.getOrUndefined(o)),
    );
    const captureModel = yield* Config.string("CAPTURE_MODEL").pipe(
      Config.option,
      Effect.map((o) => Option.getOrUndefined(o)),
    );

    return {
      fetch: Effect.gen(function* () {
        const request = yield* HttpServerRequest.HttpServerRequest;
        const [db, r2, workerEnv] = yield* Effect.all([
          dbClient.raw,
          bucketClient.raw,
          Cloudflare.Workers.WorkerEnvironment,
        ]);
        const env: Env = {
          DB: db,
          ARCHIVE_BUCKET: r2,
          ASSETS: workerEnv.ASSETS as Env["ASSETS"],
          ...(devLibraryId !== undefined ? { DEV_LIBRARY_ID: devLibraryId } : {}),
          ...(captureModel !== undefined ? { CAPTURE_MODEL: captureModel } : {}),
        };
        const response = yield* Effect.promise(() =>
          Promise.resolve(app.fetch(request.source as unknown as Request, env)),
        );
        return HttpServerResponse.fromWeb(response);
      }),
    };
  }).pipe(
    Effect.provide(
      Layer.mergeAll(Cloudflare.D1.QueryDatabaseBinding, Cloudflare.R2.ReadWriteBucketBinding),
    ),
  ),
) {}
