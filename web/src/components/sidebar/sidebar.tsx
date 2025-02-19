import type * as React from "react";

import AddBookmarkModal from "@/components/bookmarks/add-bookmark";
import { NavUser } from "@/components/sidebar-nav-user";
import ThemeToggle from "@/components/theme-toggle";
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarGroup,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarRail,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Settings2 } from "lucide-react";

export function SidebarComponent({
	children,
	props,
}: {
	children?: React.ReactElement;
	props?: React.ComponentProps<typeof Sidebar>;
}) {
	return (
		<Sidebar collapsible="icon" {...props}>
			<SidebarHeader>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton size="lg" asChild>
							<Link
								to={ROUTES.BOOKMARKS}
								search={{
									page: 1,
									filters: [],
									query: "",
								}}
							>
								<div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-background text-foreground">
									<img src="/logo.svg" alt="Recally logo" className="size-4" />
								</div>
								<div className="flex flex-col gap-0.5 leading-none">
									<span className="font-semibold">Recally</span>
								</div>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
					<SidebarMenuItem>
						<SidebarMenuButton asChild tooltip="Trigger Sidebar">
							<div className="flex items-center gap-2">
								<SidebarTrigger />
								<span>Trigger Sidebar</span>
							</div>
						</SidebarMenuButton>
					</SidebarMenuItem>
					<SidebarMenuItem>
						<SidebarMenuButton asChild tooltip="Add Bookmark">
							<div className="flex items-center gap-2">
								<AddBookmarkModal />
								<span>Add Bookmark</span>
							</div>
						</SidebarMenuButton>
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarHeader>
			<SidebarContent>
				{children}
				<SidebarGroup className="mt-auto">
					<SidebarMenu>
						<SidebarMenuItem key="preferences">
							<SidebarMenuButton asChild>
								<Link to={ROUTES.SETTINGS}>
									<Settings2 />
									<span>Preferences</span>
								</Link>
							</SidebarMenuButton>
						</SidebarMenuItem>
						<ThemeToggle />
					</SidebarMenu>
				</SidebarGroup>
			</SidebarContent>
			<SidebarFooter>
				<NavUser />
			</SidebarFooter>
			<SidebarRail />
		</Sidebar>
	);
}
