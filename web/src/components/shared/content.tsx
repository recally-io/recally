import { ArticleHeader } from "@/components/article/article-header";
import { ArticleSummary } from "@/components/article/article-summary";
import MarkdownRenderer from "@/components/markdown-render";
import { AuthBanner } from "@/components/shared/auth-banner";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { useSharedContent } from "@/lib/apis/bookmarks";
import type React from "react";

interface SharedArticleReaderProps {
	id: string;
}

export const SharedArticleReader: React.FC<SharedArticleReaderProps> = ({
	id,
}) => {
	const { data, isLoading } = useSharedContent(id);

	if (isLoading || !data) {
		return <div>Loading...</div>;
	}

	const bookmark = data;
	return (
		<>
			<AuthBanner />
			<div className="container mx-auto p-4 max-w-4xl">
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
					<ArticleSummary summary={bookmark.content.summary} />
				)}

				{/* Main Content */}
				<div className="prose dark:prose-invert prose-lg max-w-none">
					<MarkdownRenderer content={bookmark.content.content ?? ""} />
				</div>
			</div>
		</>
	);
};
