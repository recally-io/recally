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
import { useState } from "react";
interface BookmarkListProps {
	bookmarks: Bookmark[];
}

export default function BookmarkList({ bookmarks }: BookmarkListProps) {
	const [searchQuery, setSearchQuery] = useState("");

	const filteredBookmarks = bookmarks.filter(
		(bookmark) =>
			bookmark.title?.toLowerCase().includes(searchQuery.toLowerCase()) ||
			bookmark.url.toLowerCase().includes(searchQuery.toLowerCase()) ||
			bookmark.metadata?.tags?.some((tag) =>
				tag.toLowerCase().includes(searchQuery.toLowerCase()),
			),
	);

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

			<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{filteredBookmarks.map((bookmark) => (
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
				))}
			</div>

			{filteredBookmarks.length === 0 && (
				<div className="text-center py-12">
					<p className="text-muted-foreground">No bookmarks found</p>
				</div>
			)}
		</div>
	);
}
