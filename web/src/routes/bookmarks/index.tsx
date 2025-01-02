import BookmarksListView from "@/components/bookmarks/bookmarks-list-page";
import ProtectedRoute from "@/components/protected-route";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/bookmarks/")({
	component: RouteComponent,
});

function RouteComponent() {
	return (
		<ProtectedRoute>
			<BookmarksListView />
		</ProtectedRoute>
	);
}
