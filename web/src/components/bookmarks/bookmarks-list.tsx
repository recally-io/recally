import { Badge } from "@/components/ui/badge";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import type { Bookmark } from "@/lib/apis/bookmarks";
import { ROUTES } from "@/lib/router";
import { Link } from "@tanstack/react-router";
import { Search } from "lucide-react";
import { useEffect, useState } from "react";

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
	view: "grid" | "list";
	total: number;
	currentPage: number;
	onPageChange: (page: number) => void;
	itemsPerPage: number;
}

export default function BookmarkList({
	bookmarks,
	total,
	view,
	currentPage,
	onPageChange,
	itemsPerPage,
}: BookmarkListProps) {
	const [searchQuery, setSearchQuery] = useState("");
	const totalPages = Math.ceil(total / itemsPerPage);

	// Reset to first page when search query changes
	useEffect(() => {
		onPageChange(1);
	}, [searchQuery]);

	const gridView = (bookmark: Bookmark) => {
		return (
			<Card
				key={bookmark.id}
				className="overflow-hidden transition-all hover:shadow-lg hover:-translate-y-1"
			>
				<Link to={ROUTES.BOOKMARK_DETAIL} params={{ id: bookmark.id }}>
					{bookmark.metadata?.image && (
						<div className="relative h-48 overflow-hidden">
							<img
								src={bookmark.metadata.image}
								alt={bookmark.title}
								className="w-full h-full object-cover"
							/>
						</div>
					)}
					<CardHeader>
						<CardTitle className="line-clamp-2 text-lg">
							{bookmark.title}
						</CardTitle>
						<CardDescription className="line-clamp-1">
							{bookmark.url}
						</CardDescription>
					</CardHeader>
				</Link>
				<CardContent>
					<div className="flex flex-wrap gap-1.5">
						{bookmark.metadata?.tags?.map((tag) => (
							<Badge key={tag} variant="secondary" className="text-xs">
								{tag}
							</Badge>
						))}
					</div>
				</CardContent>
			</Card>
		);
	};

	const listView = (bookmark: Bookmark) => {
		return (
			<div key={bookmark.id} className="p-4 border rounded-md flex gap-4">
				{bookmark.metadata?.image && (
					<div className="relative h-24 w-24 flex-shrink-0 overflow-hidden rounded-md">
						<img
							src={bookmark.metadata.image}
							alt={bookmark.title}
							className="w-full h-full object-cover"
						/>
					</div>
				)}
				<div className="flex-grow">
					<Link to={ROUTES.BOOKMARK_DETAIL} params={{ id: bookmark.id }}>
						<h3 className="text-lg font-semibold line-clamp-1">
							{bookmark.title}
						</h3>
					</Link>
					<p className="text-sm text-muted-foreground line-clamp-1 break-all">
						{bookmark.url}
					</p>
					<div className="mt-2 flex flex-wrap gap-1">
						{bookmark.metadata?.tags?.map((tag) => (
							<Badge key={tag} variant="secondary" className="text-xs">
								{tag}
							</Badge>
						))}
					</div>
				</div>
			</div>
		);
	};

	return (
		<div className="container mx-auto px-4 py-6 space-y-6">
			<div className="w-full  relative">
				<Input
					type="search"
					placeholder="Search by title, URL, or tags..."
					value={searchQuery}
					onChange={(e) => setSearchQuery(e.target.value)}
					className="w-full pl-9"
				/>
				<Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
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
