import AuthComponent from "@/components/auth/auth";
import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/auth/register")({
	component: RouteComponent,
});

function RouteComponent() {
	return <AuthComponent mode="register" />;
}
