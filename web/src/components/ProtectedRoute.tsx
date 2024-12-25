import { useUser } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import type React from "react";
import { Navigate } from "react-router-dom";

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
		return <Navigate to={ROUTES.LOGIN} replace />;
	}

	return children;
}
