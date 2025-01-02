import AuthComponent from "@/components/auth/auth";
import { createLazyFileRoute } from "@tanstack/react-router";

export const Route = createLazyFileRoute("/auth/login")({
	component: RouteComponent,
});

function RouteComponent() {
	return <AuthComponent mode="login" />;
}
