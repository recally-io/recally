import { ArticleHeader } from "@/components/article/article-header";
import { ArticleSummary } from "@/components/article/article-summary";
import MarkdownRenderer from "@/components/markdown-render";
import PdfViewer from "@/components/pdf-viewer";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import type { Bookmark as BookmarkType } from "@/lib/apis/bookmarks";
import { useGetFile } from "@/lib/apis/file";
import type React from "react";
import { useEffect, useState } from "react";

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
	const { trigger: getFile } = useGetFile();
	const [pdfUrl, setPdfUrl] = useState<string | null>(null);

	useEffect(() => {
		const loadPdf = async () => {
			if (bookmark.content.s3_key) {
				try {
					const fileUrl = await getFile({
						object_key: bookmark.content.s3_key,
					});
					setPdfUrl(fileUrl.url);
				} catch (error) {
					console.error("Failed to load PDF:", error);
				}
			}
		};
		loadPdf();
	}, [bookmark.content.s3_key, getFile]);

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
			{bookmark.content.type === "pdf" && pdfUrl ? (
				<PdfViewer fileUrl={pdfUrl} />
			) : (
				<div className="prose dark:prose-invert prose-lg max-w-none">
					<MarkdownRenderer content={bookmark.content.content ?? ""} />
				</div>
			)}
		</>
	);
};
