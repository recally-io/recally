// Export the AuthClient class and singleton instance
export { AuthClient, authClient } from "./client";

// Export error classes
export { AuthError } from "./errors";

// Export types
export type {
	LoginInput,
	OAuthLoginResponse,
	RegisterInput,
	User,
	UserSettings,
} from "./types";
