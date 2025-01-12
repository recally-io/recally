export const ROUTES = {
	HOME: "/",

	AUTH: "/auth",
	AUTH_LOGIN: "/auth/login",
	AUTH_REGISTER: "/auth/register",
	AUTH_CALLBACK: "/auth/$provider/callback",

	BOOKMARKS: "/bookmarks",
	BOOKMARK_DETAIL: "/bookmarks/$id",

	SETTINGS: "/settings",
	SETTINGS_PROFILE: "/settings/profile",
	SETTINGS_AI: "/settings/ai",
	SETTINGS_API_KEYS: "/settings/api-keys",
} as const;

// Type for route keys
export type RouteKeys = keyof typeof ROUTES;
