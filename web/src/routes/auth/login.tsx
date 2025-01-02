import AuthComponent from "@/components/auth/auth";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/auth/login")({
	component: RouteComponent,
});

function RouteComponent() {
	return <AuthComponent mode="login" />;
}
