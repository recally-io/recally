import { ROUTES } from "@/lib/router";
import { Navigate, createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/settings/")({
	component: RouteComponent,
});

function RouteComponent() {
	return <Navigate to={ROUTES.SETTINGS_PROFILE} />;
}
