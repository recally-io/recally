import type { ToolContext } from "@recally/capture";
import type { z } from "zod";

// Our tool contract (§6). Tools expose capability only — they never chain
// fetch→browser→anything internally, never decide capture success, and never
// see the model's context. Platform adapters supply execute().

export interface ToolOutcome {
  // Bounded text the model sees; full evidence was already persisted.
  content: string;
  // Structured observation persisted to agent_events; kept small.
  details?: unknown;
  // Evidence object ref (R2 key) when the tool produced content.
  resultRef?: string;
  // finish/propose_archive end the agent loop.
  terminate?: boolean;
  // Semantic end-of-run, only meaningful with terminate: true. A rejected
  // archive proposal does NOT end the run — it returns feedback instead.
  runOutcome?:
    | { kind: "finished"; outcome: string; reasonCode: string }
    | { kind: "proposed"; proposalRef: string };
}

export interface ToolSpec {
  name: string;
  label: string;
  description: string;
  schema: z.ZodType<unknown>;
  // Server-side injected identity/budget — model args can never override it.
  execute(ctx: ToolContext, args: unknown): Promise<ToolOutcome>;
}

// What a workflow needs from the runtime per turn: run the model once,
// execute the resulting tool calls, hand back an updated serializable state
// and whether the run should continue.
export interface TurnInput {
  serializedState: string | null; // null = first turn
  systemPrompt: string; // skill body + goal framing
  goal: string;
  tools: ToolSpec[];
  ctx: ToolContext;
  model: string;
  maxTurns: number;
  turn: number;
}

export type RunOutcome =
  | { kind: "continue" } // more turns needed
  | { kind: "finished"; outcome: string; reasonCode: string }
  | { kind: "proposed"; proposalRef: string }
  | { kind: "budget_exhausted" };

export interface TurnResult {
  outcome: RunOutcome;
  serializedState: string;
  // Audit records appended this turn (decision/tool_call/tool_result).
  events: Array<{
    kind: "decision" | "tool_call" | "tool_result" | "observation";
    toolName?: string;
    reasonCode?: string;
    summary?: string;
    evidenceRef?: string;
    usage?: { input_tokens: number; output_tokens: number };
  }>;
}

export interface AgentRuntime {
  readonly runtimeVersion: string;
  runTurn(input: TurnInput): Promise<TurnResult>;
}
