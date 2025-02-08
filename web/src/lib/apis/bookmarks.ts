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

export interface BookmarkMetadata {
	reading_progress?: number;
	last_read_at?: string;

	highlights?: Highlight[];
}

export interface BookmarkContentFileMetadata {
	name?: string;
	extension?: string;
	mime_type?: string;
	size?: number;
	page_count?: number;
}

export interface BookmarkContentMetadata {
	author?: string;
	published_at?: string;
	description?: string;
	site_name?: string;
	domain?: string;

	favicon?: string;
	cover?: string;

	file?: BookmarkContentFileMetadata;
}

export interface BookmarkContent {
	type: string;
	url?: string;
	domain?: string;
	title?: string;
	description?: string;
	summary?: string;
	content?: string;
	html?: string;
	s3_key?: string;
	metadata?: BookmarkContentMetadata;
}

export interface Bookmark {
	id: string;
	userId: string;
	content_id: string;
	is_favorite: boolean;
	is_archive: boolean;
	is_public: boolean;
	created_at: string;
	updated_at: string;
	content: BookmarkContent;
	metadata?: BookmarkMetadata;
	tags?: string[];
	share?: ShareContent;
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
	title?: string;
	description?: string;
	tags?: string[];
	content?: string;
	type?: string;
	s3_key?: string;
	metadata?: BookmarkContentMetadata;
}

interface BookmarkUpdateInput {
	summary?: string;
	content?: string;
	html?: string;
	metadata?: BookmarkContentMetadata;
}

interface ShareContentUpdateInput {
	expires_at: string;
}

interface BookmarkRefreshInput {
	fetcher?: string;
	regenerate_summary?: boolean;
	is_proxy_image?: boolean;
}

export interface ShareBookmarkRequest {
	expires_at: string;
}

export interface ShareContent {
	id: string;
	bookmark_id: string;
	expires_at?: string;
	created_at: string;
}

// API Functions
const api = {
	list: (filters: string[] = [], query = "", limit = 20, offset = 0) => {
		const params = new URLSearchParams();
		params.set("limit", limit.toString());
		params.set("offset", offset.toString());
		params.set("query", query);
		// Append each filter separately
		for (const filter of filters) {
			params.append("filter", filter);
		}

		const url = `/api/v1/bookmarks?${params.toString()}`;
		return fetcher<ListBookmarksResponse>(url);
	},

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

	shareContent: (id: string, request: ShareBookmarkRequest) =>
		fetcher<ShareContent>(`/api/v1/bookmarks/${id}/share`, {
			method: "POST",
			body: JSON.stringify(request),
		}),

	getSharedContent: (token: string) =>
		fetcher<Bookmark>(`/api/v1/shared/${token}`),

	updateSharedContent: (id: string, input: ShareContentUpdateInput) =>
		fetcher<Bookmark>(`/api/v1/bookmarks/${id}/share`, {
			method: "PUT",
			body: JSON.stringify(input),
		}),

	unshareContent: (id: string) =>
		fetcher<void>(`/api/v1/bookmarks/${id}/share`, {
			method: "DELETE",
		}),
};

// SWR Hooks
export function useBookmarks(
	limit = 20,
	offset = 0,
	filters: string[] = [],
	query = "",
) {
	return useSWR<ListBookmarksResponse>(
		["bookmarks", filters, query, limit, offset],
		() => api.list(filters, query, limit, offset),
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

export function useSharedContent(token: string) {
	return useSWR<Bookmark>("sharedContent", () => api.getSharedContent(token));
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

export function useShareContentMutations() {
	const { mutate } = useSWRConfig();

	return {
		shareContent: async (id: string, request: ShareBookmarkRequest) => {
			const response = await api.shareContent(id, request);
			mutate(["bookmark", id]);
			return response;
		},

		unshareContent: async (id: string) => {
			await api.unshareContent(id);
			mutate(["bookmark", id]);
		},

		updateSharedContent: async (id: string, input: ShareContentUpdateInput) => {
			const response = await api.updateSharedContent(id, input);
			mutate(["bookmark", id]);
			return response;
		},
	};
}
