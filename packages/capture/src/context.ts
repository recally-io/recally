// Injected server-side on every tool call (§6.1). The model can never switch
// library, item, run, skill, or budget via tool arguments.
export interface ToolContext {
  libraryId: string;
  itemId: string;
  runId: string;
  attemptId: string;
  generation: number;
  skillRevision: string;
  toolsetVersion: string;
  policyVersion: string;
  budget: CaptureBudgets;
}

export interface CaptureBudgets {
  turnsLeft: number;
  modelCallsLeft: number;
  inputTokensLeft: number;
  outputTokensLeft: number;
  browserEpisodesLeft: number;
  browserMsLeft: number;
  browserActionsLeft: number;
  downloadBytesLeft: number;
  wallClockDeadlineMs: number;
}

// Internal run bookkeeping (§6.5) — runtime state, not a strategy enum.
export type CaptureRunState = {
  status: "running" | "needs_input" | "committing" | "done" | "failed";
  piStateRef: string | null;
  skillRevision: string;
  toolsetVersion: string;
  observations: string[]; // evidence refs
  activeBrowserEpisode?: string;
  budgets: CaptureBudgets;
  noProgressCount: number;
};
