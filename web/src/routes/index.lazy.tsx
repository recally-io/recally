import { ROUTES } from "@/lib/router";
import { Navigate, createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/")({
	component: Index,
});

function Index() {
	return <Navigate to={ROUTES.BOOKMARKS} />;
}
