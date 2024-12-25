import { Bookmark, ChevronRight, Rss, Settings2 } from "lucide-react";

import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
	SidebarGroup,
	SidebarGroupLabel,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubButton,
	SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";

const items = [
	{
		title: "Bookmarks",
		url: "#bookmark",
		icon: Bookmark,
		isActive: true,
		items: [
			{
				title: "Articles",
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=article`,
			},
			{
				title: "EPUBs",
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=epub`,
			},
			{
				title: "PDFs",
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=pdf`,
			},
			{
				title: "Videos",
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=video`,
			},
			{
				title: "Podcasts",
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=podcast`,
			},
		],
	},
	{
		title: "Subscriptions",
		url: `${ROUTES.BOOKMARKS}?type=feed`,
		icon: Rss,
		items: [],
	},
	{
		title: "Settings",
		url: "#Settings",
		icon: Settings2,
		items: [
			{
				title: "General",
				url: `${ROUTES.SETTINGS}#General`,
			},
		],
	},
];

export function NavMain() {
	return (
		<SidebarGroup>
			<SidebarGroupLabel>Platform</SidebarGroupLabel>
			<SidebarMenu>
				{items.map((item) => (
					<Collapsible
						key={item.title}
						asChild
						defaultOpen={item.isActive}
						className="group/collapsible"
					>
						<SidebarMenuItem>
							<CollapsibleTrigger asChild>
								<SidebarMenuButton tooltip={item.title}>
									{item.icon && <item.icon />}
									<span>{item.title}</span>
									<ChevronRight className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90" />
								</SidebarMenuButton>
							</CollapsibleTrigger>
							<CollapsibleContent>
								<SidebarMenuSub>
									{item.items?.map((subItem) => (
										<SidebarMenuSubItem key={subItem.title}>
											<SidebarMenuSubButton asChild>
												<a href={subItem.url}>
													<span>{subItem.title}</span>
												</a>
											</SidebarMenuSubButton>
										</SidebarMenuSubItem>
									))}
								</SidebarMenuSub>
							</CollapsibleContent>
						</SidebarMenuItem>
					</Collapsible>
				))}
			</SidebarMenu>
		</SidebarGroup>
	);
}
