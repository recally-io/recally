import AddBookmarkModal from "@/components/AddBookmarkModal";
import BookmarkGrid from "@/components/BookmarkGrid";
import BookmarkList from "@/components/BookmarkList";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useBookmarkMutations, useBookmarks } from "@/lib/apis/bookmarks";
import { Grid, List, PlusCircle } from "lucide-react";
import { useState } from "react";

export default function HomePage() {
	const { data: bookmarks = [] } = useBookmarks();
	const { createBookmark } = useBookmarkMutations();
	const [view, setView] = useState<"list" | "grid">("list");
	const [isModalOpen, setIsModalOpen] = useState(false);
	const [searchTerm, setSearchTerm] = useState("");

	const filteredBookmarks = Array.isArray(bookmarks)
		? bookmarks.filter(
				(bookmark) =>
					bookmark.url.toLowerCase().includes(searchTerm.toLowerCase()) ||
					bookmark.summary?.toLowerCase().includes(searchTerm.toLowerCase()),
			)
		: [];

	const addBookmark = async (newBookmark: {
		title: string;
		url: string;
		tags: string[];
	}) => {
		try {
			await createBookmark({
				url: newBookmark.url,
				metadata: {
					tags: newBookmark.tags,
				},
			});
			setIsModalOpen(false);
		} catch (error) {
			console.error("Failed to create bookmark:", error);
		}
	};

	return (
		<div className="container mx-auto p-4 max-w-4xl">
			<div className="mb-4 flex gap-4">
				<Input
					type="text"
					placeholder="Search bookmarks..."
					value={searchTerm}
					onChange={(e) => setSearchTerm(e.target.value)}
					className="flex-grow"
				/>
				<Button
					variant="outline"
					onClick={() => setView("list")}
					aria-label="List view"
				>
					<List className="h-4 w-4" />
				</Button>
				<Button
					variant="outline"
					onClick={() => setView("grid")}
					aria-label="Grid view"
				>
					<Grid className="h-4 w-4" />
				</Button>
				<Button onClick={() => setIsModalOpen(true)}>
					<PlusCircle className="h-4 w-4" />
				</Button>
			</div>

			{view === "list" ? (
				<BookmarkList bookmarks={filteredBookmarks} />
			) : (
				<BookmarkGrid bookmarks={filteredBookmarks} />
			)}

			<AddBookmarkModal
				isOpen={isModalOpen}
				onClose={() => setIsModalOpen(false)}
				onAdd={addBookmark}
			/>
		</div>
	);
}
