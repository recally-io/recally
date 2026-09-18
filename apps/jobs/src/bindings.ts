import { PiRuntime } from "@recally/agent-runtime";
import { cfStreamFn, resolveCfModel } from "@recally/platform-cloudflare";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Layer from "effect/Layer";
import { definedModelVars, type ModelVars } from "../../../infra/model-env";
import type { JobsEnv } from "./env";

export const JobsBaseLayer = Layer.mergeAll(
  Cloudflare.D1.QueryDatabaseBinding,
  Cloudflare.R2.ReadWriteBucketBinding,
  Cloudflare.Workers.AIBinding,
);

export function makeJobsEnv(
  db: JobsEnv["DB"],
  bucket: JobsEnv["ARCHIVE_BUCKET"],
  ai: JobsEnv["AI"],
  modelVars: ModelVars,
): JobsEnv {
  return { DB: db, ARCHIVE_BUCKET: bucket, AI: ai, ...definedModelVars(modelVars) };
}

export function createPiRuntime(env: JobsEnv, runtimeVersion?: string): PiRuntime {
  const modelEnv = {
    AI: env.AI as never,
    AI_GATEWAY_NAME: env.AI_GATEWAY_NAME,
    CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
  };

  const options: ConstructorParameters<typeof PiRuntime>[0] = {
    streamFn: cfStreamFn(modelEnv),
    resolveModel: (id) => resolveCfModel(id, modelEnv),
  };

  if (runtimeVersion !== undefined) options.runtimeVersion = runtimeVersion;

  return new PiRuntime(options);
}
