export type View = "grid" | "list";

export type BookmarkSearch = {
	page: number;
	// filter: site:github.com,category:url,tag:tag1
	filters: string[];
	// query: search query
	query: string;
};
