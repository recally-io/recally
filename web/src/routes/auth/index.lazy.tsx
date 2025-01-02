import { ROUTES } from "@/lib/router";
import { Navigate, createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/auth/")({
	component: RouteComponent,
});

function RouteComponent() {
	return <Navigate to={ROUTES.AUTH_LOGIN} />;
}
