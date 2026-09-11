import { sha256Hex } from "./hash";

export interface OperationKeyInput {
  libraryId: string;
  inputRevision: string;
  stage: string;
  pipelineVersion: string;
  modelConfigVersion: string;
  stageConfigVersion: string;
}

// Dedupe key per logical operation (plan §9.6). Uniqueness prevents duplicate
// effects; it does not imply the external call was billed exactly once.
export async function operationKey(input: OperationKeyInput): Promise<string> {
  return sha256Hex(
    [
      input.libraryId,
      input.inputRevision,
      input.stage,
      input.pipelineVersion,
      input.modelConfigVersion,
      input.stageConfigVersion,
    ].join("|"),
  );
}
