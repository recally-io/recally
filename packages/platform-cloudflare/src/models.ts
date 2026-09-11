import type { StreamFn } from "@earendil-works/pi-agent-core";
import type { Model } from "@earendil-works/pi-ai";
import { cloudflareWorkersAIProvider } from "@earendil-works/pi-ai/providers/cloudflare-workers-ai";

// Pi model resolution over Workers AI. Two transports, chosen by env:
//  - AI binding passthrough via AI Gateway (AI_GATEWAY_NAME set): calls
//    env.AI.fetch() so no API token is needed.
//  - Direct HTTPS openai-completions endpoint: needs CLOUDFLARE_ACCOUNT_ID +
//    a CF API token (secret, never in vars).
// Exact wiring is M0-T02 verification scope.

type CfModel = Model<"openai-completions">;

export interface ModelEnv {
  AI: {
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
  if (env.CLOUDFLARE_ACCOUNT_ID) {
    return {
      ...model,
      baseUrl: model.baseUrl?.replace("{CLOUDFLARE_ACCOUNT_ID}", env.CLOUDFLARE_ACCOUNT_ID),
    };
  }
  return model;
}

export function cfStreamFn(env: ModelEnv): StreamFn {
  const bindingFetch = env.AI.fetch?.bind(env.AI);
  return (model, context, options) =>
    provider.streamSimple(model as CfModel, context, {
      ...options,
      // Binding-routed requests are pre-authenticated; the sentinel satisfies
      // the API layer's auth check and is stripped by the gateway.
      ...(bindingFetch
        ? {
            fetch: bindingFetch,
            headers: {
              "cf-aig-authorization": "Bearer cloudflare-gateway-binding",
            },
          }
        : {}),
    });
}
