import type { StreamFn } from "@earendil-works/pi-agent-core";
import type { Model } from "@earendil-works/pi-ai";
import { cloudflareWorkersAIProvider } from "@earendil-works/pi-ai/providers/cloudflare-workers-ai";

// Pi model resolution over Workers AI. Transport picks itself by env:
//  - AI_GATEWAY_NAME set: binding passthrough through the gateway's
//    provider route (env.AI.fetch), giving log/observability in AI Gateway.
//  - Otherwise: env.AI.run() with returnRawResponse, which returns the
//    provider's native OpenAI-format Response — no gateway or API token
//    needed, and it carries tool calls + streaming identically.
// Verified against workerd 2026-09: env.AI.run answers OpenAI-shaped JSON
// for chat models; env.AI.fetch's gateway route 400s when no gateway exists.

type CfModel = Model<"openai-completions">;

export interface ModelEnv {
  AI: {
    run?: (model: string, input: unknown, options?: unknown) => Promise<Response>;
    fetch?: (input: Request | string | URL, init?: RequestInit) => Promise<Response>;
  };
  CLOUDFLARE_ACCOUNT_ID?: string | undefined;
  AI_GATEWAY_NAME?: string | undefined;
}

const provider = cloudflareWorkersAIProvider();

export function resolveCfModel(modelId: string, env: ModelEnv): CfModel {
  const model = provider.getModels().find((m: CfModel) => m.id === modelId);

  if (!model) throw new Error(`unknown Workers AI model ${modelId}`);

  if (env.AI_GATEWAY_NAME) {
    return {
      ...model,
      baseUrl: `https://workers-binding.ai/ai-gateway/gateways/${env.AI_GATEWAY_NAME}/workers-ai`,
    };
  }

  // run() mode: the fetch adapter ignores the URL, but it must parse.
  return { ...model, baseUrl: "https://workers-binding.ai/ai/v1" };
}

export function cfStreamFn(env: ModelEnv): StreamFn {
  let boundFetch: typeof fetch | undefined;

  if (env.AI_GATEWAY_NAME && env.AI.fetch) {
    boundFetch = env.AI.fetch.bind(env.AI);
  } else if (env.AI.run) {
    const run = env.AI.run.bind(env.AI);
    boundFetch = (_input, init) => {
      const { model, ...input } = JSON.parse(String(init?.body ?? "{}"));

      // The AI binding validates inputs against a strict schema: assistant
      // messages carrying tool_calls must have string content (null 400s).
      if (Array.isArray(input.messages)) {
        input.messages = input.messages.map((m: { content?: unknown }) => ({
          ...m,
          content: m.content ?? "",
        }));
      }

      return run(model, input, { returnRawResponse: true });
    };
  }

  return (model, context, options) =>
    provider.streamSimple(model as CfModel, context, {
      ...options,
      // Binding-routed calls are pre-authenticated; the sentinel satisfies
      // the API layer's auth check and is stripped before the provider.
      ...(boundFetch
        ? {
            fetch: boundFetch,
            headers: { "cf-aig-authorization": "Bearer cloudflare-gateway-binding" },
          }
        : {}),
    });
}
