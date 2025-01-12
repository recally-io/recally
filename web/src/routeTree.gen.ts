/* eslint-disable */

// @ts-nocheck

// noinspection JSUnusedGlobalSymbols

// This file was automatically generated by TanStack Router.
// You should NOT make any changes in this file as it will be overwritten.
// Additionally, you should also exclude this file from your linter and/or formatter to prevent it from being checked or modified.

import { createFileRoute } from "@tanstack/react-router";

// Import Routes

import { Route as rootRoute } from "./routes/__root";
import { Route as SettingsIndexImport } from "./routes/settings/index";
import { Route as BookmarksIndexImport } from "./routes/bookmarks/index";
import { Route as AuthIndexImport } from "./routes/auth/index";
import { Route as ShareIdImport } from "./routes/share/$id";
import { Route as SettingsProfileImport } from "./routes/settings/profile";
import { Route as SettingsApiKeysImport } from "./routes/settings/api-keys";
import { Route as SettingsAiImport } from "./routes/settings/ai";
import { Route as BookmarksIdImport } from "./routes/bookmarks/$id";
import { Route as AuthRegisterImport } from "./routes/auth/register";
import { Route as AuthLoginImport } from "./routes/auth/login";
import { Route as AuthOauthProviderCallbackImport } from "./routes/auth/oauth.$provider.callback";

// Create Virtual Routes

const IndexLazyImport = createFileRoute("/")();

// Create/Update Routes

const IndexLazyRoute = IndexLazyImport.update({
	id: "/",
	path: "/",
	getParentRoute: () => rootRoute,
} as any).lazy(() => import("./routes/index.lazy").then((d) => d.Route));

const SettingsIndexRoute = SettingsIndexImport.update({
	id: "/settings/",
	path: "/settings/",
	getParentRoute: () => rootRoute,
} as any);

const BookmarksIndexRoute = BookmarksIndexImport.update({
	id: "/bookmarks/",
	path: "/bookmarks/",
	getParentRoute: () => rootRoute,
} as any);

const AuthIndexRoute = AuthIndexImport.update({
	id: "/auth/",
	path: "/auth/",
	getParentRoute: () => rootRoute,
} as any);

const ShareIdRoute = ShareIdImport.update({
	id: "/share/$id",
	path: "/share/$id",
	getParentRoute: () => rootRoute,
} as any);

const SettingsProfileRoute = SettingsProfileImport.update({
	id: "/settings/profile",
	path: "/settings/profile",
	getParentRoute: () => rootRoute,
} as any);

const SettingsApiKeysRoute = SettingsApiKeysImport.update({
	id: "/settings/api-keys",
	path: "/settings/api-keys",
	getParentRoute: () => rootRoute,
} as any);

const SettingsAiRoute = SettingsAiImport.update({
	id: "/settings/ai",
	path: "/settings/ai",
	getParentRoute: () => rootRoute,
} as any);

const BookmarksIdRoute = BookmarksIdImport.update({
	id: "/bookmarks/$id",
	path: "/bookmarks/$id",
	getParentRoute: () => rootRoute,
} as any);

const AuthRegisterRoute = AuthRegisterImport.update({
	id: "/auth/register",
	path: "/auth/register",
	getParentRoute: () => rootRoute,
} as any);

const AuthLoginRoute = AuthLoginImport.update({
	id: "/auth/login",
	path: "/auth/login",
	getParentRoute: () => rootRoute,
} as any);

const AuthOauthProviderCallbackRoute = AuthOauthProviderCallbackImport.update({
	id: "/auth/oauth/$provider/callback",
	path: "/auth/oauth/$provider/callback",
	getParentRoute: () => rootRoute,
} as any);

// Populate the FileRoutesByPath interface

declare module "@tanstack/react-router" {
	interface FileRoutesByPath {
		"/": {
			id: "/";
			path: "/";
			fullPath: "/";
			preLoaderRoute: typeof IndexLazyImport;
			parentRoute: typeof rootRoute;
		};
		"/auth/login": {
			id: "/auth/login";
			path: "/auth/login";
			fullPath: "/auth/login";
			preLoaderRoute: typeof AuthLoginImport;
			parentRoute: typeof rootRoute;
		};
		"/auth/register": {
			id: "/auth/register";
			path: "/auth/register";
			fullPath: "/auth/register";
			preLoaderRoute: typeof AuthRegisterImport;
			parentRoute: typeof rootRoute;
		};
		"/bookmarks/$id": {
			id: "/bookmarks/$id";
			path: "/bookmarks/$id";
			fullPath: "/bookmarks/$id";
			preLoaderRoute: typeof BookmarksIdImport;
			parentRoute: typeof rootRoute;
		};
		"/settings/ai": {
			id: "/settings/ai";
			path: "/settings/ai";
			fullPath: "/settings/ai";
			preLoaderRoute: typeof SettingsAiImport;
			parentRoute: typeof rootRoute;
		};
		"/settings/api-keys": {
			id: "/settings/api-keys";
			path: "/settings/api-keys";
			fullPath: "/settings/api-keys";
			preLoaderRoute: typeof SettingsApiKeysImport;
			parentRoute: typeof rootRoute;
		};
		"/settings/profile": {
			id: "/settings/profile";
			path: "/settings/profile";
			fullPath: "/settings/profile";
			preLoaderRoute: typeof SettingsProfileImport;
			parentRoute: typeof rootRoute;
		};
		"/share/$id": {
			id: "/share/$id";
			path: "/share/$id";
			fullPath: "/share/$id";
			preLoaderRoute: typeof ShareIdImport;
			parentRoute: typeof rootRoute;
		};
		"/auth/": {
			id: "/auth/";
			path: "/auth";
			fullPath: "/auth";
			preLoaderRoute: typeof AuthIndexImport;
			parentRoute: typeof rootRoute;
		};
		"/bookmarks/": {
			id: "/bookmarks/";
			path: "/bookmarks";
			fullPath: "/bookmarks";
			preLoaderRoute: typeof BookmarksIndexImport;
			parentRoute: typeof rootRoute;
		};
		"/settings/": {
			id: "/settings/";
			path: "/settings";
			fullPath: "/settings";
			preLoaderRoute: typeof SettingsIndexImport;
			parentRoute: typeof rootRoute;
		};
		"/auth/oauth/$provider/callback": {
			id: "/auth/oauth/$provider/callback";
			path: "/auth/oauth/$provider/callback";
			fullPath: "/auth/oauth/$provider/callback";
			preLoaderRoute: typeof AuthOauthProviderCallbackImport;
			parentRoute: typeof rootRoute;
		};
	}
}

// Create and export the route tree

export interface FileRoutesByFullPath {
	"/": typeof IndexLazyRoute;
	"/auth/login": typeof AuthLoginRoute;
	"/auth/register": typeof AuthRegisterRoute;
	"/bookmarks/$id": typeof BookmarksIdRoute;
	"/settings/ai": typeof SettingsAiRoute;
	"/settings/api-keys": typeof SettingsApiKeysRoute;
	"/settings/profile": typeof SettingsProfileRoute;
	"/share/$id": typeof ShareIdRoute;
	"/auth": typeof AuthIndexRoute;
	"/bookmarks": typeof BookmarksIndexRoute;
	"/settings": typeof SettingsIndexRoute;
	"/auth/oauth/$provider/callback": typeof AuthOauthProviderCallbackRoute;
}

export interface FileRoutesByTo {
	"/": typeof IndexLazyRoute;
	"/auth/login": typeof AuthLoginRoute;
	"/auth/register": typeof AuthRegisterRoute;
	"/bookmarks/$id": typeof BookmarksIdRoute;
	"/settings/ai": typeof SettingsAiRoute;
	"/settings/api-keys": typeof SettingsApiKeysRoute;
	"/settings/profile": typeof SettingsProfileRoute;
	"/share/$id": typeof ShareIdRoute;
	"/auth": typeof AuthIndexRoute;
	"/bookmarks": typeof BookmarksIndexRoute;
	"/settings": typeof SettingsIndexRoute;
	"/auth/oauth/$provider/callback": typeof AuthOauthProviderCallbackRoute;
}

export interface FileRoutesById {
	__root__: typeof rootRoute;
	"/": typeof IndexLazyRoute;
	"/auth/login": typeof AuthLoginRoute;
	"/auth/register": typeof AuthRegisterRoute;
	"/bookmarks/$id": typeof BookmarksIdRoute;
	"/settings/ai": typeof SettingsAiRoute;
	"/settings/api-keys": typeof SettingsApiKeysRoute;
	"/settings/profile": typeof SettingsProfileRoute;
	"/share/$id": typeof ShareIdRoute;
	"/auth/": typeof AuthIndexRoute;
	"/bookmarks/": typeof BookmarksIndexRoute;
	"/settings/": typeof SettingsIndexRoute;
	"/auth/oauth/$provider/callback": typeof AuthOauthProviderCallbackRoute;
}

export interface FileRouteTypes {
	fileRoutesByFullPath: FileRoutesByFullPath;
	fullPaths:
		| "/"
		| "/auth/login"
		| "/auth/register"
		| "/bookmarks/$id"
		| "/settings/ai"
		| "/settings/api-keys"
		| "/settings/profile"
		| "/share/$id"
		| "/auth"
		| "/bookmarks"
		| "/settings"
		| "/auth/oauth/$provider/callback";
	fileRoutesByTo: FileRoutesByTo;
	to:
		| "/"
		| "/auth/login"
		| "/auth/register"
		| "/bookmarks/$id"
		| "/settings/ai"
		| "/settings/api-keys"
		| "/settings/profile"
		| "/share/$id"
		| "/auth"
		| "/bookmarks"
		| "/settings"
		| "/auth/oauth/$provider/callback";
	id:
		| "__root__"
		| "/"
		| "/auth/login"
		| "/auth/register"
		| "/bookmarks/$id"
		| "/settings/ai"
		| "/settings/api-keys"
		| "/settings/profile"
		| "/share/$id"
		| "/auth/"
		| "/bookmarks/"
		| "/settings/"
		| "/auth/oauth/$provider/callback";
	fileRoutesById: FileRoutesById;
}

export interface RootRouteChildren {
	IndexLazyRoute: typeof IndexLazyRoute;
	AuthLoginRoute: typeof AuthLoginRoute;
	AuthRegisterRoute: typeof AuthRegisterRoute;
	BookmarksIdRoute: typeof BookmarksIdRoute;
	SettingsAiRoute: typeof SettingsAiRoute;
	SettingsApiKeysRoute: typeof SettingsApiKeysRoute;
	SettingsProfileRoute: typeof SettingsProfileRoute;
	ShareIdRoute: typeof ShareIdRoute;
	AuthIndexRoute: typeof AuthIndexRoute;
	BookmarksIndexRoute: typeof BookmarksIndexRoute;
	SettingsIndexRoute: typeof SettingsIndexRoute;
	AuthOauthProviderCallbackRoute: typeof AuthOauthProviderCallbackRoute;
}

const rootRouteChildren: RootRouteChildren = {
	IndexLazyRoute: IndexLazyRoute,
	AuthLoginRoute: AuthLoginRoute,
	AuthRegisterRoute: AuthRegisterRoute,
	BookmarksIdRoute: BookmarksIdRoute,
	SettingsAiRoute: SettingsAiRoute,
	SettingsApiKeysRoute: SettingsApiKeysRoute,
	SettingsProfileRoute: SettingsProfileRoute,
	ShareIdRoute: ShareIdRoute,
	AuthIndexRoute: AuthIndexRoute,
	BookmarksIndexRoute: BookmarksIndexRoute,
	SettingsIndexRoute: SettingsIndexRoute,
	AuthOauthProviderCallbackRoute: AuthOauthProviderCallbackRoute,
};

export const routeTree = rootRoute
	._addFileChildren(rootRouteChildren)
	._addFileTypes<FileRouteTypes>();

/* ROUTE_MANIFEST_START
{
  "routes": {
    "__root__": {
      "filePath": "__root.tsx",
      "children": [
        "/",
        "/auth/login",
        "/auth/register",
        "/bookmarks/$id",
        "/settings/ai",
        "/settings/api-keys",
        "/settings/profile",
        "/share/$id",
        "/auth/",
        "/bookmarks/",
        "/settings/",
        "/auth/oauth/$provider/callback"
      ]
    },
    "/": {
      "filePath": "index.lazy.tsx"
    },
    "/auth/login": {
      "filePath": "auth/login.tsx"
    },
    "/auth/register": {
      "filePath": "auth/register.tsx"
    },
    "/bookmarks/$id": {
      "filePath": "bookmarks/$id.tsx"
    },
    "/settings/ai": {
      "filePath": "settings/ai.tsx"
    },
    "/settings/api-keys": {
      "filePath": "settings/api-keys.tsx"
    },
    "/settings/profile": {
      "filePath": "settings/profile.tsx"
    },
    "/share/$id": {
      "filePath": "share/$id.tsx"
    },
    "/auth/": {
      "filePath": "auth/index.tsx"
    },
    "/bookmarks/": {
      "filePath": "bookmarks/index.tsx"
    },
    "/settings/": {
      "filePath": "settings/index.tsx"
    },
    "/auth/oauth/$provider/callback": {
      "filePath": "auth/oauth.$provider.callback.tsx"
    }
  }
}
ROUTE_MANIFEST_END */
