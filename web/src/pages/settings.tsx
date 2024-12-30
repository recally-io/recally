import ProtectedRoute from "@/components/protected-route";

import { ProfileSettings } from "@/components/settings/profile";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Home, Menu, User } from "lucide-react";
import { parseAsString, useQueryState } from "nuqs";
import { createRoot } from "react-dom/client";
import App from "./app-basic";

const settingsEnum = {
	GENERAL: "general",
	PROFILE: "profile",
};

function SidebarNav() {
	return (
		<nav className="space-y-2 w-48">
			<a
				href="/"
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<Home className="h-4 w-4" />
				Home
			</a>
			<a
				href={`/settings?tab=${settingsEnum.PROFILE}`}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<User className="h-4 w-4" />
				Profile
			</a>
		</nav>
	);
}

function SettingsPage() {
	const [tab, _] = useQueryState(
		"tab",
		parseAsString.withDefault(settingsEnum.GENERAL),
	);

	const mainTab = () => {
		if (tab === settingsEnum.PROFILE) {
			return <ProfileSettings />;
		}
	};

	return (
		<div className="container mx-auto py-10">
			<h1 className="text-2xl font-semibold mb-8">Preferences</h1>

			<div className="flex gap-12">
				<div className="hidden md:block">
					<SidebarNav />
				</div>
				<Sheet>
					<SheetTrigger asChild>
						<Button variant="outline" size="icon" className="md:hidden">
							<Menu className="h-4 w-4" />
						</Button>
					</SheetTrigger>
					<SheetContent side="left" className="w-64">
						<SidebarNav />
					</SheetContent>
				</Sheet>
				{mainTab()}
			</div>
		</div>
	);
}

createRoot(document.getElementById("root")!).render(
	<App>
		<ProtectedRoute>
			<SettingsPage />
		</ProtectedRoute>
	</App>,
);
