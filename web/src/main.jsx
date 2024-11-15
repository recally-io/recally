import "@mantine/core/styles.css";

import { LoadingOverlay } from "@mantine/core";
import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import ReactDOM from "react-dom/client";
import {
  createBrowserRouter,
  Navigate,
  RouterProvider,
} from "react-router-dom";
import Assistants from "./components/assistants";
import { AuthenticationForm } from "./components/auth";
import { Layout, ThemeProvider } from "./components/layout";
import { queryClient } from "./libs/api";
import { AuthContextProvider, useAuthContext } from "./libs/auth-context";
import Home from "./pages/home";
import Threads from "./pages/threads";
import Bookmarks from "./pages/bookmarks";

const ProtectedRoute = ({ children }) => {
  const { checkIsLogin } = useAuthContext();

  if (checkIsLogin.isLoading) {
    return <LoadingOverlay visible={true} />;
  }

  if (!checkIsLogin.isLoading && !checkIsLogin.data) {
    return <Navigate to="/auth" replace />;
  }

  return children;
};

const router = createBrowserRouter([
  {
    path: "/",
    element: <Layout main={<Home />} />,
  },
  {
    path: "auth",
    element: <Layout main={<AuthenticationForm />} />,
  },
  {
    path: "/assistants",
    element: (
      <ProtectedRoute>
        <Layout main={<Assistants />} />
      </ProtectedRoute>
    ),
  },
  {
    path: "/assistants/:assistantId",
    element: (
      <ProtectedRoute>
        <Layout main={<Assistants />} />
      </ProtectedRoute>
    ),
  },
  {
    path: "/assistants/:assistantId/threads",
    element: (
      <ProtectedRoute>
        <Threads />
      </ProtectedRoute>
    ),
  },
  {
    path: "/assistants/:assistantId/threads/:threadId",
    element: (
      <ProtectedRoute>
        <Threads />
      </ProtectedRoute>
    ),
  },
  {
    path: "/bookmarks",
    element: (
      <ProtectedRoute>
        <Layout main={<Bookmarks />} />
      </ProtectedRoute>
    ),
  },
]);

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <ThemeProvider>
      <QueryClientProvider client={queryClient}>
        <AuthContextProvider>
          <RouterProvider router={router} />
        </AuthContextProvider>
      </QueryClientProvider>
    </ThemeProvider>
  </React.StrictMode>,
);
