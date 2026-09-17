// Error protocol per plan §12.5: { code, message, request_id, retryable, next_action }.
export const ERROR_CODES = [
  "needs_login",
  "forbidden",
  "not_found",
  "challenge_detected",
  "unsupported_type",
  "incomplete_content",
  "policy_denied",
  "budget_exceeded",
  "model_unavailable",
  "source_unavailable",
  "stale_generation",
  "cancelled",
  "conflict",
  "invalid_input",
  "not_implemented",
  "write_barrier",
  "internal",
] as const;

export type ErrorCode = (typeof ERROR_CODES)[number];

const HTTP_STATUS: Record<ErrorCode, number> = {
  needs_login: 401,
  forbidden: 403,
  not_found: 404,
  challenge_detected: 422,
  unsupported_type: 415,
  incomplete_content: 422,
  policy_denied: 403,
  budget_exceeded: 429,
  model_unavailable: 503,
  source_unavailable: 502,
  stale_generation: 409,
  cancelled: 409,
  conflict: 409,
  invalid_input: 400,
  not_implemented: 501,
  write_barrier: 503,
  internal: 500,
};

export class AppError extends Error {
  readonly code: ErrorCode;
  readonly httpStatus: number;
  readonly retryable: boolean;
  readonly nextAction?: string;

  constructor(
    code: ErrorCode,
    message: string,
    opts: { retryable?: boolean; nextAction?: string } = {},
  ) {
    super(message);
    this.name = "AppError";
    this.code = code;
    this.httpStatus = HTTP_STATUS[code];
    this.retryable = opts.retryable ?? false;

    if (opts.nextAction !== undefined) this.nextAction = opts.nextAction;
  }
}

export interface ErrorResponse {
  code: ErrorCode;
  message: string;
  request_id: string;
  retryable: boolean;
  next_action?: string;
}

export function toErrorResponse(
  err: unknown,
  requestId: string,
): { status: number; body: ErrorResponse } {
  if (err instanceof AppError) {
    const body: ErrorResponse = {
      code: err.code,
      message: err.message,
      request_id: requestId,
      retryable: err.retryable,
    };

    if (err.nextAction !== undefined) body.next_action = err.nextAction;

    return { status: err.httpStatus, body };
  }

  return {
    status: 500,
    body: {
      code: "internal",
      message: err instanceof Error ? err.message : "internal error",
      request_id: requestId,
      retryable: false,
    },
  };
}
