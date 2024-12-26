import ProtectedRoute from "@/components/ProtectedRoute";
import { ThemeProvider } from "@/components/theme-provider";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import { SWRConfig } from "swr";
import fetcher from "./lib/apis/fetcher";
import { ROUTES } from "./lib/router";
import BookmarkDetailPage from "./pages/Bookmark-detail-page";
import AuthPage from "./pages/auth";
import BookmarkPage from "./pages/bookmarks-page";

export default function App() {
	return (
		<ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
			<SWRConfig
				value={{
					// Define your global configuration options here
					fetcher: fetcher,
					// ...other global configurations...
				}}
			>
				<Router>
					<Routes>
						{/* Public routes */}
						<Route path={ROUTES.LOGIN} element={<AuthPage />} />
						<Route path={ROUTES.SIGNUP} element={<AuthPage />} />
						{/* Protected routes */}
						<Route
							path="/"
							element={
								<ProtectedRoute>
									<BookmarkPage />
								</ProtectedRoute>
							}
						/>
						<Route
							path={ROUTES.BOOKMARKS}
							element={
								<ProtectedRoute>
									<BookmarkPage />
								</ProtectedRoute>
							}
						/>
						<Route
							path={`${ROUTES.BOOKMARKS}/:id`}
							element={
								<ProtectedRoute>
									<BookmarkDetailPage />
								</ProtectedRoute>
							}
						/>
					</Routes>
				</Router>
			</SWRConfig>
		</ThemeProvider>
	);
}
