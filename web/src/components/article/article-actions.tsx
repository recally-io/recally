import { Chrome, Database, Globe, RefreshCw, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";

// Define fetcher types
export type FetcherType = "http" | "jina" | "browser";

interface ArticleActionsProps {
	onDelete: () => Promise<void>;
	onRefetch: (type: FetcherType) => Promise<void>;
	isLoading: boolean;
}

export const ArticleActions: React.FC<ArticleActionsProps> = ({
	onDelete,
	onRefetch,
	isLoading,
}) => {
	return (
		<div className="flex justify-end flex-wrap items-center gap-1 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 py-1">
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger asChild>
						<DropdownMenu>
							<DropdownMenuTrigger asChild>
								<Button
									variant="ghost"
									size="icon"
									disabled={isLoading}
									className="transition-all hover:scale-105"
								>
									<RefreshCw
										className={`h-4 w-4 ${isLoading ? "animate-spin" : ""}`}
									/>
								</Button>
							</DropdownMenuTrigger>
							<DropdownMenuContent align="end">
								<DropdownMenuItem onClick={async () => await onRefetch("http")}>
									<Globe className="mr-2 h-4 w-4" /> HTTP Fetcher
								</DropdownMenuItem>
								<DropdownMenuItem onClick={async () => await onRefetch("jina")}>
									<Database className="mr-2 h-4 w-4" /> Jina Fetcher
								</DropdownMenuItem>
								<DropdownMenuItem
									onClick={async () => await onRefetch("browser")}
								>
									<Chrome className="mr-2 h-4 w-4" /> Browser Fetcher
								</DropdownMenuItem>
							</DropdownMenuContent>
						</DropdownMenu>
					</TooltipTrigger>
					<TooltipContent>Refetch article</TooltipContent>
				</Tooltip>

				<Tooltip>
					<TooltipTrigger asChild>
						<Button
							onClick={async () => await onDelete()}
							variant="ghost"
							size="icon"
							className="text-destructive hover:bg-destructive/10 transition-all hover:scale-105"
						>
							<Trash2 className="h-4 w-4" />
						</Button>
					</TooltipTrigger>
					<TooltipContent>Delete article</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		</div>
	);
};
