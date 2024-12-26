import BookmarkList from "@/components/bookmarks-list";
import { AppSidebar } from "@/components/sidebar";
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
import { PlusCircle } from "lucide-react";
import { useState } from "react";

export default function BookmarkPage() {
	const { data: bookmarks = [] } = useBookmarks();
	const { createBookmark } = useBookmarkMutations();
	const [open, setOpen] = useState(false);
	const [url, setUrl] = useState("");

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
			<AppSidebar />
			<SidebarInset>
				<div className="flex flex-col h-full">
					<header className="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12">
						<div className="flex items-center gap-1 px-4">
							<SidebarTrigger className="-ml-1" />
							<AddBookmarkModal />
						</div>
					</header>
					<BookmarkList bookmarks={bookmarks} />
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
