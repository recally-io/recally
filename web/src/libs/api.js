import { QueryClient } from "@tanstack/react-query";
import { AssistantsApi, AuthApi, Configuration, ToolsApi } from "../sdk/index";

export const queryClient = new QueryClient();

const config = new Configuration({
    credentials: "include",
});

export const assistantAPI = new AssistantsApi(config);
export const authApi = new AuthApi(config);
export const toolsApi = new ToolsApi(config);
