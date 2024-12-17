import useSWR from "swr";

// Types
export interface User {
  id: string;
  username?: string;
  email?: string;
}

interface LoginInput {
  email: string;
  password: string;
}

interface RegisterInput {
  username: string;
  email: string;
  password: string;
}

// Utility function for API calls
async function fetchApi<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, {
    credentials: "include", // Important: this enables sending/receiving cookies
    headers: {
      "Content-Type": "application/json",
    },
    ...options,
  });

  if (!response.ok) {
    throw new Error(`HTTP error! status: ${response.status}`);
  }

  const data = await response.json();
  return data.data;
}

// API Functions
const api = {
  login: (input: LoginInput) =>
    fetchApi<User>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  register: (input: RegisterInput) =>
    fetchApi<User>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  logout: () =>
    fetchApi<void>("/api/v1/auth/logout", {
      method: "POST",
    }),

  validateToken: () => fetchApi<User>("/api/v1/auth/validate-jwt"),

  oauthLogin: (provider: string) =>
    (window.location.href = `/api/v1/oauth/${provider}/login`),
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
      api.oauthLogin(provider);
    },
  };
}
