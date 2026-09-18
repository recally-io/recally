import type { JobStatus } from "@recally/domain";

export function isJobLive(status: JobStatus): boolean {
  return status === "queued" || status === "running" || status === "retry_wait";
}

export function isJobCancellable(status: JobStatus): boolean {
  return isJobLive(status) || status === "needs_input";
}
