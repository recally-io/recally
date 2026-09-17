import { DEFAULT_LIMITS } from "@recally/domain";
import { R2EvidenceStore } from "@recally/platform-cloudflare";
import { updateJobStatus } from "@recally/storage";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import { ModelVarsConfig } from "../../../../infra/model-env";
import { ArchiveBucket, Database } from "../../../../infra/resources";
import { JobsBaseLayer, makeJobsEnv } from "../bindings";
import { markWorkflowError, type WorkflowInput } from "./common";
import { initializeCapture, finalizeCapture } from "./ingest/lifecycle";
import { commitManual } from "./ingest/manual";
import { runCaptureTurn } from "./ingest/turn";
import type { IngestEnv } from "../env";

// IngestWorkflow = durable supervisor (§9.1). It does NOT choose acquisition
// strategy — each turn delegates to Pi + capture skill; it owns persistence,
// budget fencing, retry and commit boundaries.
export class IngestWorkflow extends Cloudflare.Workflow<IngestWorkflow>()(
  "IngestWorkflow",
  Effect.gen(function* () {
    // Construction phase: register bindings + resolve Config. Native handles
    // resolve per run in the body (Runtime phase).
    const dbClient = yield* Cloudflare.D1.QueryDatabase(Database);
    const bucketClient = yield* Cloudflare.R2.ReadWriteBucket(ArchiveBucket);
    const aiClient = yield* Cloudflare.Workers.AI();
    const browserClient = yield* Cloudflare.Browser("BROWSER");
    const modelVars = yield* ModelVarsConfig;

    return Effect.fn(function* (input: WorkflowInput) {
      const [db, r2, aiRaw, browserRaw] = yield* Effect.all([
        dbClient.raw,
        bucketClient.raw,
        aiClient.raw,
        browserClient.raw,
      ]);

      const env: IngestEnv = {
        ...makeJobsEnv(db, r2, aiRaw, modelVars),
        BROWSER: browserRaw as unknown as IngestEnv["BROWSER"],
      };

      const evidence = new R2EvidenceStore(r2);
      const { jobId, libraryId } = input;
      // Set once "init" completes so the failure handler can also close out
      // the capture_runs row (best effort — the job row is the source of truth).
      let failedRunId: string | null = null;

      const pipeline = Effect.gen(function* () {
        const init = yield* Cloudflare.Workflows.task(
          "init",
          Effect.tryPromise(() => initializeCapture(db, input)).pipe(Effect.orDie),
        );

        if (!init) return;
        failedRunId = init.run.id;
        const { run, attemptId } = init;

        const payload = JSON.parse(init.job.payload) as {
          kind: string;
          url?: string;
          text?: string;
          content?: string;
        };

        const manualText = payload.text ?? payload.content;

        // Manual content: user-submitted text archives directly — no agent
        // ritual re-fetch just to satisfy a fetch-first default (§5.6).
        if (payload.kind === "manual" && manualText) {
          yield* Cloudflare.Workflows.task(
            "manual-commit",
            Effect.tryPromise(() =>
              commitManual(env, run, attemptId, jobId, manualText, payload.url),
            ).pipe(Effect.orDie),
          );

          return;
        }

        const maxTurns = DEFAULT_LIMITS.agent.maxModelCalls;

        for (let turn = 0; turn <= maxTurns; turn++) {
          const result = yield* Cloudflare.Workflows.task(
            `turn-${turn}`,
            Effect.tryPromise(() =>
              runCaptureTurn(env, evidence, { jobId, libraryId, run, attemptId, turn, maxTurns }),
            ).pipe(Effect.orDie),
          );

          if (result.kind === "cancelled") {
            yield* Cloudflare.Workflows.task(
              "mark-cancelled",
              Effect.tryPromise(() => updateJobStatus(db, jobId, "cancelled")).pipe(Effect.orDie),
            );

            return;
          }

          const outcome = result.outcome;

          if (outcome.kind === "continue") continue;

          yield* Cloudflare.Workflows.task(
            "finalize",
            Effect.tryPromise(() => finalizeCapture(db, jobId, run.id, outcome)).pipe(Effect.orDie),
          );

          return;
        }
      });

      // A step that exhausts retries would otherwise leave the job 'running'
      // forever (observed: agent_events UNIQUE collision errored a run while
      // its job row stayed running). Mark it failed, then re-die so the
      // workflow instance still reports the error.
      yield* pipeline.pipe(
        Effect.catchDefect((err) =>
          Effect.gen(function* () {
            yield* markWorkflowError(
              db,
              jobId,
              failedRunId,
              err instanceof Error ? err : String(err),
            );

            return yield* Effect.die(err);
          }),
        ),
      );
    });
  }).pipe(Effect.provide(Layer.mergeAll(JobsBaseLayer, Cloudflare.Workers.BrowserBinding))),
) {}
