import useSWR from "swr";
import fetcher from "./fetcher";
import { authClient } from "@/lib/auth-client";
import type { LoginInput, RegisterInput, User } from "@/lib/auth-client";

export interface ApiKey {
	id: string;
	name: string;
	prefix: string;
	hash: string;
	scopes: string[];
	expires_at: number;
	created_at: number;
}

export interface CreateApiKeyInput {
	name: string;
	prefix?: string;
	scopes?: string[];
	expires_at: Date;
}

// API Functions
const api = {
	// API Key operations (unchanged)
	createApiKey: (input: CreateApiKeyInput) =>
		fetcher<ApiKey>("/api/v1/auth/keys", {
			method: "POST",
			body: JSON.stringify(input),
		}),

	listApiKeys: (prefix?: string, isActive?: boolean) => {
		const params = new URLSearchParams();
		if (prefix) params.append("prefix", prefix);
		if (isActive !== undefined) params.append("is_active", String(isActive));
		return fetcher<ApiKey[]>(`/api/v1/auth/keys?${params.toString()}`);
	},

	deleteApiKey: (id: string) =>
		fetcher<void>(`/api/v1/auth/keys/${id}`, {
			method: "DELETE",
		}),

	// OAuth callback (still needed for callback route)
	OAuthCallback: (provider: string, code: string) => {
		return fetcher<User>(`/api/v1/oauth/${provider}/callback?code=${code}`);
	},
};

// SWR Hooks
export function useUser() {
	const { data, error, mutate } = useSWR<User>(
		"auth-user",
		() => authClient.validateSession(),
		{
			// Adjust SWR options for caching
			revalidateOnFocus: false,
			revalidateIfStale: false,
			revalidateOnReconnect: false,
			dedupingInterval: 60000, // 1 minute
			shouldRetryOnError: false,
		},
	);

	return {
		user: data,
		isLoading: !error && !data,
		isError: error,
		mutate,
	};
}

export function useApiKeys(prefix?: string, isActive?: boolean) {
	const { data, error, mutate } = useSWR<ApiKey[]>(
		["api-keys", prefix, isActive],
		() => api.listApiKeys(prefix, isActive),
	);

	return {
		keys: data,
		isLoading: !error && !data,
		isError: error,
		mutate,
	};
}

export function useApiKeysMutations() {
	const { mutate: mutateAll } = useApiKeys();

	return {
		createApiKey: async (input: CreateApiKeyInput) => {
			const key = await api.createApiKey(input);
			// Mutate all related queries
			await mutateAll();
			return key;
		},

		deleteApiKey: async (id: string) => {
			await api.deleteApiKey(id);
			// Mutate all related queries
			await mutateAll();
		},
	};
}

// Mutation Hooks
export function useAuth() {
	const { mutate } = useUser();

	return {
		login: async (input: LoginInput) => {
			const user = await authClient.login(input.email, input.password);
			await mutate(user);
			return user;
		},

		register: async (input: RegisterInput) => {
			const user = await authClient.register(
				input.username,
				input.email,
				input.password,
			);
			await mutate(user);
			return user;
		},

		logout: async () => {
			await authClient.logout();
			await mutate(undefined);
		},

		oauthLogin: async (provider: string) => {
			return authClient.getOAuthURL(provider);
		},

		oauthCallback: async (provider: string, code: string, state: string) => {
			console.log("OAuth Callback", provider, code, state);
			const data = await api.OAuthCallback(provider, code);
			return data;
		},
	};
}
