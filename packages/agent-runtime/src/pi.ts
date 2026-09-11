import type { AgentMessage, AgentTool, StreamFn } from "@earendil-works/pi-agent-core";
import { Agent } from "@earendil-works/pi-agent-core";
import type { Api, Model } from "@earendil-works/pi-ai";
import type { ToolContext } from "@recally/capture";
import { Type } from "typebox";
import { z } from "zod";
import type {
  AgentRuntime,
  RunOutcome,
  ToolOutcome,
  ToolSpec,
  TurnInput,
  TurnResult,
} from "./spec";

// Pi adapter (§5.2, ADR-02/06). Boundary verified by M0-T02: bundle size,
// compat flags, and state serialize/restore on a real Worker. Our code never
// sees Pi internals beyond this file.

interface PersistedState {
  messages: unknown[]; // Pi AgentMessage[] — serialized verbatim
}

export interface PiRuntimeOptions {
  streamFn: StreamFn;
  resolveModel: (modelId: string) => Model<Api>;
  runtimeVersion?: string;
}

interface Collected {
  events: TurnResult["events"];
  endRun?: Extract<RunOutcome, { kind: "finished" } | { kind: "proposed" }>;
}

function toPiTool(spec: ToolSpec, ctx: ToolContext, collected: Collected): AgentTool {
  return {
    name: spec.name,
    label: spec.label,
    description: spec.description,
    // Type.Unsafe carries the raw JSON Schema; our execute() re-validates with
    // zod regardless, so a model can never reach a side effect unvalidated.
    parameters: Type.Unsafe(z.toJSONSchema(spec.schema)),
    executionMode: "sequential",
    execute: async (
      _toolCallId,
      params,
    ): Promise<import("@earendil-works/pi-agent-core").AgentToolResult<unknown>> => {
      collected.events.push({ kind: "tool_call", toolName: spec.name });
      let outcome: ToolOutcome;
      try {
        const args = spec.schema.parse(params);
        outcome = await spec.execute(ctx, args);
      } catch (err) {
        // ZodError.message is a JSON blob the model can't act on — flatten
        // issues into `path: message` so a bad-args retry can self-correct.
        const message =
          err instanceof z.ZodError
            ? err.issues.map((i) => `${i.path.join(".") || "input"}: ${i.message}`).join("; ")
            : err instanceof Error
              ? err.message
              : String(err);
        collected.events.push({
          kind: "tool_result",
          toolName: spec.name,
          reasonCode: "error",
          summary: `${spec.name} failed: ${message}`.slice(0, 500),
        });
        throw err instanceof z.ZodError
          ? new Error(`${spec.name} invalid arguments: ${message}`)
          : err;
      }
      collected.events.push({
        kind: "tool_result",
        toolName: spec.name,
        summary: outcome.content.slice(0, 500),
        ...(outcome.resultRef !== undefined ? { evidenceRef: outcome.resultRef } : {}),
      });
      if (outcome.terminate && outcome.runOutcome) {
        collected.endRun = outcome.runOutcome;
      }
      return {
        content: [{ type: "text" as const, text: outcome.content }],
        details: outcome.details,
        ...(outcome.terminate !== undefined ? { terminate: outcome.terminate } : {}),
      };
    },
  };
}

// One durable turn = restore Pi state → one model turn (with tool execution)
// → serialize state. The workflow wraps each call in step.do("turn-N"), so a
// crash between turns replays at most one turn (plan §9.4).
export class PiRuntime implements AgentRuntime {
  readonly runtimeVersion: string;
  private readonly streamFn: StreamFn;
  private readonly resolveModel: (modelId: string) => Model<Api>;

  constructor(opts: PiRuntimeOptions) {
    this.streamFn = opts.streamFn;
    this.resolveModel = opts.resolveModel;
    this.runtimeVersion = opts.runtimeVersion ?? "pi-0.85.1";
  }

  async runTurn(input: TurnInput): Promise<TurnResult> {
    if (input.turn >= input.maxTurns) {
      const state = input.serializedState ?? JSON.stringify({ messages: [] });
      return {
        outcome: { kind: "budget_exhausted" },
        serializedState: state,
        events: [],
      };
    }

    const collected: Collected = { events: [] };
    const restored: PersistedState = input.serializedState
      ? (JSON.parse(input.serializedState) as PersistedState)
      : { messages: [] };

    const agent = new Agent({
      streamFn: this.streamFn,
      initialState: {
        systemPrompt: input.systemPrompt,
        model: this.resolveModel(input.model),
        tools: input.tools.map((t) => toPiTool(t, input.ctx, collected)),
        messages: restored.messages as AgentMessage[],
      },
      // One model turn per durable step. Tool calls inside the turn still run
      // through our wrapped execute() — evidence lands before results return.
      shouldStopAfterTurn: () => true,
      // Policy hook: schema/permission violations block before side effects.
      beforeToolCall: async ({ toolCall }) => {
        if (!input.tools.some((t) => t.name === toolCall.name)) {
          return {
            block: true,
            reason: `tool ${toolCall.name} is not available to this agent`,
          };
        }
        return undefined;
      },
    });

    agent.subscribe((event) => {
      if (event.type === "message_end" || event.type === "turn_end") {
        collected.events.push({ kind: "decision" });
      }
    });

    try {
      const last = restored.messages[restored.messages.length - 1] as { role?: string } | undefined;
      if (restored.messages.length === 0) {
        await agent.prompt(input.goal);
      } else if (last?.role === "assistant") {
        // The model ended its turn talking instead of calling a tool — a
        // continue() can't resume from an assistant message, so nudge it.
        await agent.prompt(
          "Continue the task. Use tools to make progress; call propose_archive or finish when done.",
        );
      } else {
        agent.state.messages = restored.messages as AgentMessage[];
        await agent.continue();
      }
      await agent.waitForIdle();
    } catch (err) {
      collected.events.push({
        kind: "tool_result",
        summary: `agent error: ${err instanceof Error ? err.message : String(err)}`,
      });
      return {
        outcome: {
          kind: "finished",
          outcome: "failed",
          reasonCode: "agent_error",
        },
        serializedState: JSON.stringify({ messages: agent.state.messages }),
        events: collected.events,
      };
    }

    const serializedState = JSON.stringify({ messages: agent.state.messages });
    const outcome: RunOutcome = collected.endRun ?? { kind: "continue" };
    return { outcome, serializedState, events: collected.events };
  }
}
