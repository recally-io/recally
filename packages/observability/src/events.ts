// Per-task audit fields (plan §18.1). Logs carry ids and refs, never bodies,
// full prompts, cookies, signed urls, or share tokens (plan §14.5).
export interface AuditEvent {
  request_id: string;
  job_id?: string;
  run_id?: string;
  attempt_id?: string;
  stage: string;
  tool_name?: string;
  input_hash?: string;
  output_ref?: string;
  reason_code?: string;
  duration_ms?: number;
  model_id?: string;
  usage?: { input_tokens: number; output_tokens: number };
  browser_ms?: number;
  policy_version?: string;
}

export function audit(event: AuditEvent): void {
  console.log(JSON.stringify({ level: "info", ...event }));
}
