import { useAuth } from "@/lib/apis/auth";
import { ROUTES } from "@/lib/router";
import { createLazyFileRoute, useRouter } from "@tanstack/react-router";
import { useEffect } from "react";

type CallbackSchema = {
	state: string;
	code: string;
};

export const Route = createLazyFileRoute("/auth/oauth/$provider/callback")({
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
				console.log("OAuth Callback:", state, code);
				await handleOAuthCallback();
			}
		};
		handleCallback();
	}, [state, code]);

	const handleOAuthCallback = async () => {
		if (state && code) {
			const data = await oauthCallback(provider.toLowerCase(), code, state);
			console.log("OAuth Callback data:", data);
			router.navigate({ to: ROUTES.HOME });
		}
	};

	return <></>;
}
