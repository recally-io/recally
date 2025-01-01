import { useSWRConfig } from "swr";
import fetcher from "./fetcher";

export type SummaryConfig = {
	model?: string;
	prompt?: string;
	language?: string;
};

export type UserSettings = {
	summary_options?: SummaryConfig;
};

export interface User {
	id: string;
	avatar?: string;
	username?: string;
	email?: string;
	phone?: string;
	Status?: string;
	Settings?: UserSettings;
}

// API Functions
const api = {
	updateInfo: (
		userId: string,
		data: { username?: string; email?: string; phone?: string },
	) =>
		fetcher<User>(`/api/v1/users/${userId}/info`, {
			method: "PUT",
			body: JSON.stringify(data),
		}),

	updateSettings: (userId: string, settings: User["Settings"]) =>
		fetcher<User>(`/api/v1/users/${userId}/settings`, {
			method: "PUT",
			body: JSON.stringify({ settings }),
		}),

	updatePassword: (userId: string, currentPassword: string, password: string) =>
		fetcher<User>(`/api/v1/users/${userId}/password`, {
			method: "PUT",
			body: JSON.stringify({ currentPassword, password }),
		}),
};

// Mutation Hooks
export function useUsers() {
	const { mutate } = useSWRConfig();

	const invalidateUserData = (userId: string) => {
		// Invalidate auth-user cache
		mutate("auth-user");
		// Invalidate specific user cache if needed
		mutate(["user", userId]);
	};

	return {
		updateInfo: async (
			userId: string,
			data: { username?: string; email?: string; phone?: string },
		) => {
			const user = await api.updateInfo(userId, data);
			invalidateUserData(userId);
			return user;
		},

		updateSettings: async (userId: string, settings: User["Settings"]) => {
			const user = await api.updateSettings(userId, settings);
			invalidateUserData(userId);
			return user;
		},

		updatePassword: async (
			userId: string,
			currentPassword: string,
			password: string,
		) => {
			const user = await api.updatePassword(userId, currentPassword, password);
			invalidateUserData(userId);
			return user;
		},
	};
}
