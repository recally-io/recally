import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Bot, Home, Menu, User } from "lucide-react";

function SidebarNav() {
	return (
		<nav className="space-y-2 w-48">
			<Link
				to={ROUTES.HOME}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<Home className="h-4 w-4" />
				Home
			</Link>
			<Link
				to={ROUTES.SETTINGS_PROFILE}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<User className="h-4 w-4" />
				Profile
			</Link>
			<Link
				to={ROUTES.SETTINGS_SUMMARY}
				className="px-3 py-1 text-sm text-muted-foreground hover:text-foreground flex items-center gap-2"
			>
				<Bot className="h-4 w-4" />
				Summary
			</Link>
		</nav>
	);
}

export function SettingsPageComponenrt({
	children,
}: { children: React.ReactElement }) {
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
				<div className="w-full">{children}</div>
			</div>
		</div>
	);
}
