import MarkdownRenderer from "@/components/markdown-render";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import type { Bookmark as BookmarkType } from "@/lib/apis/bookmarks";
import { Bookmark, Share2, ThumbsUp } from "lucide-react";
import type React from "react";
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/components/ui/accordion";
interface ArticleReaderProps {
	bookmark: BookmarkType;
}

export const ArticleReader: React.FC<ArticleReaderProps> = ({ bookmark }) => {
	const { title, summary, content, metadata } = bookmark;

	const handleShare = () => {
		// Implement share functionality
		console.log("Sharing article:", bookmark.url);
	};

	const handleLike = () => {
		// Implement like functionality
		console.log("Liking article:", bookmark.id);
	};

	const handleSaveBookmark = () => {
		// Implement save bookmark functionality
		console.log("Saving bookmark:", bookmark.id);
	};

	return (
		<div className="container mx-auto px-4 py-8 max-w-screen-lg">
			<Card className="p-8 shadow-lg">
				<h1 className="text-3xl font-bold mb-4">{title}</h1>

				{/* Summary Section using Accordion */}
				<Accordion type="single" collapsible className="mb-8">
					<AccordionItem value="summary" className="border-none">
						<AccordionTrigger className="bg-background rounded-t-lg px-6 py-4 hover:no-underline">
							<h2 className="text-lg font-semibold">Summary</h2>
						</AccordionTrigger>
						<AccordionContent className="bg-muted rounded-b-lg px-6 pb-6">
							<div className="prose dark:prose-invert prose-sm max-w-none pt-4">
								<p style={{ whiteSpace: "pre-line" }}>{summary}</p>
							</div>
						</AccordionContent>
					</AccordionItem>
				</Accordion>

				{/* {metadata?.tags  && (
          <div className="flex flex-wrap gap-2 mb-6">
            {metadata.tags.map((tag, index) => (
              <Badge key={index} variant="secondary">
                {tag}
              </Badge>
            ))}
          </div>
        )} */}

				<div className="flex space-x-4 mb-8">
					<Button onClick={handleShare} variant="outline">
						<Share2 className="mr-2 h-4 w-4" /> Share
					</Button>
					<Button onClick={handleLike} variant="outline">
						<ThumbsUp className="mr-2 h-4 w-4" /> Like
					</Button>
					<Button onClick={handleSaveBookmark} variant="outline">
						<Bookmark className="mr-2 h-4 w-4" /> Save
					</Button>
				</div>

				{metadata?.highlights && (
					<div className="mb-8">
						<h2 className="text-xl font-semibold mb-4">Highlights</h2>
						{metadata.highlights.map((highlight) => (
							<div
								key={highlight.id}
								className="bg-yellow-100 dark:bg-yellow-900/30 p-4 rounded-md mb-4"
							>
								<p className="text-gray-800 dark:text-gray-200">
									{highlight.text}
								</p>
								{highlight.note && (
									<p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
										Note: {highlight.note}
									</p>
								)}
							</div>
						))}
					</div>
				)}

				{/* Main Content */}
				<div className="mt-8 border-t pt-8">
					<div className="prose dark:prose-invert prose-lg max-w-none">
						<MarkdownRenderer content={content ?? ""} />
					</div>
				</div>
			</Card>
		</div>
	);
};
