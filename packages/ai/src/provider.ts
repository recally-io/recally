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

    // gpt-oss models answer in Responses API shape (output[] with output_text),
    // chat models in {response} — normalize all observed shapes to one text.
    const res = raw as {
      response?: unknown;
      output?: Array<{ type?: string; content?: Array<{ type?: string; text?: string }> }>;
      choices?: Array<{ message?: { content?: unknown } }>;
      usage?: {
        prompt_tokens?: number;
        completion_tokens?: number;
        input_tokens?: number;
        output_tokens?: number;
      };
    };

    let text = "";

    if (typeof res?.response === "string" && res.response) {
      text = res.response;
    } else if (res?.response && typeof res.response === "object") {
      text = JSON.stringify(res.response);
    } else {
      const msg = res?.output?.find((o) => o.type === "message");
      const outputText = msg?.content?.find((c) => c.type === "output_text")?.text;
      const choice = res?.choices?.[0]?.message?.content;
      text =
        outputText ?? (typeof choice === "string" ? choice : choice ? JSON.stringify(choice) : "");
    }

    if (!text) {
      throw new AppError(
        "model_unavailable",
        `model ${args.model} returned no text; raw: ${JSON.stringify(raw).slice(0, 300)}`,
        { retryable: true },
      );
    }

    let parsed: unknown;

    try {
      parsed = JSON.parse(text.replace(/<think>[\s\S]*?<\/think>/g, "").trim());
    } catch {
      throw new AppError(
        "model_unavailable",
        `model ${args.model} returned non-JSON output: ${text.slice(0, 300)}`,
        { retryable: true },
      );
    }

    const checked = args.schema.safeParse(parsed);

    if (!checked.success) {
      const issues = checked.error.issues
        .map((i) => `${i.path.join(".") || "output"}: ${i.message}`)
        .join("; ");

      // Include a raw slice — off-schema output is a model-compat bug you
      // can't fix blind (e.g. enum casing, think-tag residue).
      throw new AppError(
        "model_unavailable",
        `model ${args.model} output failed schema: ${issues}. raw: ${text.slice(0, 300)}`,
        { retryable: true },
      );
    }

    return {
      value: checked.data,
      model: args.model,
      usage: {
        inputTokens: res?.usage?.prompt_tokens ?? res?.usage?.input_tokens ?? 0,
        outputTokens: res?.usage?.completion_tokens ?? res?.usage?.output_tokens ?? 0,
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
