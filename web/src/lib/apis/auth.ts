import useSWR from "swr";
import fetcher from "./fetcher";
import type { User } from "./users";

interface LoginInput {
	email: string;
	password: string;
}

interface RegisterInput {
	username: string;
	email: string;
	password: string;
}

interface OAuthLoginResponse {
	url: string;
}

// API Functions
const api = {
	login: (input: LoginInput) =>
		fetcher<User>("/api/v1/auth/login", {
			method: "POST",
			body: JSON.stringify(input),
		}),

	register: (input: RegisterInput) =>
		fetcher<User>("/api/v1/auth/register", {
			method: "POST",
			body: JSON.stringify(input),
		}),

	logout: () =>
		fetcher<void>("/api/v1/auth/logout", {
			method: "POST",
		}),

	validateToken: () => fetcher<User>("/api/v1/auth/validate-jwt"),

	oauthLogin: (provider: string) => {
		return fetcher<OAuthLoginResponse>(
			`/api/v1/oauth/${provider.toLowerCase()}/login`,
		);
	},

	OAuthCallback: (provider: string, code: string) => {
		return fetcher<User>(`/api/v1/oauth/${provider}/callback?code=${code}`);
	},
};

// SWR Hooks
export function useUser() {
	const { data, error, mutate } = useSWR<User>("auth-user", api.validateToken, {
		// Adjust SWR options for caching
		revalidateOnFocus: false,
		revalidateIfStale: false,
		revalidateOnReconnect: false,
		dedupingInterval: 60000, // 1 minute
		shouldRetryOnError: false,
	});

	return {
		user: data,
		isLoading: !error && !data,
		isError: error,
		mutate,
	};
}

// Mutation Hooks
export function useAuth() {
	const { mutate } = useUser();

	return {
		login: async (input: LoginInput) => {
			const user = await api.login(input);
			await mutate(user);
			return user;
		},

		register: async (input: RegisterInput) => {
			const user = await api.register(input);
			await mutate(user);
			return user;
		},

		logout: async () => {
			await api.logout();
			await mutate(undefined);
		},

		oauthLogin: (provider: string) => {
			return api.oauthLogin(provider);
		},

		oauthCallback: async (provider: string, code: string) => {
			const data = await api.OAuthCallback(provider, code);
			return data;
		},
	};
}
