import { BaseLayout } from "@/components/layout/BaseLayout";
import { ThemeProvider } from "@/components/theme-provider";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import AuthPage from "./pages/auth";
import BookmarkPage from "./pages/BookmarkPage";
import HomePage from "./pages/HomePage";

export default function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
      <Router>
        <BaseLayout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/bookmarks/:id" element={<BookmarkPage />} />
            <Route path="/accounts/login" element={<AuthPage />} />
            <Route path="/accounts/signup" element={<AuthPage />} />
          </Routes>
        </BaseLayout>
      </Router>
    </ThemeProvider>
  );
}
