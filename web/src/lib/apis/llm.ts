import useSWR from "swr";
import fetcher from "./fetcher";

// Types
export interface Model {
	id: string;
	name: string;
}

export interface Tool {
	name: string;
	description: string;
	parameters: object;
}

// API Functions
const api = {
	listModels: () => fetcher<Model[]>("/api/v1/llm/models"),

	listTools: () => fetcher<Tool[]>("/api/v1/llm/tools"),
};

// SWR Hooks
export function useLLMs() {
	const models = useSWR<Model[]>("llm-models", api.listModels, {
		revalidateOnFocus: false,
		revalidateIfStale: false,
		revalidateOnReconnect: false,
		dedupingInterval: 3600000, // 1 hour
	});

	const tools = useSWR<Tool[]>("llm-tools", api.listTools, {
		revalidateOnFocus: false,
		revalidateIfStale: false,
		revalidateOnReconnect: false,
		dedupingInterval: 3600000, // 1 hour
	});

	return {
		models: models.data,
		tools: tools.data,
		isLoading: !models.data || !tools.data,
		isError: models.error || tools.error,
	};
}
