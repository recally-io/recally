import ProtectedRoute from "@/components/protected-route";

import { ProfileSettings } from "@/components/settings/profile";
import { SummarySettings } from "@/components/settings/summary";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ROUTES } from "@/lib/router";
import { Bot, Home, Menu, User } from "lucide-react";
import { parseAsString, useQueryState } from "nuqs";
import { createRoot } from "react-dom/client";
import App from "./app-basic";

const settingsEnum = {
	GENERAL: "general",
	PROFILE: "profile",
	SUMMARY: "summary",
};

function SidebarNav() {
	return (
		<nav className="space-y-2 w-48">
			<a
				href={ROUTES.HOME}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<Home className="h-4 w-4" />
				Home
			</a>
			<a
				href={`${ROUTES.SETTINGS}?tab=${settingsEnum.PROFILE}`}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<User className="h-4 w-4" />
				Profile
			</a>
			<a
				href={`${ROUTES.SETTINGS}?tab=${settingsEnum.SUMMARY}`}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<Bot className="h-4 w-4" />
				Summary
			</a>
		</nav>
	);
}

function SettingsPage() {
	const [tab, _] = useQueryState(
		"tab",
		parseAsString.withDefault(settingsEnum.PROFILE),
	);

	const mainTab = () => {
		if (tab === settingsEnum.PROFILE) {
			return <ProfileSettings />;
		} else if (tab === settingsEnum.SUMMARY) {
			return <SummarySettings />;
		}
	};

	return (
		<div className="container mx-auto py-10 px-4">
			<h1 className="text-2xl font-semibold mb-8">Preferences</h1>

			<div className="flex gap-4">
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
				<div className="w-full">{mainTab()}</div>
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
