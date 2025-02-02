import { ArticleHeader } from "@/components/article/article-header";
import { ArticleSummary } from "@/components/article/article-summary";
import MarkdownRenderer from "@/components/markdown-render";
import { Badge } from "@/components/ui/badge";
import type { Bookmark as BookmarkType } from "@/lib/apis/bookmarks";
import type React from "react";
import { Separator } from "../ui/separator";

interface ArticleReaderProps {
	bookmark: BookmarkType;
	onDelete?: (id: string) => Promise<void>;
	onRefetch?: (id: string, fetcherType: string) => Promise<void>;
	onRegenerateSummary?: (id: string) => Promise<void>;
}

export const ArticleReader: React.FC<ArticleReaderProps> = ({
	bookmark,
	onRegenerateSummary,
}) => {
	return (
		<>
			<ArticleHeader
				title={bookmark.content.title ?? "Title"}
				url={bookmark.content.url || ""}
				publishedAt={
					bookmark.content.metadata?.published_at ?? bookmark.created_at
				}
			/>

			{/* Tags */}
			{bookmark.tags && bookmark.tags.length > 0 && (
				<div className="flex flex-wrap gap-2 mt-4">
					{bookmark.tags.map((tag) => (
						<Badge key={tag} variant="secondary">
							{tag}
						</Badge>
					))}
				</div>
			)}

			<Separator className="my-4" />

			{bookmark.content.summary && (
				<ArticleSummary
					summary={bookmark.content.summary}
					onRegenerate={
						onRegenerateSummary
							? () => onRegenerateSummary(bookmark.id)
							: undefined
					}
				/>
			)}

			{/* Main Content */}
			<div className="prose dark:prose-invert prose-lg max-w-none">
				<MarkdownRenderer content={bookmark.content.content ?? ""} />
			</div>
		</>
	);
};
