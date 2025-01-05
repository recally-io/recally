import useSWR, { useSWRConfig } from "swr";
import fetcher from "./fetcher";

// Types
export interface Highlight {
	id: string;
	text: string;
	startOffset: number;
	endOffset: number;
	note?: string;
}

export interface Metadata {
	tags?: string[];
	highlights?: Highlight[];

	author?: string;
	published_at?: string;
	description?: string;
	site_name?: string;
	domain?: string;

	favicon?: string;
	cover?: string;
	image?: string;
}

export interface Bookmark {
	id: string;
	userId: string;
	type: string;
	url?: string;
	domain?: string;
	title?: string;
	description?: string;
	summary?: string;
	content?: string;
	tags?: string[];
	html?: string;
	metadata?: Metadata;
	created_at: string;
	updated_at: string;
}

export interface ListBookmarksResponse {
	bookmarks: Bookmark[];
	total: number;
	limit: number;
	offset: number;
}

export interface Tag {
	name: string;
	count: number;
}

export interface Domain {
	name: string;
	count: number;
}

interface BookmarkCreateInput {
	url: string;
	metadata?: Metadata;
}

interface BookmarkUpdateInput {
	summary?: string;
	content?: string;
	html?: string;
	metadata?: Metadata;
}

interface BookmarkRefreshInput {
	fetcher?: "http" | "jina" | "browser";
	regenerate_summary?: boolean;
}

// API Functions
const api = {
	list: (filter = "", query = "", limit = 20, offset = 0) =>
		// filter=category:article;type:rss
		fetcher<ListBookmarksResponse>(
			`/api/v1/bookmarks?limit=${limit}&offset=${offset}&query=${query}&filter=${filter}`,
		),

	create: (input: BookmarkCreateInput) =>
		fetcher<Bookmark>("/api/v1/bookmarks", {
			method: "POST",
			body: JSON.stringify(input),
		}),

	get: (id: string) => fetcher<Bookmark>(`/api/v1/bookmarks/${id}`),

	update: (id: string, input: BookmarkUpdateInput) =>
		fetcher<Bookmark>(`/api/v1/bookmarks/${id}`, {
			method: "PUT",
			body: JSON.stringify(input),
		}),

	delete: (id: string) =>
		fetcher<void>(`/api/v1/bookmarks/${id}`, {
			method: "DELETE",
		}),

	deleteAll: (userId: string) =>
		fetcher<void>(`/api/v1/bookmarks?user-id=${userId}`, {
			method: "DELETE",
		}),

	refresh: (id: string, input: BookmarkRefreshInput) =>
		fetcher<Bookmark>(`/api/v1/bookmarks/${id}/refresh`, {
			method: "POST",
			body: JSON.stringify(input),
		}),

	listTags: () => fetcher<Tag[]>("/api/v1/bookmarks/tags"),

	listDomains: () => fetcher<Domain[]>("/api/v1/bookmarks/domains"),
};

// SWR Hooks
export function useBookmarks(limit = 20, offset = 0, filter = "", query = "") {
	return useSWR<ListBookmarksResponse>(
		["bookmarks", filter, query, limit, offset],
		() => api.list(filter, query, limit, offset),
	);
}

export function useBookmark(id: string) {
	return useSWR<Bookmark>(id ? ["bookmark", id] : null, () => api.get(id));
}

export function useTags() {
	return useSWR<Tag[]>("bookmarkTags", () => api.listTags());
}

export function useDomains() {
	return useSWR<Domain[]>("bookmarkDomains", () => api.listDomains());
}

// Mutation Hooks
export function useBookmarkMutations() {
	const { mutate } = useSWRConfig();

	const invalidateBookmarks = () => {
		mutate((key: unknown) => Array.isArray(key) && key[0] === "bookmarks");
	};

	return {
		createBookmark: async (input: BookmarkCreateInput) => {
			const bookmark = await api.create(input);
			invalidateBookmarks();
			return bookmark;
		},

		updateBookmark: async (id: string, input: BookmarkUpdateInput) => {
			const bookmark = await api.update(id, input);
			mutate(["bookmark", id]);
			invalidateBookmarks();
			return bookmark;
		},

		deleteBookmark: async (id: string) => {
			await api.delete(id);
			mutate(["bookmark", id], null);
			invalidateBookmarks();
		},

		deleteAllBookmarks: async (userId: string) => {
			await api.deleteAll(userId);
			invalidateBookmarks();
		},

		refreshBookmark: async (id: string, input: BookmarkRefreshInput) => {
			const bookmark = await api.refresh(id, input);
			mutate(["bookmark", id]);
			invalidateBookmarks();
			return bookmark;
		},
	};
}
