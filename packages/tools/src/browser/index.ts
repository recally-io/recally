import type { ToolSpec } from "@recally/agent-runtime";
import { checkUrlTarget } from "@recally/capture";
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

interface SessionRegistry {
  get(runId: string): BrowserSessionHandle | null;
  set(runId: string, s: BrowserSessionHandle): void;
  delete(runId: string): void;
}

// In-memory registry: valid inside one Workflow step, which is exactly the
// lifetime of a browser episode. Never reused across runs.
const sessions: SessionRegistry = (() => {
  const map = new Map<string, BrowserSessionHandle>();
  return {
    get: (id) => map.get(id) ?? null,
    set: (id, s) => void map.set(id, s),
    delete: (id) => void map.delete(id),
  };
})();

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
      const ref = await deps.runStore.saveObservation(ctx, {
        url: obs.url,
        fetchedAt: nowIso(),
        pageKind: "unknown",
        title: obs.title,
        textPreview: obs.text.slice(0, 2000),
        nodeRefs: obs.nodeRefs,
      });
      return {
        content: `browser open: ${obs.title || obs.url}\n${obs.text.slice(0, 1500)}`,
        details: { url: obs.url, nodeRefs: obs.nodeRefs },
        resultRef: ref,
      };
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
      const session = sessions.get(ctx.runId);
      if (!session) throw new AppError("invalid_input", "no browser session — call browser_open");
      const obs = await session.observe();
      const ref = await deps.runStore.saveObservation(ctx, {
        url: obs.url,
        fetchedAt: nowIso(),
        textPreview: obs.text.slice(0, 2000),
        nodeRefs: obs.nodeRefs,
      });
      return {
        content: `${obs.title || obs.url}\n${obs.text.slice(0, 1500)}`,
        details: { url: obs.url, nodeRefs: obs.nodeRefs },
        resultRef: ref,
      };
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
      const session = sessions.get(ctx.runId);
      if (!session) throw new AppError("invalid_input", "no browser session — call browser_open");
      const { action, target } = browserActInput.parse(args);
      switch (action) {
        case "scroll":
          await session.scroll();
          break;
        case "wait_for":
          if (!target) throw new AppError("invalid_input", "wait_for requires target selector");
          await session.waitFor(target);
          break;
        case "expand":
          if (!target) throw new AppError("invalid_input", "expand requires an observed target");
          await session.click(target);
          break;
        case "navigate_same_article":
          if (!target) throw new AppError("invalid_input", "navigate requires an observed link");
          checkUrlTarget(target);
          await session.navigate(target);
          break;
      }
      const obs = await session.observe();
      return {
        content: `acted ${action}: now ${obs.title || obs.url}\n${obs.text.slice(0, 1000)}`,
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
      const session = sessions.get(ctx.runId);
      if (!session) throw new AppError("invalid_input", "no browser session — call browser_open");
      const { mode } = browserCaptureInput.parse(args);
      if (mode === "rendered_dom" || mode === "text") {
        const dom = await session.renderedDom();
        const source = await deps.runStore.saveSource(ctx, {
          url: "browser:session",
          kind: "rendered_dom",
          contentType: "text/html",
          body: dom,
        });
        return {
          content: `captured rendered DOM as source ${source.sourceId}`,
          resultRef: source.evidenceRef,
        };
      }
      const shot = await session.screenshot();
      const source = await deps.runStore.saveSource(ctx, {
        url: "browser:session",
        kind: "screenshot",
        contentType: "image/webp",
        body: shot,
      });
      return {
        content: `captured screenshot as source ${source.sourceId}`,
        resultRef: source.evidenceRef,
      };
    },
  };
}

export function browserCloseTool(_deps: ToolDeps): ToolSpec {
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
    browserCloseTool(deps),
  ];
}
