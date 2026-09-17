import { AppError } from "@recally/domain";
import { type CaptureRunRow, updateJobStatus } from "@recally/storage";
import { assembleRun } from "../../context";
import type { IngestEnv } from "../../env";

export async function commitManual(
  env: IngestEnv,
  run: CaptureRunRow,
  attemptId: string,
  jobId: string,
  manualText: string,
  url?: string,
) {
  const { deps, ctx } = await assembleRun(env, run, attemptId);

  const source = await deps.runStore.saveSource(ctx, {
    url: url ?? "manual:input",
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
}
