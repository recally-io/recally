import { ArticleActions } from "@/components/article/article-actions";
import { ArticleReader } from "@/components/bookmarks/bookmark-content";
import { BookmarksSidebar } from "@/components/bookmarks/sidebar";
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
import {
	useBookmark,
	useBookmarkMutations,
	useShareContentMutations,
} from "@/lib/apis/bookmarks";
import { ROUTES } from "@/lib/router";
import { useState } from "react";

// Add this near the top of the file, before the component
const getShareUrl = (shareId?: string) => {
	const host = window.location.origin;
	return shareId ? `${host}/share/${shareId}` : undefined;
};

export default function BookmarkDetailPage({ id }: { id: string }) {
	const { toast } = useToast();
	const [isLoading, setIsLoading] = useState(false);
	const [showDeleteDialog, setShowDeleteDialog] = useState(false);
	const { deleteBookmark, refreshBookmark } = useBookmarkMutations();
	const { shareContent, unshareContent, updateSharedContent } =
		useShareContentMutations();

	const { data: bookmark, error } = useBookmark(id);
	if (error) {
		return <div className="container mx-auto p-4">Error loading bookmark</div>;
	}
	if (!bookmark) {
		return <div className="container mx-auto p-4">Loading...</div>;
	}

	const handleRefetch = async (fetcherType: string) => {
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
			setIsLoading(true);
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
			window.location.href = ROUTES.BOOKMARKS;
		} catch (error) {
			toast({
				title: "Error",
				description: `Failed to delete bookmark ${error instanceof Error ? error.message : "Unknown error"}`,
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
			setShowDeleteDialog(false);
		}
	};

	const handleShare = async () => {
		try {
			setIsLoading(true);
			const expiresAt = new Date();
			expiresAt.setDate(expiresAt.getDate() + 7); // Share for 7 days
			await shareContent(bookmark.id, {
				expires_at: expiresAt.toISOString(),
			});
			toast({
				title: "Success",
				description: "Article shared successfully",
			});
			handleCopyLink();
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to share article",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const handleUnshare = async () => {
		try {
			setIsLoading(true);
			await unshareContent(bookmark.id);
			toast({
				title: "Success",
				description: "Article unshared successfully",
			});
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to unshare article",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const handleCopyLink = async () => {
		try {
			setIsLoading(true);
			const shareUrl = getShareUrl(bookmark.metadata?.share?.id);
			if (shareUrl) {
				await navigator.clipboard.writeText(shareUrl);
				toast({
					title: "Success",
					description: "Share link copied to clipboard",
				});
			}
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to copy share link",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const handleUpdateExpiration = async (date: Date) => {
		try {
			setIsLoading(true);
			await updateSharedContent(bookmark.id, {
				expires_at: date.toISOString(),
			});
			toast({
				title: "Success",
				description: "Share expiration updated successfully",
			});
		} catch (error) {
			toast({
				title: "Error",
				description: "Failed to update share expiration",
				variant: "destructive",
			});
		} finally {
			setIsLoading(false);
		}
	};

	const shareStatus = bookmark.metadata?.share
		? {
				isShared: true,
				isExpired: bookmark.metadata.share.expires_at
					? new Date(bookmark.metadata.share.expires_at) < new Date()
					: false,
			}
		: {
				isShared: false,
				isExpired: false,
			};

	const shareExpireTime = bookmark.metadata?.share?.expires_at
		? new Date(bookmark.metadata.share.expires_at)
		: undefined;

	return (
		<SidebarProvider defaultOpen={false}>
			<BookmarksSidebar />
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
								onRegenerateSummary={handleRegenerateSummary}
								onShare={handleShare}
								onUnshare={handleUnshare}
								isLoading={isLoading}
								shareStatus={shareStatus}
								copyLink={
									bookmark.metadata?.share?.id ? handleCopyLink : undefined
								}
								shareExpireTime={shareExpireTime}
								onUpdateExpiration={handleUpdateExpiration}
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
