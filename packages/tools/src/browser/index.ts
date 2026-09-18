import type { ToolOutcome, ToolSpec } from "@recally/agent-runtime";
import { checkUrlTarget, type ToolContext } from "@recally/capture";
import {
  browserActInput,
  browserCaptureInput,
  browserCloseInput,
  browserObserveInput,
  browserOpenInput,
} from "@recally/contracts";
import { AppError, nowIso } from "@recally/domain";
import type { BrowserFactory, BrowserSessionHandle, ToolDeps } from "../deps";

// browser_* tools (§6.3): atomic actions the agent picks; the session handle
// never crosses a step and is managed by the platform adapter. Per-episode
// state lives in a per-run session registry so consecutive calls share one
// session inside a browser episode (§9.5).

// Valid inside one Workflow step (one browser episode), isolated by run id.
const sessions = new Map<string, BrowserSessionHandle>();

type BrowserObservation = Awaited<ReturnType<BrowserSessionHandle["observe"]>>;

function requireSession(runId: string): BrowserSessionHandle {
  const session = sessions.get(runId);

  if (!session) throw new AppError("invalid_input", "no browser session — call browser_open");

  return session;
}

function formatBrowserObservation(obs: BrowserObservation, prefix: string, limit: number): string {
  return `${prefix}${obs.title || obs.url}\n${obs.text.slice(0, limit)}`;
}

async function saveBrowserObservation(
  deps: ToolDeps,
  ctx: ToolContext,
  obs: BrowserObservation,
  pageKind?: "unknown",
): Promise<ToolOutcome> {
  const observation = {
    url: obs.url,
    fetchedAt: nowIso(),
    textPreview: obs.text.slice(0, 2000),
    nodeRefs: obs.nodeRefs,
  };

  const resultRef = await deps.runStore.saveObservation(
    ctx,
    pageKind ? { ...observation, pageKind, title: obs.title } : observation,
  );

  return {
    content: formatBrowserObservation(obs, pageKind ? "browser open: " : "", 1500),
    details: { url: obs.url, nodeRefs: obs.nodeRefs },
    resultRef,
  };
}

async function saveBrowserCapture(
  deps: ToolDeps,
  ctx: ToolContext,
  capture: {
    body: string | Uint8Array;
    kind: "rendered_dom" | "screenshot";
    contentType: string;
    label: string;
  },
): Promise<ToolOutcome> {
  const { label, ...input } = capture;
  const source = await deps.runStore.saveSource(ctx, { url: "browser:session", ...input });

  return {
    content: `captured ${label} as source ${source.sourceId}`,
    resultRef: source.evidenceRef,
  };
}

function requireBrowser(deps: ToolDeps): BrowserFactory {
  if (!deps.browser) {
    throw new AppError("not_implemented", "browser capability not configured");
  }

  return deps.browser;
}

export function browserOpenTool(deps: ToolDeps): ToolSpec {
  return {
    name: "browser_open",
    label: "Open browser",
    description:
      "Open a controlled browser session on a URL you already observed. Expensive — only when fetched/adapted evidence has a concrete gap a browser could reach.",
    schema: browserOpenInput,
    async execute(ctx, args) {
      // ctx.budget is rebuilt per durable turn, so enforce against persisted
      // tool_call events — that count survives step/isolate boundaries.
      const episodes = await deps.runStore.countToolCalls(ctx, "browser_open");

      if (episodes >= ctx.budget.browserEpisodesLeft) {
        throw new AppError("budget_exceeded", "browser episode budget exhausted");
      }

      const { url_ref } = browserOpenInput.parse(args);
      checkUrlTarget(url_ref);
      const session = await requireBrowser(deps).open(url_ref);
      sessions.set(ctx.runId, session);
      const obs = await session.observe();

      return saveBrowserObservation(deps, ctx, obs, "unknown");
    },
  };
}

export function browserObserveTool(deps: ToolDeps): ToolSpec {
  return {
    name: "browser_observe",
    label: "Observe page",
    description:
      "Re-observe the current page: title, url, text summary, node refs. Old refs die after any act.",
    schema: browserObserveInput,
    async execute(ctx) {
      const session = requireSession(ctx.runId);
      const obs = await session.observe();

      return saveBrowserObservation(deps, ctx, obs);
    },
  };
}

export function browserActTool(deps: ToolDeps): ToolSpec {
  return {
    name: "browser_act",
    label: "Browser action",
    description:
      "One reading-scoped action: wait_for a selector, scroll, expand (click) an observed target, or navigate to an observed same-article link. Never login/post/purchase/fill forms.",
    schema: browserActInput,
    async execute(ctx, args) {
      const actions = await deps.runStore.countToolCalls(ctx, "browser_act");

      if (actions >= ctx.budget.browserActionsLeft) {
        throw new AppError("budget_exceeded", "browser action budget exhausted");
      }

      const session = requireSession(ctx.runId);
      const { action, target } = browserActInput.parse(args);

      if (action === "scroll") {
        await session.scroll();
      } else {
        const targetErrors = {
          wait_for: "wait_for requires target selector",
          expand: "expand requires an observed target",
          navigate_same_article: "navigate requires an observed link",
        };

        if (!target) throw new AppError("invalid_input", targetErrors[action]);

        switch (action) {
          case "wait_for":
            await session.waitFor(target);
            break;
          case "expand":
            await session.click(target);
            break;
          case "navigate_same_article":
            checkUrlTarget(target);
            await session.navigate(target);
            break;
        }
      }

      const obs = await session.observe();

      return {
        content: formatBrowserObservation(obs, `acted ${action}: now `, 1000),
        details: { action, url: obs.url },
      };
    },
  };
}

export function browserCaptureTool(deps: ToolDeps): ToolSpec {
  return {
    name: "browser_capture",
    label: "Capture page evidence",
    description:
      "Persist current rendered DOM, a screenshot, or text as evidence. Returns a source ref you can propose for archive.",
    schema: browserCaptureInput,
    async execute(ctx, args) {
      const session = requireSession(ctx.runId);
      const { mode } = browserCaptureInput.parse(args);

      if (mode === "rendered_dom" || mode === "text") {
        return saveBrowserCapture(deps, ctx, {
          body: await session.renderedDom(),
          kind: "rendered_dom",
          contentType: "text/html",
          label: "rendered DOM",
        });
      }

      return saveBrowserCapture(deps, ctx, {
        body: await session.screenshot(),
        kind: "screenshot",
        contentType: "image/webp",
        label: "screenshot",
      });
    },
  };
}

export function browserCloseTool(): ToolSpec {
  return {
    name: "browser_close",
    label: "Close browser",
    description: "End the browser episode. Always close when done.",
    schema: browserCloseInput,
    async execute(ctx) {
      const session = sessions.get(ctx.runId);
      sessions.delete(ctx.runId);

      if (session) await session.close();

      return { content: "browser session closed" };
    },
  };
}

export function browserTools(deps: ToolDeps): ToolSpec[] {
  return [
    browserOpenTool(deps),
    browserObserveTool(deps),
    browserActTool(deps),
    browserCaptureTool(deps),
    browserCloseTool(),
  ];
}
