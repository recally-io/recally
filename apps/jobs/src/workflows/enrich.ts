import { WorkflowEntrypoint, type WorkflowEvent, type WorkflowStep } from "cloudflare:workers";
import { PiRuntime } from "@recally/agent-runtime";
import { resolveModels } from "@recally/ai";
import { newId, nowIso } from "@recally/domain";
import {
  cfStreamFn,
  D1R2RunStore,
  R2EvidenceStore,
  resolveCfModel,
} from "@recally/platform-cloudflare";
import { loadSkill } from "@recally/skills";
import { r2Keys } from "@recally/storage";
import type { ToolDeps } from "@recally/tools";
import { analysisTools } from "@recally/tools";
import type { Env } from "../env";

interface EnrichParams {
  jobId: string;
  libraryId: string;
}

const MAX_TURNS = 16;

// EnrichWorkflow runs the AnalysisAgent over one committed ContentRevision.
// Its toolset is read-only (§10.1) — a summary failure can never corrupt or
// re-fetch the archived original.
export class EnrichWorkflow extends WorkflowEntrypoint<Env, EnrichParams> {
  async run(event: WorkflowEvent<EnrichParams>, step: WorkflowStep) {
    const { jobId, libraryId } = event.payload;
    const env = this.env;
    const evidence = new R2EvidenceStore(env.ARCHIVE_BUCKET);

    const job = await step.do("load", async () => {
      const j = await env.DB.prepare("SELECT * FROM jobs WHERE id = ? AND library_id = ?")
        .bind(jobId, libraryId)
        .first<{ id: string; item_id: string; payload: string }>();
      if (!j) throw new Error(`job ${jobId} not found`);
      await env.DB.prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?")
        .bind(nowIso(), jobId)
        .run();
      return j;
    });
    const { content_revision_id } = JSON.parse(job.payload) as {
      content_revision_id: string;
    };

    const skill = await step.do("skill", () => loadSkill("analyze"));

    const deps: ToolDeps = {
      evidence,
      runStore: new D1R2RunStore(env.DB, evidence),
      readRevisionBlocks: (_ctx, rev) => evidence.getText(r2Keys.contentBlocks(libraryId, rev)),
      persistArtifact: async (_ctx, artifact) => {
        const artifactId = newId();
        const now = nowIso();
        await evidence.put(
          r2Keys.artifact(libraryId, artifactId),
          JSON.stringify(artifact.output),
          "application/json",
        );
        await env.DB.prepare(
          `INSERT INTO ai_artifacts (id, library_id, item_id, content_revision_id, type, model_id,
            prompt_version, output_key, coverage, created_at)
           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        )
          .bind(
            artifactId,
            libraryId,
            job.item_id,
            artifact.contentRevisionId,
            artifact.type,
            artifact.modelId,
            artifact.promptVersion,
            r2Keys.artifact(libraryId, artifactId),
            JSON.stringify(artifact.coverage ?? null),
            now,
          )
          .run();
        return artifactId;
      },
    };

    const models = resolveModels(env);
    const runtime = new PiRuntime({
      streamFn: cfStreamFn(env as never),
      resolveModel: (id) => resolveCfModel(id, env as never),
    });
    const ctx = {
      libraryId,
      itemId: job.item_id,
      runId: jobId,
      attemptId: jobId,
      generation: 1,
      skillRevision: skill.revision,
      toolsetVersion: "toolset-0001",
      policyVersion: "policy-0001",
      budget: {
        turnsLeft: MAX_TURNS,
        modelCallsLeft: MAX_TURNS,
        inputTokensLeft: 60_000,
        outputTokensLeft: 8_000,
        browserEpisodesLeft: 0,
        browserMsLeft: 0,
        browserActionsLeft: 0,
        downloadBytesLeft: 0,
        wallClockDeadlineMs: Date.now() + 10 * 60_000,
      },
    };

    let serialized: string | null = null;
    for (let turn = 0; turn <= MAX_TURNS; turn++) {
      const result = await step.do(`turn-${turn}`, async () => {
        const j = await env.DB.prepare("SELECT status FROM jobs WHERE id = ?")
          .bind(jobId)
          .first<{ status: string }>();
        if (j?.status === "cancelled") return { done: "cancelled" as const };
        const r = await runtime.runTurn({
          serializedState: serialized,
          systemPrompt: skill.prompt,
          goal: `Summarize content revision ${content_revision_id} for the user's reading archive. Read the outline first, cover the body, cite block ids.`,
          tools: analysisTools(deps),
          ctx,
          model: models.summary,
          maxTurns: MAX_TURNS,
          turn,
        });
        serialized = r.serializedState;
        return { done: r.outcome.kind === "continue" ? null : r.outcome.kind };
      });
      if (result.done) break;
    }

    await step.do("finish", async () => {
      await env.DB.prepare(
        "UPDATE jobs SET status = 'succeeded', result = 'summarized', updated_at = ? WHERE id = ?",
      )
        .bind(nowIso(), jobId)
        .run();
    });
  }
}
