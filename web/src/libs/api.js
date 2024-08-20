import { AssistantsApi, AuthApi, Configuration, ToolsApi } from "../sdk/index";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

export const queryClient = new QueryClient();

const config = new Configuration({
  // apiKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQxNDAwNTQsInVzZXJfaWQiOiJmZDM5YjM5Zi1hOTc1LTQyYzYtODdjMy1mNDFiMjdlZGY1NTcifQ.6o7N3RoBltVsfHaYo1rju4oFnpZ1J64zGEhbbSckjgc",
  // accessToken:
  //     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQxNDAwNTQsInVzZXJfaWQiOiJmZDM5YjM5Zi1hOTc1LTQyYzYtODdjMy1mNDFiMjdlZGY1NTcifQ.6o7N3RoBltVsfHaYo1rju4oFnpZ1J64zGEhbbSckjgc",
  credentials: "include",
});

export const assistantAPI = new AssistantsApi(config);
export const authApi = new AuthApi(config);
export const toolsApi = new ToolsApi(config);
