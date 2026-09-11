import { sha256Hex } from "./hash";

// Bumping a version changes operation keys, so re-runs produce new artifacts
// instead of overwriting old ones (plan §9.6).
export const PIPELINE_VERSION = "pipe-0001";
export const POLICY_VERSION = "pol-0001";
export const EXTRACTOR_VERSION = "ext-0001";
export const TOKENIZER_VERSION = "tok-0001";
export const EMBEDDING_VERSION = "emb-0001";
export const PROMPT_VERSIONS = {
  captureDecide: "cap-decide-0001",
  captureVerify: "cap-verify-0001",
  summarize: "sum-0001",
  answer: "ans-0001",
} as const;

// Model config is fixed at job start; embedding upgrades must not invalidate
// unrelated stages (plan §10.1).
export async function modelConfigVersion(models: Record<string, string>): Promise<string> {
  const stable = Object.keys(models)
    .sort()
    .map((k) => `${k}=${models[k]}`)
    .join(";");
  return `mcfg-${(await sha256Hex(stable)).slice(0, 16)}`;
}
