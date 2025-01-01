import MarkdownRenderer from "@/components/markdown-render";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter } from "@/components/ui/card";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import { ChevronsUpDown, Clock, RefreshCw } from "lucide-react";
import { useState } from "react";

interface ArticleSummaryProps {
	summary: string;
	onRegenerate?: () => Promise<void>;
}

export const ArticleSummary = ({
	summary,
	onRegenerate,
}: ArticleSummaryProps) => {
	const [isOpen, setIsOpen] = useState(true);
	const [isRegenerating, setIsRegenerating] = useState(false);
	const readingTime = Math.ceil(summary.split(/\s+/).length / 200); // Approx. reading time in minutes

	const handleRegenerate = async () => {
		if (!onRegenerate) return;
		setIsRegenerating(true);
		await onRegenerate();
		setIsRegenerating(false);
	};

	return (
		<Card className="mb-6 bg-muted/50">
			<Collapsible open={isOpen} onOpenChange={setIsOpen}>
				<div className="flex items-center justify-between px-4 py-2">
					<div className="flex items-center gap-2">
						<h2 className="text-lg font-semibold">AI Summary</h2>
						<span className="flex items-center text-xs text-muted-foreground">
							<Clock className="mr-1 h-3 w-3" />
							{readingTime} min read
						</span>
					</div>
					<CollapsibleTrigger asChild>
						<Button variant="ghost" size="sm" className="w-9 p-0">
							<ChevronsUpDown className="h-4 w-4" />
							<span className="sr-only">Toggle summary</span>
						</Button>
					</CollapsibleTrigger>
				</div>

				<CollapsibleContent>
					<CardContent className="pt-2">
						<div className="prose dark:prose-invert prose-sm max-w-none prose-h1:text-xl">
							<MarkdownRenderer content={summary} />
						</div>
					</CardContent>
					{onRegenerate && (
						<CardFooter className="justify-end py-2">
							<Button
								variant="ghost"
								size="sm"
								onClick={handleRegenerate}
								disabled={isRegenerating}
							>
								<RefreshCw
									className={`mr-2 h-4 w-4 ${isRegenerating ? "animate-spin" : ""}`}
								/>
								Regenerate
							</Button>
						</CardFooter>
					)}
				</CollapsibleContent>
			</Collapsible>
		</Card>
	);
};
