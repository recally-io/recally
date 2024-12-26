export const ROUTES = {
	HOME: "/",
	LOGIN: "/auth.html?mode=login",
	SIGNUP: "/auth.html?mode=register",
	SETTINGS: "/settings.html",
	BOOKMARKS: "/bookmarks.html",
} as const;

// Type for route keys
export type RouteKeys = keyof typeof ROUTES;

// Helper function to get path
export const getPath = (path: RouteKeys): string => ROUTES[path];
