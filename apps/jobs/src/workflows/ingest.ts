import { AppError, DEFAULT_LIMITS, newId, nowIso } from "@recally/domain";
import { R2EvidenceStore } from "@recally/platform-cloudflare";
import { type CaptureRunRow, getJob, r2Keys, updateJobStatus } from "@recally/storage";
import { captureTools } from "@recally/tools";
import * as Cloudflare from "alchemy/Cloudflare";
import * as Effect from "effect/Effect";
import * as Layer from "effect/Layer";
import { definedModelVars, ModelVarsConfig } from "../../../../infra/model-env";
import { ArchiveBucket, Database } from "../../../../infra/resources";
import { assembleRun } from "../context";
import type { IngestEnv } from "../env";

interface IngestParams {
  jobId: string;
  libraryId: string;
}

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

    const markWorkflowError = (
      db: IngestEnv["DB"],
      jobId: string,
      runId: string | null,
      err: unknown,
    ) =>
      Effect.tryPromise(async () => {
        const now = nowIso();

        const stmts = [
          db
            .prepare(
              "UPDATE jobs SET status = 'failed', result = 'workflow_error', error = ?, updated_at = ? WHERE id = ?",
            )
            .bind(err instanceof Error ? err.message.slice(0, 500) : String(err), now, jobId),
        ];

        if (runId) {
          stmts.push(
            db
              .prepare(
                "UPDATE capture_runs SET status = 'failed', outcome_code = 'failed', updated_at = ? WHERE id = ?",
              )
              .bind(now, runId),
          );
        }

        await db.batch(stmts);
      }).pipe(Effect.ignore);

    return Effect.fn(function* (input: IngestParams) {
      const [db, r2, aiRaw, browserRaw] = yield* Effect.all([
        dbClient.raw,
        bucketClient.raw,
        aiClient.raw,
        browserClient.raw,
      ]);

      const env: IngestEnv = {
        DB: db,
        ARCHIVE_BUCKET: r2,
        AI: aiRaw,
        BROWSER: browserRaw as unknown as IngestEnv["BROWSER"],
        ...definedModelVars(modelVars),
      };

      const evidence = new R2EvidenceStore(r2);
      const { jobId, libraryId } = input;
      // Set once "init" completes so the failure handler can also close out
      // the capture_runs row (best effort — the job row is the source of truth).
      let failedRunId: string | null = null;

      const pipeline = Effect.gen(function* () {
        const init = yield* Cloudflare.Workflows.task(
          "init",
          Effect.tryPromise(async () => {
            const job = await getJob(db, libraryId, jobId);

            if (!job) throw new AppError("not_found", `job ${jobId}`);

            const run = await db
              .prepare("SELECT * FROM capture_runs WHERE job_id = ? AND library_id = ?")
              .bind(jobId, libraryId)
              .first<CaptureRunRow>();

            if (!run) throw new AppError("not_found", `run for job ${jobId}`);

            const item = await db
              .prepare(
                "SELECT deleted_at, capture_generation FROM items WHERE id = ? AND library_id = ?",
              )
              .bind(run.item_id, libraryId)
              .first<{ deleted_at: string | null; capture_generation: number }>();

            if (!item || item.deleted_at || item.capture_generation !== run.generation) {
              await updateJobStatus(db, jobId, "failed", {
                error: "stale_generation",
              });

              return null;
            }

            const attemptId = newId();
            const now = nowIso();
            await db.batch([
              db
                .prepare(
                  "INSERT INTO capture_attempts (id, library_id, run_id, attempt_no, kind, status, started_at) VALUES (?, ?, ?, 1, 'fetch', 'running', ?)",
                )
                .bind(attemptId, libraryId, run.id, now),
              db
                .prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?")
                .bind(now, jobId),
              db
                .prepare("UPDATE capture_runs SET status = 'running', updated_at = ? WHERE id = ?")
                .bind(now, run.id),
            ]);

            return { run, job, attemptId };
          }).pipe(Effect.orDie),
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
            Effect.tryPromise(async () => {
              const { deps, ctx } = await assembleRun(env, run, attemptId);

              const source = await deps.runStore.saveSource(ctx, {
                url: payload.url ?? "manual:input",
                kind: "manual",
                contentType: "text/plain",
                body: manualText,
              });

              const result = await deps.committer!.commit(ctx, {
                sourceIds: [source.sourceId],
                selectedBlockRanges: [],
                requiredAssetIds: [],
                optionalAssetIds: [],
                metadataEvidence: {},
                missingParts: [],
                claimedQuality: "complete",
              });

              if (result.status !== "accepted")
                throw new AppError("internal", "manual commit rejected");
              await updateJobStatus(db, jobId, "succeeded", {
                result: "complete",
              });
            }).pipe(Effect.orDie),
          );

          return;
        }

        const maxTurns = DEFAULT_LIMITS.agent.maxModelCalls;

        for (let turn = 0; turn <= maxTurns; turn++) {
          const result = yield* Cloudflare.Workflows.task(
            `turn-${turn}`,
            Effect.tryPromise(async () => {
              const { ctx, deps, runtime, skill } = await assembleRun(env, run, attemptId);

              // Cancellation/generation fence before every model call (§9.7).
              const job = await getJob(db, libraryId, jobId);

              if (!job || job.status === "cancelled" || job.cancel_requested) {
                return { kind: "cancelled" as const };
              }

              const freshRun = await db
                .prepare("SELECT pi_state_key FROM capture_runs WHERE id = ?")
                .bind(run.id)
                .first<{ pi_state_key: string | null }>();

              const serializedState = freshRun?.pi_state_key
                ? await evidence.getText(freshRun.pi_state_key)
                : null;

              const turnResult = await runtime.runTurn({
                serializedState,
                systemPrompt: skill.prompt,
                goal: `Archive the content the user intended by this URL: ${run.target_url}`,
                tools: captureTools(deps),
                ctx,
                model: env.CAPTURE_MODEL ?? "@cf/qwen/qwen3-30b-a3b-fp8",
                maxTurns,
                turn,
              });

              // Persist Pi state + audit events inside the same step (§9.4).
              // Sequences come from MAX+1 like saveObservation's mid-turn writes —
              // a separate turn*100 namespace collides with them under UNIQUE.
              const stateKey = r2Keys.piState(libraryId, run.id, turn);
              await evidence.put(stateKey, turnResult.serializedState, "application/json");
              const now = nowIso();

              const seqBase =
                (
                  await db
                    .prepare(
                      "SELECT COALESCE(MAX(sequence), -1) + 1 AS s FROM agent_events WHERE run_id = ? AND attempt_id = ?",
                    )
                    .bind(run.id, attemptId)
                    .first<{ s: number }>()
                )?.s ?? 0;

              await db.batch([
                db
                  .prepare(
                    "INSERT INTO agent_states (id, library_id, run_id, sequence, state_key, created_at) VALUES (?, ?, ?, ?, ?, ?)",
                  )
                  .bind(newId(), libraryId, run.id, turn, stateKey, now),
                db
                  .prepare("UPDATE capture_runs SET pi_state_key = ?, updated_at = ? WHERE id = ?")
                  .bind(stateKey, now, run.id),
                ...turnResult.events.map((e, i) =>
                  db
                    .prepare(
                      `INSERT INTO agent_events (id, library_id, run_id, attempt_id, sequence, kind, reason_code, summary, evidence_ref, usage, created_at)
               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
                    )
                    .bind(
                      newId(),
                      libraryId,
                      run.id,
                      attemptId,
                      seqBase + i,
                      e.kind,
                      e.reasonCode ?? null,
                      e.toolName ? `${e.toolName}: ${e.summary ?? ""}` : (e.summary ?? null),
                      e.evidenceRef ?? null,
                      e.usage ? JSON.stringify(e.usage) : null,
                      now,
                    ),
                ),
              ]);

              return { kind: "outcome" as const, outcome: turnResult.outcome };
            }).pipe(Effect.orDie),
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
            Effect.tryPromise(async () => {
              const now = nowIso();

              if (outcome.kind === "proposed") {
                // ArchiveService already published the snapshot.
                await db.batch([
                  db
                    .prepare(
                      "UPDATE jobs SET status = 'succeeded', result = 'complete', updated_at = ? WHERE id = ?",
                    )
                    .bind(now, jobId),
                  db
                    .prepare(
                      "UPDATE capture_runs SET status = 'succeeded', outcome_code = 'complete', updated_at = ? WHERE id = ?",
                    )
                    .bind(now, run.id),
                ]);
              } else if (outcome.kind === "finished") {
                const o = outcome as { outcome: string; reasonCode: string };

                const jobStatus =
                  o.outcome === "needs_input"
                    ? "needs_input"
                    : o.outcome === "complete" || o.outcome === "partial"
                      ? "succeeded"
                      : "failed";

                await db.batch([
                  db
                    .prepare("UPDATE jobs SET status = ?, result = ?, updated_at = ? WHERE id = ?")
                    .bind(jobStatus, o.reasonCode, now, jobId),
                  db
                    .prepare(
                      "UPDATE capture_runs SET status = ?, outcome_code = ?, updated_at = ? WHERE id = ?",
                    )
                    .bind(
                      jobStatus === "succeeded" ? "succeeded" : jobStatus,
                      o.outcome,
                      now,
                      run.id,
                    ),
                ]);
              } else {
                await db.batch([
                  db
                    .prepare(
                      "UPDATE jobs SET status = 'failed', result = 'budget_exceeded', updated_at = ? WHERE id = ?",
                    )
                    .bind(now, jobId),
                  db
                    .prepare(
                      "UPDATE capture_runs SET status = 'failed', outcome_code = 'budget_exhausted', updated_at = ? WHERE id = ?",
                    )
                    .bind(now, run.id),
                ]);
              }
            }).pipe(Effect.orDie),
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
            yield* markWorkflowError(db, jobId, failedRunId, err);

            return yield* Effect.die(err);
          }),
        ),
      );
    });
  }).pipe(
    Effect.provide(
      Layer.mergeAll(
        Cloudflare.D1.QueryDatabaseBinding,
        Cloudflare.R2.ReadWriteBucketBinding,
        Cloudflare.Workers.AIBinding,
        Cloudflare.Workers.BrowserBinding,
      ),
    ),
  ),
) {}
