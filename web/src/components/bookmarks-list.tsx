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

interface BookmarkListProps {
	bookmarks: Bookmark[];
}

export default function BookmarkList({ bookmarks }: BookmarkListProps) {
	return (
		<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
			{bookmarks.map((bookmark) => (
				<Card
					key={bookmark.id}
					className="overflow-hidden transition-transform transform hover:-translate-y-1 mx-2"
				>
					<a
						href={`${ROUTES.BOOKMARKS}?id=${bookmark.id}`}
						rel="noreferrer"
					>
						{bookmark.metadata?.image && (
							<img
								src={bookmark.metadata.image}
								alt={bookmark.title}
								className="w-full h-48 object-cover"
							/>
						)}
						<CardHeader>
							<CardTitle className="flex items-center justify-between gap-2">
								{/* <a
									href={`${ROUTES.BOOKMARKS}?id=${bookmark.id}`}
									target="_blank"
									rel="noreferrer"
								> */}
								<span className="flex items-center gap-2 truncate">
									{bookmark.title}
								</span>
								{/* </a> */}
							</CardTitle>
							<CardDescription className="truncate">
								{bookmark.url}
							</CardDescription>
						</CardHeader>
					</a>
					<CardContent>
						<div className="flex flex-wrap gap-2">
							{bookmark.metadata?.tags?.map((tag) => (
								<Badge key={tag} variant="secondary">
									{tag}
								</Badge>
							))}
						</div>
					</CardContent>
				</Card>
			))}
		</div>
	);
}
