import { SidebarMenuButton, SidebarTrigger } from "@/components/ui/sidebar";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Home } from "lucide-react";

export function SidebarHeaderTrigger() {
	return (
		<>
			<div className="hidden items-center justify-start md:flex"></div>
			<div className="flex items-center justify-start md:hidden">
				<SidebarMenuButton size="sm" asChild>
					<Link
						to={ROUTES.BOOKMARKS}
						search={{
							page: 1,
							filters: [],
							query: "",
						}}
					>
						<Home />
					</Link>
				</SidebarMenuButton>
				<SidebarTrigger />
			</div>
		</>
	);
}
