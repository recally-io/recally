import MarkdownRenderer from "@/components/markdown-render";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import {
	Collapsible,
	CollapsibleContent,
	CollapsibleTrigger,
} from "@/components/ui/collapsible";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import { ChevronDown, RefreshCw } from "lucide-react";
import { useState } from "react";
interface ArticleSummaryProps {
	summary: string;
	onRegenerateSummary: () => Promise<void>;
	isLoading?: boolean;
}

export const ArticleSummary: React.FC<ArticleSummaryProps> = ({
	summary,
	onRegenerateSummary,
	isLoading,
}) => {
	const [isOpen, setIsOpen] = useState(false);
	const [isGenerating, setIsGenerating] = useState(false);

	return (
		<Card className="bg-secondary/50 border-none shadow-sm">
			<Collapsible open={isOpen} onOpenChange={setIsOpen}>
				<div className="p-6">
					<div className="flex items-center justify-between">
						<div className="flex items-center gap-2">
							<span className="ml-2 inline-flex items-center rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary">
								AI Generated Summary
							</span>
						</div>
						<div className="flex items-center gap-2">
							<TooltipProvider>
								<Tooltip>
									<TooltipTrigger asChild>
										<Button
											variant="ghost"
											size="icon"
											className="transition-all hover:scale-105"
											onClick={async () => {
												setIsGenerating(true);
												await onRegenerateSummary();
												setIsGenerating(false);
											}}
											disabled={isLoading || isGenerating}
										>
											<RefreshCw
												className={`h-4 w-4 ${isLoading || isGenerating ? "animate-spin" : ""}`}
											/>
										</Button>
									</TooltipTrigger>
									<TooltipContent>Regenerate summary</TooltipContent>
								</Tooltip>
							</TooltipProvider>

							<CollapsibleTrigger asChild>
								<Button variant="ghost" size="sm">
									<ChevronDown
										className={`h-4 w-4 transition-transform duration-200 ${
											isOpen ? "transform rotate-180" : ""
										}`}
									/>
								</Button>
							</CollapsibleTrigger>
						</div>
					</div>
				</div>
				<CollapsibleContent>
					<div className="px-6 pb-6">
						<Card className="bg-background p-4 shadow-sm">
							{isLoading || isGenerating ? (
								<div className="flex items-center justify-center py-8">
									<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
								</div>
							) : (
								<div className="prose dark:prose-invert prose-sm max-w-none">
									<MarkdownRenderer content={summary} />
								</div>
							)}
						</Card>
					</div>
				</CollapsibleContent>
			</Collapsible>
		</Card>
	);
};
