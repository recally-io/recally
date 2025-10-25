/**
 * Custom error class for authentication-related errors
 */
export class AuthError extends Error {
	constructor(
		public code: string,
		public status: number,
		message: string,
	) {
		super(message);
		this.name = "AuthError";
		// Maintains proper stack trace for where our error was thrown (only available on V8)
		if (Error.captureStackTrace) {
			Error.captureStackTrace(this, AuthError);
		}
	}

	/**
	 * Creates an AuthError from an API response
	 * @param response - The error response from the API
	 * @returns A new AuthError instance
	 */
	static fromResponse(response: {
		code?: string;
		status?: number;
		message?: string;
	}): AuthError {
		return new AuthError(
			response.code || "AUTH_ERROR",
			response.status || 500,
			response.message || "Authentication failed",
		);
	}
}
