import BookmarksListView from "@/components/bookmarks/bookmarks-list-page";
import type { BookmarkSearch } from "@/components/bookmarks/types";
import ProtectedRoute from "@/components/protected-route";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/bookmarks/")({
	component: RouteComponent,
	validateSearch: (search: Record<string, unknown>): BookmarkSearch => {
		// validate and parse the search params into a typed state
		return {
			page: Number(search?.page ?? 1),
			filters: (search.filters as string[]) || [],
			query: (search.query as string) || "",
		};
	},
});

function RouteComponent() {
	const search = Route.useSearch();
	return (
		<ProtectedRoute>
			<BookmarksListView search={search} />
		</ProtectedRoute>
	);
}
