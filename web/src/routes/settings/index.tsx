import { ROUTES } from "@/lib/router";
import { Navigate, createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/settings/")({
	component: RouteComponent,
});

function RouteComponent() {
	return <Navigate to={ROUTES.SETTINGS_PROFILE} />;
}
