import ProtectedRoute from "@/components/ProtectedRoute";
import { BaseLayout } from "@/components/layout/BaseLayout";
import { ThemeProvider } from "@/components/theme-provider";
import { Route, BrowserRouter as Router, Routes } from "react-router-dom";
import { SWRConfig } from "swr";
import fetcher from "./lib/apis/fetcher";
import { ROUTES } from "./lib/router";
import BookmarkPage from "./pages/BookmarkPage";
import HomePage from "./pages/HomePage";
import AuthPage from "./pages/auth";

export default function App() {
	ROUTES.BOOKMARKS
	return (
		<ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
			<SWRConfig
				value={{
					// Define your global configuration options here
					fetcher: fetcher,
					// ...other global configurations...
				}}
			>
				<Router>
					<BaseLayout>
						<Routes>
							{/* Public routes */}
							<Route path={ROUTES.LOGIN} element={<AuthPage />} />
							<Route path={ROUTES.SIGNUP} element={<AuthPage />} />
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
								path={`${ROUTES.BOOKMARKS}/:id`}
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
