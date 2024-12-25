export const ROUTES = {
  HOME: '/',
  DASHBOARD: '/dashboard',
  LOGIN: '/auth/login',
  SIGNUP: '/auth/register',
  PROFILE: '/profile',
  SETTINGS: '/settings',
  BOOKMARKS: '/bookmarks',
} as const;

// Type for route keys
export type RouteKeys = keyof typeof ROUTES;

// Helper function to get path
export const getPath = (path: RouteKeys): string => ROUTES[path];

