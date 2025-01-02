import type * as React from "react";

import { NavUser } from "@/components/sidebar-nav-user";
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
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Bookmark, Settings2 } from "lucide-react";
import ThemeToggle from "../theme-toggle";
import { BookmarksSidebarContent } from "./bookmarks-sidebar-content";

export function BookmarksSidebar({
	...props
}: React.ComponentProps<typeof Sidebar>) {
	return (
		<Sidebar collapsible="icon" {...props}>
			<SidebarHeader>
				<SidebarMenu>
					<SidebarMenuItem>
						<SidebarMenuButton size="lg" asChild>
							<Link to={ROUTES.HOME}>
								<div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
									<Bookmark className="size-4" />
								</div>
								<div className="flex flex-col gap-0.5 leading-none">
									<span className="font-semibold">Recally</span>
								</div>
							</Link>
						</SidebarMenuButton>
					</SidebarMenuItem>
				</SidebarMenu>
			</SidebarHeader>
			<SidebarContent>
				<BookmarksSidebarContent />
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
