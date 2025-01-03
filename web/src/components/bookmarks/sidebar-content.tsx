import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubButton,
	SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Bookmark, ChevronRight, Newspaper } from "lucide-react";

const items = [
	{
		title: "Bookmarks",
		url: "#bookmark",
		icon: Bookmark,
		isActive: true,
		items: [
			{
				title: "Articles",
				icon: Newspaper,
				url: ROUTES.BOOKMARKS,
				type: "bookmark",
				category: "article",
			},
			// {
			// 	title: "EPUBs",
			// 	icon: BookOpen,
			// 	url: `${ROUTES.BOOKMARKS}?type=bookmark&category=epub`,
			// },
			// {
			// 	title: "PDFs",
			// 	icon: BookImage,
			// 	url: `${ROUTES.BOOKMARKS}?type=bookmark&category=pdf`,
			// },
			// {
			// 	title: "Videos",
			// 	url: `${ROUTES.BOOKMARKS}?type=bookmark&category=video`,
			// },
			// {
			// 	title: "Podcasts",
			// 	url: `${ROUTES.BOOKMARKS}?type=bookmark&category=podcast`,
			// },
		],
	},
];

export function BookmarksSidebarContent() {
	return (
		<SidebarGroup>
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
												<Link
													to={subItem.url}
													search={{
														type: subItem.type,
														category: subItem.category,
													}}
												>
													{subItem.icon && <subItem.icon />}
													<span>{subItem.title}</span>
												</Link>
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
