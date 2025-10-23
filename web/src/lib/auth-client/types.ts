// Re-export User and UserSettings from the users module to maintain consistency
export type { User, UserSettings } from "../apis/users";

// Login input interface
export interface LoginInput {
	email: string;
	password: string;
}

// Register input interface
export interface RegisterInput {
	username: string;
	email: string;
	password: string;
}

// OAuth response interface
export interface OAuthLoginResponse {
	url: string;
}
