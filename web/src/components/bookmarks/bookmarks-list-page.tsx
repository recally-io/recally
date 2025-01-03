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
import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { useBookmarkMutations, useBookmarks } from "@/lib/apis/bookmarks";
import { useRouter } from "@tanstack/react-router";
import { List, Loader2, PlusCircle, Table } from "lucide-react";
import { useState } from "react";

type View = "grid" | "list";

export type BookmarkSearch = {
	page: number;
	// filter: site:github.com,category:url,tag:tag1
	filter: string;
	// query: search query
	query: string;
};

export default function BookmarksListView({
	search,
}: { search: BookmarkSearch }) {
	const limit = 12; // max 2 columns
	const offset = (search.page - 1) * limit;

	const router = useRouter();
	const { data, isLoading } = useBookmarks(
		limit,
		offset,
		search.filter,
		search.query,
	);
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

	const handlePageChange = (page: number) => {
		router.navigate({
			search: (prev) => ({ ...prev, page: page }),
		});
	};

	const handleSearch = (query: string) => {
		router.navigate({
			search: (prev) => ({ ...prev, query: query, page: 1 }),
		});
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
							currentPage={search.page}
							onPageChange={handlePageChange}
							onSearch={handleSearch}
							itemsPerPage={limit}
						/>
					)}
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
