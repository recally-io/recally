import {
	ArticleActions,
	type FetcherType,
} from "@/components/article/article-actions";
import { BookmarkSidebar } from "@/components/bookmark-sidebar";
import { ArticleReader } from "@/components/content-reader";
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
import {
	SidebarInset,
	SidebarProvider,
	SidebarTrigger,
} from "@/components/ui/sidebar";
import { useToast } from "@/hooks/use-toast";
import { useBookmark, useBookmarkMutations } from "@/lib/apis/bookmarks";
import { useState } from "react";

export default function BookmarkDetailPage({ id }: { id: string }) {
	const { toast } = useToast();
	const [isLoading, setIsLoading] = useState(false);
	const [showDeleteDialog, setShowDeleteDialog] = useState(false);
	const { deleteBookmark, refreshBookmark } = useBookmarkMutations();

	const { data: bookmark, error } = useBookmark(id);
	if (error) {
		return <div className="container mx-auto p-4">Error loading bookmark</div>;
	}
	if (!bookmark) {
		return <div className="container mx-auto p-4">Loading...</div>;
	}

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
		<SidebarProvider>
			<BookmarkSidebar />
			<SidebarInset className="overflow-auto">
				<div className="flex flex-col h-full w-full">
					<header className="flex h-12 shrink-0 items-center justify-between bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 border-b">
						<div className="flex items-center px-4">
							<SidebarTrigger className="-ml-1" />
						</div>
						<div className="flex items-center px-4">
							<ArticleActions
								onDelete={async () => setShowDeleteDialog(true)}
								onRefetch={handleRefetch}
								isLoading={isLoading}
							/>
						</div>
					</header>
					<main className="container mx-auto p-4 max-w-5xl">
						<ArticleReader bookmark={bookmark} />

						{/* Delete Confirmation Dialog */}
						<AlertDialog
							open={showDeleteDialog}
							onOpenChange={setShowDeleteDialog}
						>
							<AlertDialogContent>
								<AlertDialogHeader>
									<AlertDialogTitle>Are you sure?</AlertDialogTitle>
									<AlertDialogDescription>
										This action cannot be undone. This will permanently delete
										the bookmark and remove it from your library.
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
					</main>
				</div>
			</SidebarInset>
		</SidebarProvider>
	);
}
