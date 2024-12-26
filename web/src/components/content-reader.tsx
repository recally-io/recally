import type { FetcherType } from "@/components/article/article-actions";
import { ArticleHeader } from "@/components/article/article-header";
import MarkdownRenderer from "@/components/markdown-render";
import type { Bookmark as BookmarkType } from "@/lib/apis/bookmarks";
import type React from "react";

interface ArticleReaderProps {
	bookmark: BookmarkType;
	onDelete?: (id: string) => Promise<void>;
	onRefetch?: (id: string, fetcherType: FetcherType) => Promise<void>;
	onRegenerateSummary?: (id: string) => Promise<void>;
}

export const ArticleReader: React.FC<ArticleReaderProps> = ({ bookmark }) => {
	return (
		<>
			<ArticleHeader
				title={bookmark.title ?? "Title"}
				url={bookmark.url}
				publishedAt={bookmark.created_at}
			/>

			{/* Main Content */}
			<div className="prose dark:prose-invert prose-lg max-w-none">
				<MarkdownRenderer content={bookmark.content ?? ""} />
			</div>
		</>
	);
};
