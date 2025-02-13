import { useAuth } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useEffect } from "react";

type CallbackSchema = {
	state: string;
	code: string;
};

export const Route = createFileRoute("/auth/oauth/$provider/callback")({
	component: RouteComponent,
	// validateSearch: (search: URLSearchParams): CallbackSchema => {
	//   const state = search.get('state')
	//   const code = search.get('code')
	//   if (!state || !code) {
	//     throw new Error('Invalid OAuth callback')
	//   }
	//   return { state, code }
	// },
});

function RouteComponent() {
	const router = useRouter();
	const { state, code }: CallbackSchema = Route.useSearch();
	const { provider } = Route.useParams();

	const { oauthCallback } = useAuth();

	useEffect(() => {
		const handleCallback = async () => {
			if (state && code) {
				await handleOAuthCallback();
			}
		};
		handleCallback();
	}, [state, code]);

	const handleOAuthCallback = async () => {
		if (state && code) {
			await oauthCallback(provider.toLowerCase(), code, state);
			router.navigate({ to: ROUTES.HOME });
		}
	};

	return <></>;
}
