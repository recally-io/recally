import { useTheme } from "@/components/theme-provider";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
	SidebarContent,
	SidebarGroup,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarMenuSub,
	SidebarMenuSubButton,
	SidebarMenuSubItem,
} from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import {
	BookImage,
	BookOpen,
	Bookmark,
	ChevronRight,
	Moon,
	Newspaper,
	Settings2,
	Sun,
} from "lucide-react";
import { Link } from "react-router-dom";

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
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=article`,
			},
			{
				title: "EPUBs",
				icon: BookOpen,
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=epub`,
			},
			{
				title: "PDFs",
				icon: BookImage,
				url: `${ROUTES.BOOKMARKS}?type=bookmark&category=pdf`,
			},
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
	// {
	// 	title: "Subscriptions",
	// 	url: `${ROUTES.BOOKMARKS}?type=feed`,
	// 	icon: Rss,
	// 	items: [],
	// }
];

function ThemeToggle() {
	const { setTheme } = useTheme();
	return (
		<SidebarMenuItem key="theme-toggle">
			<DropdownMenu>
				<DropdownMenuTrigger asChild>
					<SidebarMenuButton>
						<Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
						<Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
						<span>Theme</span>
					</SidebarMenuButton>
				</DropdownMenuTrigger>
				<DropdownMenuContent align="end">
					<DropdownMenuItem onClick={() => setTheme("light")}>
						Light
					</DropdownMenuItem>
					<DropdownMenuItem onClick={() => setTheme("dark")}>
						Dark
					</DropdownMenuItem>
					<DropdownMenuItem onClick={() => setTheme("system")}>
						System
					</DropdownMenuItem>
				</DropdownMenuContent>
			</DropdownMenu>
		</SidebarMenuItem>
	);
}

export function NavContent() {
	return (
		<SidebarContent>
			<SidebarGroup>
				{/* <SidebarGroupLabel>Platform</SidebarGroupLabel> */}
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
													<Link to={subItem.url}>
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
	);
}
