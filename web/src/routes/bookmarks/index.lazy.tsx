import BookmarksListView from "@/components/bookmarks/bookmarks-list";
import ProtectedRoute from "@/components/protected-route";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/bookmarks/")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<ProtectedRoute>
			<BookmarksListView />
		</ProtectedRoute>
	);
}
