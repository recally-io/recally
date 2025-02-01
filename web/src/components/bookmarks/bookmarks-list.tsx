import { Badge } from "@/components/ui/badge";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import type { Bookmark } from "@/lib/apis/bookmarks";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";

import SearchBox, { type SearchToken } from "@/components/bookmarks/search";
import type { BookmarkSearch, View } from "@/components/bookmarks/types";
import {
	Pagination,
	PaginationContent,
	PaginationItem,
	PaginationLink,
	PaginationNext,
	PaginationPrevious,
} from "@/components/ui/pagination";

interface BookmarkListProps {
	bookmarks: Bookmark[];
	view: View;
	search: BookmarkSearch;
	total: number;
	currentPage: number;
	onPageChange: (page: number) => void;
	onSearch: (tokens: SearchToken[], query: string) => void;
	itemsPerPage: number;
}

export default function BookmarkList({
	bookmarks,
	total,
	view,
	search,
	currentPage,
	onPageChange,
	onSearch,
	itemsPerPage,
}: BookmarkListProps) {
	const totalPages = Math.ceil(total / itemsPerPage);

	const gridView = (bookmark: Bookmark) => {
		return (
			<Card
				key={bookmark.id}
				className="overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1"
			>
				<Link to={ROUTES.BOOKMARK_DETAIL} params={{ id: bookmark.id }}>
					{bookmark.content.metadata?.cover && (
						<div className="relative h-48 overflow-hidden">
							<img
								src={bookmark.content.metadata.cover}
								alt={bookmark.content.title}
								className="w-full h-full object-cover"
							/>
						</div>
					)}
					<CardHeader>
						<CardTitle className="line-clamp-2 text-lg">
							{bookmark.content.title}
						</CardTitle>
						<CardDescription className="line-clamp-1">
							{bookmark.content.url}
						</CardDescription>
					</CardHeader>

					<CardContent>
						<div className="flex flex-wrap gap-1.5">
							{bookmark.tags?.map((tag) => (
								<Badge key={tag} variant="secondary" className="text-xs">
									{tag}
								</Badge>
							))}
						</div>
					</CardContent>
				</Link>
			</Card>
		);
	};

	const listView = (bookmark: Bookmark) => {
		return (
			<Link
				key={bookmark.id}
				to={ROUTES.BOOKMARK_DETAIL}
				params={{ id: bookmark.id }}
			>
				<div className="p-4 border rounded-md flex gap-4">
					{bookmark.content.metadata?.cover && (
						<div className="relative h-24 w-24 flex-shrink-0 overflow-hidden rounded-md">
							<img
								src={bookmark.content.metadata.cover}
								alt={bookmark.content.title}
								className="w-full h-full object-cover"
							/>
						</div>
					)}
					<div className="flex-grow">
						<h3 className="text-lg font-semibold line-clamp-1">
							{bookmark.content.title}
						</h3>
						<p className="text-sm text-muted-foreground line-clamp-1 break-all">
							{bookmark.content.url}
						</p>
						<div className="mt-2 flex flex-wrap gap-1">
							{bookmark.tags?.map((tag) => (
								<Badge key={tag} variant="secondary" className="text-xs">
									{tag}
								</Badge>
							))}
						</div>
					</div>
				</div>
			</Link>
		);
	};

	return (
		<div className="container mx-auto px-4 py-6 space-y-6">
			<div className="mb-4">
				<SearchBox onSearch={onSearch} search={search} />
			</div>

			{view === "grid" ? (
				<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6 gap-6">
					{bookmarks.map((bookmark) => gridView(bookmark))}
				</div>
			) : (
				<div className="space-y-4">
					{bookmarks.map((bookmark) => listView(bookmark))}
				</div>
			)}

			{bookmarks.length === 0 ? (
				<div className="text-center py-12">
					<p className="text-muted-foreground">No bookmarks found</p>
				</div>
			) : (
				<div className="mt-6">
					<Pagination>
						<PaginationContent>
							{currentPage !== 1 && (
								<PaginationItem>
									<PaginationPrevious
										onClick={() => onPageChange(currentPage - 1)}
									/>
								</PaginationItem>
							)}

							{Array.from({ length: totalPages }, (_, i) => i + 1).map(
								(page) => (
									<PaginationItem key={page}>
										<PaginationLink
											onClick={() => onPageChange(page)}
											isActive={currentPage === page}
										>
											{page}
										</PaginationLink>
									</PaginationItem>
								),
							)}

							{currentPage !== totalPages && (
								<PaginationItem>
									<PaginationNext
										onClick={() => onPageChange(currentPage + 1)}
									/>
								</PaginationItem>
							)}
						</PaginationContent>
					</Pagination>
				</div>
			)}
		</div>
	);
}
