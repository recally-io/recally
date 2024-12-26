import { BookmarkSidebarContent } from "@/components/bookmark-sidebar-content";
import { BookmarksSidebar } from "@/components/bookmarks-sidebar";
import { ArticleReader } from "@/components/content-reader";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { useBookmark } from "@/lib/apis/bookmarks";
import { useParams } from "react-router-dom";

export default function BookmarkDetailPage() {
	const { id } = useParams<{ id: string }>();
	const { data: bookmark, error } = useBookmark(id!);

	if (error) {
		return <div className="container mx-auto p-4">Error loading bookmark</div>;
	}

	if (!bookmark) {
		return <div className="container mx-auto p-4">Loading...</div>;
	}

	return (
		<SidebarProvider>
			<BookmarksSidebar>
				<BookmarkSidebarContent />
			</BookmarksSidebar>
			<SidebarInset>
				<div className="flex flex-col h-full">
					<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
						<div className="flex items-center gap-1 px-4">
							<SidebarTrigger className="-ml-1" />
						</div>
					</header>
					<ArticleReader bookmark={bookmark} />
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
