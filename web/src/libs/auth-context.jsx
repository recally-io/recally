import { useMutation, useQuery } from "@tanstack/react-query";
import Cookie from "js-cookie";
import { createContext, useContext, useEffect, useState } from "react";
import { toastError, toastInfo } from "./alert";
import { post } from "./api";
import { checkIsLogin } from "./auth";

export const AuthContext = createContext();

export function useAuthContext() {
  return useContext(AuthContext);
}

export function AuthContextProvider({ children }) {
  const url = new URL(window.location.href);
  const ref = url.searchParams.get("ref");
  const [isLogin, setIsLogin] = useState(false);
  const checkIsLoginQuery = useQuery({
    queryKey: ["check-login"],
    queryFn: checkIsLogin,
    enabled: !isLogin,
  });

  useEffect(() => {
    if (!checkIsLoginQuery.isLoading && checkIsLoginQuery.data) {
      setIsLogin(checkIsLoginQuery.data);
    }
  }, [checkIsLoginQuery.data]);

  const loginMutation = useMutation({
    mutationFn: async (credentials) => {
      const response = await post("/api/v1/auth/login", null, credentials);
      return response.data;
    },
    onSuccess: (data) => {
      setIsLogin(true);
      toastInfo("Logged in successfully");
      setTimeout(() => {
        window.location.href = ref || "/";
      }, 1000);
    },
    onError: (error) => {
      toastError(`Login failed: ${error.message}`);
    },
  });

  const registerMutation = useMutation({
    mutationFn: async (userData) => {
      const response = await post("/api/v1/auth/register", null, userData);
      return response.data;
    },
    onSuccess: (data) => {
      setIsLogin(true);
      toastInfo("Registered successfully");
      setTimeout(() => {
        window.location.href = ref || "/";
      }, 1000);
    },
    onError: (error) => {
      toastError(`Registration failed: ${error.message}`);
    },
  });

  const logout = () => {
    Cookie.remove("token");
    setIsLogin(false);
    toastInfo("Logged out successfully");
    setTimeout(() => {
      window.location.href = "/";
    }, 1000);
  };

  return (
    <AuthContext.Provider
      value={{
        isLogin,
        checkIsLogin: checkIsLoginQuery,
        login: loginMutation,
        register: registerMutation,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
