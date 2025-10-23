import { AuthError } from "./errors";
import type {
	LoginInput,
	OAuthLoginResponse,
	RegisterInput,
	User,
} from "./types";

/**
 * Client for authentication operations
 * Handles login, registration, logout, session validation, and OAuth flows
 */
export class AuthClient {
	constructor(private baseURL = "/api/v1") {}

	/**
	 * Authenticate user with email and password
	 * @param email - User's email address
	 * @param password - User's password
	 * @returns Authenticated user object
	 */
	async login(email: string, password: string): Promise<User> {
		return this.post<User>("/auth/login", { email, password });
	}

	/**
	 * Register a new user account
	 * @param username - Desired username
	 * @param email - User's email address
	 * @param password - User's password
	 * @returns Newly created user object
	 */
	async register(
		username: string,
		email: string,
		password: string,
	): Promise<User> {
		return this.post<User>("/auth/register", { username, email, password });
	}

	/**
	 * Log out the current user
	 * Clears the authentication cookie on the server
	 */
	async logout(): Promise<void> {
		await this.post<void>("/auth/logout", {});
	}

	/**
	 * Validate the current JWT session
	 * @returns Current user object if session is valid
	 * @throws AuthError if session is invalid or expired
	 */
	async validateSession(): Promise<User> {
		return this.get<User>("/auth/validate-jwt");
	}

	/**
	 * Get OAuth login URL for a provider
	 * @param provider - OAuth provider name (e.g., 'github', 'google')
	 * @returns Object containing the OAuth login URL
	 */
	async getOAuthURL(provider: string): Promise<OAuthLoginResponse> {
		return this.get<OAuthLoginResponse>(
			`/oauth/${provider.toLowerCase()}/login`,
		);
	}

	/**
	 * Redirect to OAuth login page for a provider
	 * @param provider - OAuth provider name (e.g., 'github', 'google')
	 */
	redirectToOAuth(provider: string): void {
		window.location.href = `${this.baseURL}/oauth/${provider.toLowerCase()}/login`;
	}

	/**
	 * Perform a GET request
	 * @param path - API endpoint path
	 * @returns Typed response data
	 */
	private async get<T>(path: string): Promise<T> {
		const response = await fetch(`${this.baseURL}${path}`, {
			credentials: "include",
			headers: {
				"Content-Type": "application/json",
			},
		});

		return this.handleResponse<T>(response);
	}

	/**
	 * Perform a POST request
	 * @param path - API endpoint path
	 * @param body - Request body
	 * @returns Typed response data
	 */
	private async post<T>(path: string, body: unknown): Promise<T> {
		const response = await fetch(`${this.baseURL}${path}`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			credentials: "include",
			body: JSON.stringify(body),
		});

		return this.handleResponse<T>(response);
	}

	/**
	 * Handle the fetch response, extracting data or throwing errors
	 * Matches the existing fetcher pattern in the codebase
	 * @param response - Fetch response object
	 * @returns Parsed response data
	 */
	private async handleResponse<T>(response: Response): Promise<T> {
		if (!response.ok) {
			let errorInfo: { code?: string; status?: number; message?: string };

			try {
				errorInfo = await response.json();
			} catch {
				errorInfo = {
					status: response.status,
					message: await response.text(),
				};
			}

			throw AuthError.fromResponse({
				...errorInfo,
				status: response.status,
			});
		}

		// Handle 204 No Content
		if (response.status === 204) {
			return {} as T;
		}

		// Parse response and extract data field (matches existing fetcher pattern)
		const result = await response.json();
		return result.data;
	}
}

// Singleton instance for use throughout the application
export const authClient = new AuthClient();
