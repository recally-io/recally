import type { R2EvidenceStore } from "@recally/platform-cloudflare";
import { newId, nowIso } from "@recally/domain";
import { type CaptureRunRow, getJob, r2Keys } from "@recally/storage";
import { captureTools } from "@recally/tools";
import { assembleRun } from "../../context";
import type { IngestEnv } from "../../env";
import type { WorkflowInput } from "../common";

interface CaptureTurnInput extends WorkflowInput {
  run: CaptureRunRow;
  attemptId: string;
  turn: number;
  maxTurns: number;
}

export async function runCaptureTurn(
  env: IngestEnv,
  evidence: R2EvidenceStore,
  { jobId, libraryId, run, attemptId, turn, maxTurns }: CaptureTurnInput,
) {
  const db = env.DB;
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
}
