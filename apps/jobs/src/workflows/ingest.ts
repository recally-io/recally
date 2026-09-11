import { WorkflowEntrypoint, type WorkflowEvent, type WorkflowStep } from "cloudflare:workers";
import { AppError, DEFAULT_LIMITS, newId, nowIso } from "@recally/domain";
import { R2EvidenceStore } from "@recally/platform-cloudflare";
import { type CaptureRunRow, getJob, r2Keys, updateJobStatus } from "@recally/storage";
import { captureTools } from "@recally/tools";
import { assembleRun } from "../context";
import type { Env } from "../env";

interface IngestParams {
  jobId: string;
  libraryId: string;
}

// IngestWorkflow = durable supervisor (§9.1). It does NOT choose acquisition
// strategy — each turn delegates to Pi + capture skill; it owns persistence,
// budget fencing, retry and commit boundaries.
export class IngestWorkflow extends WorkflowEntrypoint<Env, IngestParams> {
  async run(event: WorkflowEvent<IngestParams>, step: WorkflowStep) {
    const { jobId, libraryId } = event.payload;
    const env = this.env;
    const evidence = new R2EvidenceStore(env.ARCHIVE_BUCKET);

    const init = await step.do("init", async () => {
      const job = await getJob(env.DB, libraryId, jobId);
      if (!job) throw new AppError("not_found", `job ${jobId}`);
      const run = await env.DB.prepare(
        "SELECT * FROM capture_runs WHERE job_id = ? AND library_id = ?",
      )
        .bind(jobId, libraryId)
        .first<CaptureRunRow>();
      if (!run) throw new AppError("not_found", `run for job ${jobId}`);

      const item = await env.DB.prepare(
        "SELECT deleted_at, capture_generation FROM items WHERE id = ? AND library_id = ?",
      )
        .bind(run.item_id, libraryId)
        .first<{ deleted_at: string | null; capture_generation: number }>();
      if (!item || item.deleted_at || item.capture_generation !== run.generation) {
        await updateJobStatus(env.DB, jobId, "failed", {
          error: "stale_generation",
        });
        return null;
      }

      const attemptId = newId();
      const now = nowIso();
      await env.DB.batch([
        env.DB.prepare(
          "INSERT INTO capture_attempts (id, library_id, run_id, attempt_no, kind, status, started_at) VALUES (?, ?, ?, 1, 'fetch', 'running', ?)",
        ).bind(attemptId, libraryId, run.id, now),
        env.DB.prepare("UPDATE jobs SET status = 'running', updated_at = ? WHERE id = ?").bind(
          now,
          jobId,
        ),
        env.DB.prepare(
          "UPDATE capture_runs SET status = 'running', updated_at = ? WHERE id = ?",
        ).bind(now, run.id),
      ]);
      return { run, job, attemptId };
    });
    if (!init) return;
    const { run, attemptId } = init;
    const payload = JSON.parse(init.job.payload) as {
      kind: string;
      url?: string;
      text?: string;
      content?: string;
    };
    const manualText = payload.text ?? payload.content;

    // Manual content: user-submitted text archives directly — no agent ritual
    // re-fetch just to satisfy a fetch-first default (§5.6).
    if (payload.kind === "manual" && manualText) {
      await step.do("manual-commit", async () => {
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
        if (result.status !== "accepted") throw new AppError("internal", "manual commit rejected");
        await updateJobStatus(env.DB, jobId, "succeeded", {
          result: "complete",
        });
      });
      return;
    }

    const maxTurns = DEFAULT_LIMITS.agent.maxModelCalls;
    try {
      for (let turn = 0; turn <= maxTurns; turn++) {
      const result = await step.do(`turn-${turn}`, async () => {
        const { ctx, deps, runtime, skill } = await assembleRun(env, run, attemptId);

        // Cancellation/generation fence before every model call (§9.7).
        const job = await getJob(env.DB, libraryId, jobId);
        if (!job || job.status === "cancelled" || job.cancel_requested) {
          return { kind: "cancelled" as const };
        }

        const freshRun = await env.DB.prepare("SELECT pi_state_key FROM capture_runs WHERE id = ?")
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
            await env.DB.prepare(
              "SELECT COALESCE(MAX(sequence), -1) + 1 AS s FROM agent_events WHERE run_id = ? AND attempt_id = ?",
            )
              .bind(run.id, attemptId)
              .first<{ s: number }>()
          )?.s ?? 0;
        await env.DB.batch([
          env.DB.prepare(
            "INSERT INTO agent_states (id, library_id, run_id, sequence, state_key, created_at) VALUES (?, ?, ?, ?, ?, ?)",
          ).bind(newId(), libraryId, run.id, turn, stateKey, now),
          env.DB.prepare(
            "UPDATE capture_runs SET pi_state_key = ?, updated_at = ? WHERE id = ?",
          ).bind(stateKey, now, run.id),
          ...turnResult.events.map((e, i) =>
            env.DB.prepare(
              `INSERT INTO agent_events (id, library_id, run_id, attempt_id, sequence, kind, reason_code, summary, evidence_ref, usage, created_at)
               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
            ).bind(
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
      });

      if (result.kind === "cancelled") {
        await step.do("mark-cancelled", () => updateJobStatus(env.DB, jobId, "cancelled"));
        return;
      }
      const outcome = result.outcome;
      if (outcome.kind === "continue") continue;

      await step.do("finalize", async () => {
        const now = nowIso();
        if (outcome.kind === "proposed") {
          // ArchiveService already published the snapshot.
          await env.DB.batch([
            env.DB.prepare(
              "UPDATE jobs SET status = 'succeeded', result = 'complete', updated_at = ? WHERE id = ?",
            ).bind(now, jobId),
            env.DB.prepare(
              "UPDATE capture_runs SET status = 'succeeded', outcome_code = 'complete', updated_at = ? WHERE id = ?",
            ).bind(now, run.id),
          ]);
        } else if (outcome.kind === "finished") {
          const o = outcome as { outcome: string; reasonCode: string };
          const jobStatus =
            o.outcome === "needs_input"
              ? "needs_input"
              : o.outcome === "complete" || o.outcome === "partial"
                ? "succeeded"
                : "failed";
          await env.DB.batch([
            env.DB.prepare(
              "UPDATE jobs SET status = ?, result = ?, updated_at = ? WHERE id = ?",
            ).bind(jobStatus, o.reasonCode, now, jobId),
            env.DB.prepare(
              "UPDATE capture_runs SET status = ?, outcome_code = ?, updated_at = ? WHERE id = ?",
            ).bind(jobStatus === "succeeded" ? "succeeded" : jobStatus, o.outcome, now, run.id),
          ]);
        } else {
          await env.DB.batch([
            env.DB.prepare(
              "UPDATE jobs SET status = 'failed', result = 'budget_exceeded', updated_at = ? WHERE id = ?",
            ).bind(now, jobId),
            env.DB.prepare(
              "UPDATE capture_runs SET status = 'failed', outcome_code = 'budget_exhausted', updated_at = ? WHERE id = ?",
            ).bind(now, run.id),
          ]);
        }
      });
      return;
      }
    } catch (err) {
      // A step that exhausts retries would otherwise leave the job 'running'
      // forever (observed: agent_events UNIQUE collision errored a run while
      // its job row stayed running). Mark it failed, then rethrow so the
      // workflow instance still reports the error.
      const now = nowIso();
      await env.DB.batch([
        env.DB.prepare(
          "UPDATE jobs SET status = 'failed', result = 'workflow_error', error = ?, updated_at = ? WHERE id = ?",
        ).bind(err instanceof Error ? err.message.slice(0, 500) : String(err), now, jobId),
        env.DB.prepare(
          "UPDATE capture_runs SET status = 'failed', outcome_code = 'failed', updated_at = ? WHERE id = ?",
        ).bind(now, run.id),
      ]);
      throw err;
    }
  }
}
