// Model roles per plan §10.1. Config is fixed at job start; upgrading the
// capture model never regenerates embeddings.
interface ModelRoles {
  capture: string;
  verify: string;
  summary: string;
  answer: string;
  embedding: string;
}

const DEFAULT_MODELS: ModelRoles = {
  // Compatibility baseline, not a quality claim — M0 evaluates alternatives.
  capture: "@cf/qwen/qwen3-30b-a3b-fp8",
  verify: "@cf/qwen/qwen3-30b-a3b-fp8",
  summary: "@cf/qwen/qwen3-30b-a3b-fp8",
  answer: "@cf/qwen/qwen3-30b-a3b-fp8",
  embedding: "@cf/baai/bge-m3",
};

export interface ModelEnv {
  CAPTURE_MODEL?: string | undefined;
  VERIFY_MODEL?: string | undefined;
  SUMMARY_MODEL?: string | undefined;
  ANSWER_MODEL?: string | undefined;
  EMBEDDING_MODEL?: string | undefined;
}

export function resolveModels(env: ModelEnv): ModelRoles {
  return {
    capture: env.CAPTURE_MODEL ?? DEFAULT_MODELS.capture,
    verify: env.VERIFY_MODEL ?? DEFAULT_MODELS.verify,
    summary: env.SUMMARY_MODEL ?? DEFAULT_MODELS.summary,
    answer: env.ANSWER_MODEL ?? DEFAULT_MODELS.answer,
    embedding: env.EMBEDDING_MODEL ?? DEFAULT_MODELS.embedding,
  };
}
