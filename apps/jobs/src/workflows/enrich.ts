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
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import { definedModelVars, ModelVarsConfig } from "../../../../infra/model-env";
import { ArchiveBucket, Database } from "../../../../infra/resources";
import type { JobsEnv } from "../env";

interface EnrichParams {
  jobId: string;
  libraryId: string;
}

const MAX_TURNS = 16;

// EnrichWorkflow runs the AnalysisAgent over one committed ContentRevision.
// Its toolset is read-only (§10.1) — a summary failure can never corrupt or
// re-fetch the archived original.
export class EnrichWorkflow extends Cloudflare.Workflow<EnrichWorkflow>()(
  "EnrichWorkflow",
  Effect.gen(function* () {
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);
    const bucketClient = yield* Cloudflare.R2.ReadWriteBucket(ArchiveBucket);
    const aiClient = yield* Cloudflare.Workers.AI();
    const modelVars = yield* ModelVarsConfig;

    return Effect.fn(function* (input: EnrichParams) {
      const [db, r2, aiRaw] = yield* Effect.all([dbClient.raw, bucketClient.raw, aiClient.raw]);
      const env: JobsEnv = {
        DB: db,
        ARCHIVE_BUCKET: r2,
        AI: aiRaw,
        ...definedModelVars(modelVars),
      };
      const evidence = new R2EvidenceStore(r2);
      const { jobId, libraryId } = input;

      const job = yield* Cloudflare.Workflows.task(
        "load",
        Effect.tryPromise(async () => {
          const j = await db
            .prepare("SELECT * FROM jobs WHERE id = ? AND library_id = ?")
            .bind(jobId, libraryId)
            .first<{ id: string; item_id: string; payload: string }>();
          if (!j) throw new Error(`job ${jobId} not found`);
          await db
            .prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?")
            .bind(nowIso(), jobId)
            .run();
          return j;
        }).pipe(Effect.orDie),
      );
      const { content_revision_id } = JSON.parse(job.payload) as {
        content_revision_id: string;
      };

      const skill = yield* Cloudflare.Workflows.task(
        "skill",
        Effect.tryPromise(() => loadSkill("analyze")).pipe(Effect.orDie),
      );

      const deps: ToolDeps = {
        evidence,
        runStore: new D1R2RunStore(db, evidence),
        readRevisionBlocks: (_ctx, rev) => evidence.getText(r2Keys.contentBlocks(libraryId, rev)),
        persistArtifact: async (_ctx, artifact) => {
          const artifactId = newId();
          const now = nowIso();
          await evidence.put(
            r2Keys.artifact(libraryId, artifactId),
            JSON.stringify(artifact.output),
            "application/json",
          );
          await db
            .prepare(
              `INSERT INTO ai_artifacts (id, library_id, item_id, input_content_revision_id, type, model_id,
            prompt_version, output, coverage, status, created_at)
           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 'verified', ?)`,
            )
            .bind(
              artifactId,
              libraryId,
              job.item_id,
              artifact.contentRevisionId,
              artifact.type,
              artifact.modelId,
              artifact.promptVersion,
              JSON.stringify(artifact.output),
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

      // Per-run budget: wall-clock deadlines must anchor to the run, not to
      // isolate cold start, so the ctx is built inside the body.
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

      // One durable step per turn (§9.4 replay fence). Serialized Pi state
      // threads through step outputs, so a replayed run restores state from
      // the journal instead of a mutable closure.
      let serialized: string | null = null;
      for (let turn = 0; turn <= MAX_TURNS; turn++) {
        const result = yield* Cloudflare.Workflows.task(
          `turn-${turn}`,
          Effect.tryPromise(async () => {
            const j = await db
              .prepare("SELECT status FROM jobs WHERE id = ?")
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
            return {
              done: r.outcome.kind === "continue" ? null : r.outcome.kind,
              serializedState: r.serializedState,
            };
          }).pipe(Effect.orDie),
        );
        serialized = result.serializedState ?? serialized;
        if (result.done) break;
      }
      void serialized;

      yield* Cloudflare.Workflows.task(
        "finish",
        Effect.tryPromise(async () => {
          await db
            .prepare(
              "UPDATE jobs SET status = 'succeeded', result = 'summarized', updated_at = ? WHERE id = ?",
            )
            .bind(nowIso(), jobId)
            .run();
        }).pipe(Effect.orDie),
      );
    });
  }).pipe(
    Effect.provide(
      Layer.mergeAll(
        Cloudflare.D1.QueryDatabaseBinding,
        Cloudflare.R2.ReadWriteBucketBinding,
        Cloudflare.Workers.AIBinding,
      ),
    ),
  ),
) {}
