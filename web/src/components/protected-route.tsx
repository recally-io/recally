import { useUser } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import type React from "react";

interface ProtectedRouteProps {
	children: React.ReactElement;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
	const { user, isLoading } = useUser();

	if (isLoading) {
		// Render a loading indicator or null
		return null; // or your loading component
	}

	if (!user) {
		window.location.href = ROUTES.LOGIN;
		return null;
	}

	return children;
}
