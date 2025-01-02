import BookmarkDetailPage from "@/components/bookmarks/bookmark-detail";
import ProtectedRoute from "@/components/protected-route";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/bookmarks/$id")({
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
