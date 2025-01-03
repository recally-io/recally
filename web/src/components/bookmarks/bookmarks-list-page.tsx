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
import { PlusCircle, List, Table, Loader2 } from "lucide-react";
import { useMemo, useState } from "react";
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

type View = "grid" | "list";

export default function BookmarksListView() {
	const [currentPage, setCurrentPage] = useState(1);
	const limit = 12; // max 2 columns
	const offset = useMemo(() => (currentPage - 1) * limit, [currentPage, limit]);

	const { data, isLoading } = useBookmarks(limit, offset);
	const bookmarks = data?.bookmarks ?? [];
	const total = data?.total ?? 0;

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
					{isLoading ? (
						<div className="flex items-center justify-center h-full">
							<Loader2 className="size-8 animate-spin" />
						</div>
					) : (
						<BookmarkList
							bookmarks={bookmarks}
							total={total}
							view={view}
							currentPage={currentPage}
							onPageChange={setCurrentPage}
							itemsPerPage={limit}
						/>
					)}
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
