// Workflow dispatch adapter (§9.3). Instance IDs are deterministic so a retried
// dispatch reads the same instance instead of creating a duplicate.

export interface WorkflowInstanceRef {
  id: string;
}

export interface WorkflowBinding {
  create(options: { id: string; params?: unknown }): Promise<WorkflowInstanceRef>;
  get(id: string): Promise<{ status(): Promise<{ status: string }> }>;
}

export function workflowInstanceId(jobId: string, restartGeneration: number): string {
  return `job-${jobId}-r${restartGeneration}`;
}

export async function dispatchJob(
  binding: WorkflowBinding,
  jobId: string,
  restartGeneration: number,
  params: unknown,
): Promise<WorkflowInstanceRef> {
  const id = workflowInstanceId(jobId, restartGeneration);

  try {
    return await binding.create({ id, params });
  } catch (err) {
    // Duplicate id = already dispatched; report the existing instance.
    if (String(err).includes("already")) return { id };
    throw err;
  }
}
