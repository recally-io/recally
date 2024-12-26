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
import { Highlighter } from "lucide-react";
import { Link } from "react-router-dom";

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
					<Link to={`${ROUTES.BOOKMARKS}/${bookmark.id}`} className="block">
						{bookmark.metadata?.image && (
							<img
								src={bookmark.metadata.image}
								alt={bookmark.title}
								className="w-full h-48 object-cover"
							/>
						)}
						<CardHeader>
							<CardTitle className="flex items-center justify-between gap-2">
								<span className="flex items-center gap-2 truncate">
									{bookmark.title}
									{bookmark.metadata?.highlights?.length ? (
										<Highlighter className="h-4 w-4 text-yellow-500" />
									) : null}
								</span>
							</CardTitle>
							<CardDescription className="truncate">
								{bookmark.url}
							</CardDescription>
						</CardHeader>
						<CardContent>
							<div className="flex flex-wrap gap-2">
								{bookmark.metadata?.tags?.map((tag) => (
									<Badge key={tag} variant="secondary">
										{tag}
									</Badge>
								))}
							</div>
						</CardContent>
					</Link>
				</Card>
			))}
		</div>
	);
}
