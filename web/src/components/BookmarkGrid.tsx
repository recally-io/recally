import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { Bookmark } from "@/lib/apis/bookmarks";
import { ExternalLink, Highlighter } from "lucide-react";
import { Link } from "react-router-dom";

interface BookmarkGridProps {
	bookmarks: Bookmark[];
}

export default function BookmarkGrid({ bookmarks }: BookmarkGridProps) {
	return (
		<div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
			{bookmarks.map((bookmark) => (
				<Link to={`/bookmarks/${bookmark.id}`} key={bookmark.id}>
					<Card className="h-full cursor-pointer hover:shadow-md transition-all duration-300 ease-in-out transform hover:-translate-y-1">
						{bookmark.metadata?.image && (
							<img
								src={bookmark.metadata?.image}
								alt={`Thumbnail for ${bookmark.title}`}
								className="w-full h-40 object-cover"
							/>
						)}
						<CardHeader>
							<CardTitle className="flex items-center justify-between">
								<span className="flex items-center gap-2 truncate">
									{bookmark.title}
									{bookmark.metadata?.highlights &&
										bookmark.metadata?.highlights.length > 0 && (
											<Highlighter className="h-4 w-4 text-yellow-500 flex-shrink-0" />
										)}
								</span>
								<a
									href={bookmark.url}
									target="_blank"
									rel="noopener noreferrer"
									className="text-blue-500 hover:text-blue-700 transition-colors"
									onClick={(e) => e.stopPropagation()}
								>
									<ExternalLink className="h-4 w-4" />
								</a>
							</CardTitle>
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
					</Card>
				</Link>
			))}
		</div>
	);
}
