import type { ToolSpec } from "@recally/agent-runtime";
import { checkUrlTarget } from "@recally/capture";
import { archiveAssetInput, archiveProposalInput, finishInput } from "@recally/contracts";
import { AppError } from "@recally/domain";
import type { ToolDeps } from "../deps";

// archive_asset (§6.4): fetch an already-observed asset into evidence storage.
// Policy re-checks the target — "observed" is not a network bypass.
export function archiveAssetTool(deps: ToolDeps): ToolSpec {
  return {
    name: "archive_asset",
    label: "Archive asset",
    description:
      "Download a key resource (in-article image, media) you observed into the archive. Policy limits on type/size/redirects still apply.",
    schema: archiveAssetInput,
    async execute(ctx, args) {
      const { asset_ref } = archiveAssetInput.parse(args);
      const url = checkUrlTarget(asset_ref);
      const { safeFetch } = await import("../web-fetch/index.js");
      const fetched = await safeFetch(url.toString(), deps.fetchFn ?? fetch, {
        maxBytes: 5 * 1024 * 1024,
      });
      const source = await deps.runStore.saveSource(ctx, {
        url: fetched.finalUrl,
        kind: "response_body",
        contentType: fetched.contentType,
        body: fetched.body,
      });
      return {
        content: `asset archived: ${source.sourceId} (${source.byteSize} bytes, ${fetched.contentType})`,
        details: { asset: source.sourceId },
        resultRef: source.evidenceRef,
      };
    },
  };
}

// propose_archive: the only path to publishing. Rejection returns reasons as
// feedback so the agent can revise or finish partial.
export function proposeArchiveTool(deps: ToolDeps): ToolSpec {
  return {
    name: "propose_archive",
    label: "Propose archive",
    description:
      "Submit your archive proposal: which sources and block ranges are the real content, which assets are required, what is missing. sourceIds must be ids from earlier tool results (e.g. the 'source <uuid>' in a capture result), and block ranges use block ids from read_source/extract_content (e.g. 'b1'). adapter_record and manual sources archive whole-body — pass no block ranges for them. The archive service validates and may reject — revise or finish partial instead of retrying the same proposal.",
    schema: archiveProposalInput,
    async execute(ctx, args) {
      if (!deps.committer) throw new AppError("internal", "committer not configured");
      const proposal = archiveProposalInput.parse(args);
      const result = await deps.committer.commit(ctx, proposal);
      if (result.status === "accepted") {
        return {
          content: `archive accepted: snapshot ${result.snapshotId}`,
          terminate: true,
          runOutcome: { kind: "proposed", proposalRef: result.snapshotId! },
          details: result,
        };
      }
      return {
        content: `archive ${result.status}: ${(result.problems ?? []).join("; ")}`,
        details: result,
      };
    },
  };
}

export function finishTool(_deps: ToolDeps): ToolSpec {
  return {
    name: "finish",
    label: "Finish run",
    description:
      "End this capture run with a structured outcome: complete | partial | needs_input | failed | budget_exhausted, plus a reason code and evidence refs.",
    schema: finishInput,
    async execute(_ctx, args) {
      const f = finishInput.parse(args);
      return {
        content: `run finished: ${f.outcome} (${f.reason_code})`,
        terminate: true,
        runOutcome: {
          kind: "finished",
          outcome: f.outcome,
          reasonCode: f.reason_code,
        },
        details: f,
      };
    },
  };
}

export function archiveTools(deps: ToolDeps): ToolSpec[] {
  return [archiveAssetTool(deps), proposeArchiveTool(deps), finishTool(deps)];
}
