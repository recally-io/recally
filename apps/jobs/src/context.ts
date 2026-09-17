import { PiRuntime } from "@recally/agent-runtime";
import {
  CAPTURE_VERIFY_SYSTEM,
  resolveModels,
  verificationSchema,
  WorkersAIProvider,
} from "@recally/ai";
import type { CaptureBudgets, ToolContext } from "@recally/capture";
import { checkUrlTarget } from "@recally/capture";
import { DEFAULT_LIMITS } from "@recally/domain";
import {
  ArchiveService,
  BrowserSession,
  cfStreamFn,
  D1R2RunStore,
  R2EvidenceStore,
  resolveCfModel,
} from "@recally/platform-cloudflare";
import { type LoadedSkill, loadSkill } from "@recally/skills";
import type { CaptureRunRow } from "@recally/storage";
import type { ToolDeps } from "@recally/tools";
import { safeFetch } from "@recally/tools";
import type { IngestEnv } from "./env";

export interface RunAssembly {
  ctx: ToolContext;
  deps: ToolDeps;
  runtime: PiRuntime;
  skill: LoadedSkill;
  run: CaptureRunRow;
}

export function budgetFor(run: CaptureRunRow): CaptureBudgets {
  // Run row carries budget JSON when customized; defaults otherwise (§16.1).
  const b = (run.budget ? JSON.parse(run.budget) : {}) as Record<string, number | undefined>;
  const L = DEFAULT_LIMITS;

  return {
    turnsLeft: b.turns ?? L.agent.maxModelCalls,
    modelCallsLeft: b.model_calls ?? L.agent.maxModelCalls,
    inputTokensLeft: b.input_tokens ?? L.agent.maxInputTokens,
    outputTokensLeft: b.output_tokens ?? L.agent.maxOutputTokens,
    browserEpisodesLeft: b.browser_episodes ?? L.browser.maxEpisodesPerRun,
    browserMsLeft: b.browser_ms ?? L.browser.maxSessionMs,
    browserActionsLeft: b.browser_actions ?? L.browser.maxActions,
    downloadBytesLeft: b.download_bytes ?? 50 * 1024 * 1024,
    wallClockDeadlineMs: Date.now() + L.run.maxWallClockMs,
  };
}

// Adapter network access re-uses the same SSRF policy as web_fetch (§5.7).
async function adapterFetch(url: string): Promise<unknown> {
  checkUrlTarget(url);
  const res = await safeFetch(url, fetch, { maxBytes: 2 * 1024 * 1024 });

  return JSON.parse(new TextDecoder().decode(res.body));
}

export async function assembleRun(
  env: IngestEnv,
  run: CaptureRunRow,
  attemptId: string,
): Promise<RunAssembly> {
  const evidence = new R2EvidenceStore(env.ARCHIVE_BUCKET);
  const runStore = new D1R2RunStore(env.DB, evidence);
  const models = resolveModels(env);
  const ai = new WorkersAIProvider(env.AI as never);

  const ctx: ToolContext = {
    libraryId: run.library_id,
    itemId: run.item_id,
    runId: run.id,
    attemptId,
    generation: run.generation,
    skillRevision: run.skill_revision,
    toolsetVersion: run.toolset_version,
    policyVersion: run.policy_version,
    budget: budgetFor(run),
  };

  const deps: ToolDeps = {
    evidence,
    runStore,
    adapterFetch,
    browser: {
      open: (url) => BrowserSession.open({ BROWSER: env.BROWSER }, url),
    },
    committer: new ArchiveService({
      db: env.DB,
      evidence,
      runStore,
      // Independent verifier (§7.3): judges content evidence, not the agent's
      // self-report. Own budget reservation lands with usage settlement (M2).
      verify: async (input) =>
        (
          await ai.generate({
            model: models.verify,
            system: CAPTURE_VERIFY_SYSTEM,
            prompt: [
              `claimed quality: ${input.claimedQuality}`,
              input.title ? `title: ${input.title}` : null,
              `missing parts reported: ${input.missingParts.join("; ") || "none"}`,
              "--- head ---",
              input.headText,
              "--- tail ---",
              input.tailText,
            ]
              .filter(Boolean)
              .join("\n"),
            schema: verificationSchema,
          })
        ).value,
    }),
  };

  const skill = await loadSkill("capture");

  const runtime = new PiRuntime({
    streamFn: cfStreamFn({
      AI: env.AI as never,
      AI_GATEWAY_NAME: env.AI_GATEWAY_NAME,
      CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
    }),
    resolveModel: (id) =>
      resolveCfModel(id, {
        AI: env.AI as never,
        AI_GATEWAY_NAME: env.AI_GATEWAY_NAME,
        CLOUDFLARE_ACCOUNT_ID: env.CLOUDFLARE_ACCOUNT_ID,
      }),
    runtimeVersion: run.agent_runtime_version,
  });

  return { ctx, deps, runtime, skill, run };
}
