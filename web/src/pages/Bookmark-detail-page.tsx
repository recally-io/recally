import { ArticleReader } from "@/components/content-reader";
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
		<div className="container mx-auto p-4 max-w-4xl">
			<ArticleReader bookmark={bookmark} />
		</div>
	);
}
