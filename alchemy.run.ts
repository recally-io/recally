import * as Alchemy from "alchemy";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import Api from "./apps/api/src/worker";
import Jobs from "./apps/jobs/src/worker";
import { ArchiveBucket, Database, VectorIndex } from "./infra/resources";

// Recally — personal reading archive on Cloudflare (plan §3).
// Single Stack: both workers + shared D1/R2/Vectorize in one plan.
//
// Every resource (including both worker scripts) is named from the stage, so
// `alchemy deploy --stage prod` provisions a physically separate stack. The
// stage defaults to `live_$USER`; pin it with `--stage` or `$ALCHEMY_STAGE`.
export default Alchemy.Stack(
  "Recally",
  {
    providers: Cloudflare.providers(),
    state: Cloudflare.state(),
  },
  Effect.gen(function* () {
    // Shared data plane — declared once, bound into both workers.
    yield* Database;
    yield* ArchiveBucket;
    yield* VectorIndex;

    const api = yield* Api;
    const jobs = yield* Jobs;

    return {
      apiUrl: api.url.as<string>(),
      jobsUrl: jobs.url.as<string>(),
    };
  }),
);
