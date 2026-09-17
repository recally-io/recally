// Model-role strings and Workers AI transport knobs, shared by both workers.
//
// These are read with effect/Config during the Construction phase, which makes
// alchemy bind them onto the worker environment (docs: environments/secrets).
// Only keys that are actually set get bound; unset keys stay unset and the
// code-level defaults in @recally/ai apply — identical semantics to the old
// wrangler vars + .dev.vars setup.
//
// This module is imported by worker files in apps/*, so it must stay free of
// @recally/* imports (root-level files aren't workspace packages).

import type { ConfigError } from "effect/Config";
import * as Config from "effect/Config";
import * as Effect from "effect/Effect";
import * as Option from "effect/Option";

export interface ModelVars {
  CAPTURE_MODEL?: string | undefined;
  VERIFY_MODEL?: string | undefined;
  SUMMARY_MODEL?: string | undefined;
  ANSWER_MODEL?: string | undefined;
  EMBEDDING_MODEL?: string | undefined;
  AI_GATEWAY_NAME?: string | undefined;
  CLOUDFLARE_ACCOUNT_ID?: string | undefined;
}

const readOptional = (key: string) =>
  Config.string(key).pipe(
    Config.option,
    Effect.map((o) => Option.getOrUndefined(o)),
  );

// Spread into a worker/workflow env literal: only keys that are actually set
// become properties (exactOptionalPropertyTypes-safe).
export function definedModelVars(vars: ModelVars): Partial<ModelVars> {
  const out: Partial<ModelVars> = {};

  if (vars.CAPTURE_MODEL !== undefined) out.CAPTURE_MODEL = vars.CAPTURE_MODEL;

  if (vars.VERIFY_MODEL !== undefined) out.VERIFY_MODEL = vars.VERIFY_MODEL;

  if (vars.SUMMARY_MODEL !== undefined) out.SUMMARY_MODEL = vars.SUMMARY_MODEL;

  if (vars.ANSWER_MODEL !== undefined) out.ANSWER_MODEL = vars.ANSWER_MODEL;

  if (vars.EMBEDDING_MODEL !== undefined) out.EMBEDDING_MODEL = vars.EMBEDDING_MODEL;

  if (vars.AI_GATEWAY_NAME !== undefined) out.AI_GATEWAY_NAME = vars.AI_GATEWAY_NAME;

  if (vars.CLOUDFLARE_ACCOUNT_ID !== undefined) {
    out.CLOUDFLARE_ACCOUNT_ID = vars.CLOUDFLARE_ACCOUNT_ID;
  }

  return out;
}

export const ModelVarsConfig: Effect.Effect<ModelVars, ConfigError> = Effect.gen(function* () {
  return {
    CAPTURE_MODEL: yield* readOptional("CAPTURE_MODEL"),
    VERIFY_MODEL: yield* readOptional("VERIFY_MODEL"),
    SUMMARY_MODEL: yield* readOptional("SUMMARY_MODEL"),
    ANSWER_MODEL: yield* readOptional("ANSWER_MODEL"),
    EMBEDDING_MODEL: yield* readOptional("EMBEDDING_MODEL"),
    AI_GATEWAY_NAME: yield* readOptional("AI_GATEWAY_NAME"),
    CLOUDFLARE_ACCOUNT_ID: yield* readOptional("CLOUDFLARE_ACCOUNT_ID"),
  };
});
