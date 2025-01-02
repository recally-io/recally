import BookmarkDetailPage from "@/components/bookmarks/bookmark-detail";
import ProtectedRoute from "@/components/protected-route";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/bookmarks/$id")({
	component: RouteComponent,
});

function RouteComponent() {
	const { id } = Route.useParams();
	return (
		<ProtectedRoute>
			<BookmarkDetailPage id={id} />
		</ProtectedRoute>
	);
}
