import { AppError } from "@recally/domain";
import { z } from "zod";

// Narrow adapter per plan §10.1/ADR-06: one interface over Workers AI today,
// external providers are explicitly-enabled additions later — not a general
// runtime abstraction.

export interface ModelUsage {
  inputTokens: number;
  outputTokens: number;
}

export interface ModelResult<T> {
  value: T;
  model: string;
  usage: ModelUsage;
}

export interface ModelProvider {
  // Structured generation. Callers never trust raw model output — the schema
  // is enforced here, and a parse failure is a distinct error class.
  generate<T>(args: {
    model: string;
    system: string;
    prompt: string;
    schema: z.ZodType<T>;
    maxOutputTokens?: number;
  }): Promise<ModelResult<T>>;

  embed(args: {
    model: string;
    texts: string[];
  }): Promise<{ vectors: number[][]; usage: ModelUsage }>;
}

interface AiBinding {
  run(model: string, input: Record<string, unknown>): Promise<unknown>;
}

export class WorkersAIProvider implements ModelProvider {
  constructor(private readonly ai: AiBinding) {}

  async generate<T>(args: {
    model: string;
    system: string;
    prompt: string;
    schema: z.ZodType<T>;
    maxOutputTokens?: number;
  }): Promise<ModelResult<T>> {
    let raw: unknown;
    try {
      raw = await this.ai.run(args.model, {
        messages: [
          { role: "system", content: args.system },
          { role: "user", content: args.prompt },
        ],
        response_format: {
          type: "json_schema",
          json_schema: z.toJSONSchema(args.schema),
        },
        max_tokens: args.maxOutputTokens ?? 2048,
      });
    } catch (err) {
      throw new AppError(
        "model_unavailable",
        `model ${args.model} failed: ${err instanceof Error ? err.message : String(err)}`,
        { retryable: true },
      );
    }

    const res = raw as {
      response?: unknown;
      usage?: { prompt_tokens?: number; completion_tokens?: number };
    };
    const text =
      typeof res?.response === "string" ? res.response : JSON.stringify(res?.response ?? {});
    let parsed: unknown;
    try {
      parsed = JSON.parse(text);
    } catch {
      throw new AppError("model_unavailable", "model returned non-JSON output");
    }
    const value = args.schema.parse(parsed);
    return {
      value,
      model: args.model,
      usage: {
        inputTokens: res?.usage?.prompt_tokens ?? 0,
        outputTokens: res?.usage?.completion_tokens ?? 0,
      },
    };
  }

  async embed(args: { model: string; texts: string[] }) {
    const res = (await this.ai.run(args.model, { text: args.texts })) as {
      data?: number[][];
      usage?: { prompt_tokens?: number };
    };
    return {
      vectors: res?.data ?? [],
      usage: { inputTokens: res?.usage?.prompt_tokens ?? 0, outputTokens: 0 },
    };
  }
}
