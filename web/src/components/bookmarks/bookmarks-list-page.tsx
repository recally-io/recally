import BookmarkList from "@/components/bookmarks/bookmarks-list";
import { BookmarksSidebar } from "@/components/bookmarks/sidebar";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { useBookmarkMutations, useBookmarks } from "@/lib/apis/bookmarks";
import { PlusCircle, List, Table } from "lucide-react";
import { useState } from "react";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

type View = "grid" | "list";

export default function BookmarksListView() {
	const { data: bookmarks = [] } = useBookmarks();
	const { createBookmark } = useBookmarkMutations();
	const [open, setOpen] = useState(false);
	const [url, setUrl] = useState("");
	const [view, setView] = useState<View>("grid");

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();
		if (!url) return;
		await createBookmark({ url });
		setUrl("");
		setOpen(false);
	};

	const AddBookmarkModal = () => {
		return (
			<Dialog open={open} onOpenChange={setOpen}>
				<DialogTrigger asChild>
					<Button variant="ghost" size="icon" className="h-7 w-7">
						<PlusCircle className="size-6" />
					</Button>
				</DialogTrigger>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Add New Bookmark</DialogTitle>
					</DialogHeader>
					<form onSubmit={handleSubmit} className="space-y-4">
						<Input
							placeholder="Enter URL"
							value={url}
							onChange={(e) => setUrl(e.target.value)}
						/>
						<Button type="submit">Add Bookmark</Button>
					</form>
				</DialogContent>
			</Dialog>
		);
	};

	return (
		<SidebarProvider>
			<BookmarksSidebar />
			<SidebarInset>
				<div className="flex flex-col h-full">
					<header className="flex h-16 shrink-0 items-center justify-between gap-2 px-4 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
						<div className="flex items-center gap-1">
							<SidebarTrigger className="-ml-1" />
							<AddBookmarkModal />
						</div>
						<ToggleGroup
							type="single"
							value={view}
							onValueChange={(value) => setView(value as View)}
							size="sm"
						>
							<ToggleGroupItem value="grid" aria-label="Grid">
								<Table />
							</ToggleGroupItem>
							<ToggleGroupItem value="list" aria-label="List">
								<List />
							</ToggleGroupItem>
						</ToggleGroup>
					</header>
					<BookmarkList bookmarks={bookmarks} view={view} />
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
