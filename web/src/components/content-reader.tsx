import {
	ArticleActions,
	type FetcherType,
} from "@/components/article/article-actions";
import { ArticleHeader } from "@/components/article/article-header";
import { ArticleSummary } from "@/components/article/article-summary";
import MarkdownRenderer from "@/components/markdown-render";
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useToast } from "@/hooks/use-toast";
import {
	type Bookmark as BookmarkType,
	useBookmarkMutations,
} from "@/lib/apis/bookmarks";
import type React from "react";
import { useState } from "react";

interface ArticleReaderProps {
	bookmark: BookmarkType;
	onDelete?: (id: string) => Promise<void>;
	onRefetch?: (id: string, fetcherType: FetcherType) => Promise<void>;
	onRegenerateSummary?: (id: string) => Promise<void>;
}

export const ArticleReader: React.FC<ArticleReaderProps> = ({ bookmark }) => {
	const { toast } = useToast();
	const [isLoading, setIsLoading] = useState(false);
	const [showDeleteDialog, setShowDeleteDialog] = useState(false);

	const { deleteBookmark, refreshBookmark } = useBookmarkMutations();

	const handleRefetch = async (fetcherType: FetcherType) => {
		try {
			setIsLoading(true);
			await refreshBookmark(bookmark.id, {
				fetcher: fetcherType,
				regenerate_summary: false,
			});
			toast({
				title: "Success",
				description: `Article refetched using ${fetcherType} fetcher`,
			});
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to refetch article",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const handleRegenerateSummary = async () => {
		try {
			await refreshBookmark(bookmark.id, {
				regenerate_summary: true,
			});
			toast({
				title: "Success",
				description: "Summary regenerated successfully",
			});
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to regenerate summary",
				variant: "destructive",
			});
		}
	};

	const handleDelete = async () => {
		try {
			setIsLoading(true);
			await deleteBookmark(bookmark.id);
			toast({
				title: "Success",
				description: "Bookmark deleted successfully",
			});
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to delete bookmark",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
			setShowDeleteDialog(false);
		}
	};

	return (
		<div className="container mx-auto px-4 py-8 ">
			<ArticleHeader
				title={bookmark.title ?? "Title"}
				url={bookmark.url}
				publishedAt={bookmark.created_at}
				readingTime={"5 min read"}
				tags={bookmark.metadata?.tags}
			/>

			<ArticleActions
				onDelete={async () => setShowDeleteDialog(true)}
				onRefetch={handleRefetch}
				isLoading={isLoading}
			/>
			<div className="my-8">
				<ArticleSummary
					summary={bookmark.summary ?? ""}
					onRegenerateSummary={handleRegenerateSummary}
					isLoading={isLoading}
				/>
			</div>

			{/* Delete Confirmation Dialog */}
			<AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Are you sure?</AlertDialogTitle>
						<AlertDialogDescription>
							This action cannot be undone. This will permanently delete the
							bookmark and remove it from your library.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction
							onClick={handleDelete}
							className="bg-destructive text-destructive-foreground"
						>
							Delete
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>

			{/* Main Content */}
			<div className="mt-8 border-t pt-8">
				<div className="prose dark:prose-invert prose-lg max-w-none">
					<MarkdownRenderer content={bookmark.content ?? ""} />
				</div>
			</div>
			{/* </Card> */}
		</div>
	);
};
