import { SidebarComponent } from "@/components/sidebar/sidebar";
import {
	SidebarGroup,
	SidebarInset,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Bot, Key, User } from "lucide-react";
import ProtectedRoute from "../protected-route";

const items = [
	{
		title: "Profile",
		url: ROUTES.SETTINGS_PROFILE,
		icon: User,
	},
	{
		title: "API Keys",
		url: ROUTES.SETTINGS_API_KEYS,
		icon: Key,
	},
	{
		title: "AI",
		url: ROUTES.SETTINGS_AI,
		icon: Bot,
	},
];

function sidebarNavContent() {
	return (
		<SidebarGroup>
			<SidebarMenu>
				{items.map((item) => (
					<SidebarMenuItem key={item.title}>
						<SidebarMenuButton tooltip={item.title}>
							<Link to={item.url} className="flex items-center gap-2">
								{item.icon && <item.icon />}
								<span>{item.title}</span>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
				))}
			</SidebarMenu>
		</SidebarGroup>
	);
}

export function SettingsPageComponenrt({
	children,
}: { children: React.ReactElement }) {
	return (
		<ProtectedRoute>
			<SidebarProvider defaultOpen={false}>
				<SidebarComponent>{sidebarNavContent()}</SidebarComponent>
				<SidebarInset>
					<div className="flex flex-col container mx-auto h-full">
						<header className="flex h-16 shrink-0 items-center justify-between gap-2 px-4 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
							<div className="flex items-center gap-1">
								<SidebarTrigger className="-ml-1" />
								<h1 className="text-2xl font-semibold">Preferences</h1>
							</div>
						</header>
						<div className="w-full p-2">{children}</div>
					</div>
				</SidebarInset>
			</SidebarProvider>
		</ProtectedRoute>
	);
}
