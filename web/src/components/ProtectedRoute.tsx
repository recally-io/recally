import React from "react";
import { useUser } from "@/lib/apis/auth";
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
		return <Navigate to="/accounts/login" replace />;
	}

	return children;
}
