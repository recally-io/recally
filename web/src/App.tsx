import { BaseLayout } from "@/components/layout/BaseLayout";
import { ThemeProvider } from "@/components/theme-provider";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import AuthPage from "./pages/auth";
import BookmarkPage from "./pages/BookmarkPage";
import HomePage from "./pages/HomePage";
import ProtectedRoute from "@/components/ProtectedRoute";
import { SWRConfig } from "swr";

export default function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <SWRConfig
        value={{
          // Define your global configuration options here
          fetcher: (resource, init) =>
            fetch(resource, init).then((res) => res.json()),
          // ...other global configurations...
        }}
      >
        <Router>
          <BaseLayout>
            <Routes>
              {/* Public routes */}
              <Route path="/accounts/login" element={<AuthPage />} />
              <Route path="/accounts/signup" element={<AuthPage />} />
              {/* Protected routes */}
              <Route
                path="/"
                element={
                  <ProtectedRoute>
                    <HomePage />
                  </ProtectedRoute>
                }
              />
              <Route
                path="/bookmarks/:id"
                element={
                  <ProtectedRoute>
                    <BookmarkPage />
                  </ProtectedRoute>
                }
              />
            </Routes>
          </BaseLayout>
        </Router>
      </SWRConfig>
    </ThemeProvider>
  );
}
