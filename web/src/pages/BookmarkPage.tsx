import BookmarkDetail from "@/components/BookmarkDetail";
import { useBookmark, useBookmarkMutations } from "@/lib/apis/bookmarks";
import { useParams } from "react-router-dom";

export default function BookmarkPage() {
	const { id } = useParams<{ id: string }>();
	const { data: bookmark, error } = useBookmark(id!);
	const { updateBookmark } = useBookmarkMutations();

	if (error) {
		return <div className="container mx-auto p-4">Error loading bookmark</div>;
	}

	if (!bookmark) {
		return <div className="container mx-auto p-4">Loading...</div>;
	}

	return (
		<div className="container mx-auto p-4 max-w-4xl">
			<BookmarkDetail
				bookmark={bookmark}
				onUpdateBookmark={async (id, highlights) => {
					try {
						await updateBookmark(id, {
							metadata: {
								...bookmark.metadata,
								highlights,
							},
						});
					} catch (error) {
						console.error("Failed to update bookmark:", error);
					}
				}}
			/>
		</div>
	);
}
