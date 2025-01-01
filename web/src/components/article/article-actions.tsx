import { Bot, Chrome, Database, Globe, RefreshCw, Trash2 } from "lucide-react";

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
	onRegenerateSummary: () => Promise<void>;
	isLoading: boolean;
}

interface RefreshDropdownMenuProps {
	isLoading: boolean;
	onRefetch: (type: FetcherType) => Promise<void>;
	onRegenerateSummary: () => Promise<void>;
}

const RefreshDropdownMenu = (props: RefreshDropdownMenuProps) => {
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button
					variant="ghost"
					size="icon"
					disabled={props.isLoading}
					className="transition-all hover:scale-105"
				>
					<RefreshCw
						className={`h-4 w-4 ${props.isLoading ? "animate-spin" : ""}`}
					/>
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end">
				<DropdownMenuItem onClick={async () => await props.onRefetch("http")}>
					<Globe className="mr-2 h-4 w-4" /> HTTP Fetcher
				</DropdownMenuItem>
				<DropdownMenuItem onClick={async () => await props.onRefetch("jina")}>
					<Database className="mr-2 h-4 w-4" /> Jina Fetcher
				</DropdownMenuItem>
				<DropdownMenuItem
					onClick={async () => await props.onRefetch("browser")}
				>
					<Chrome className="mr-2 h-4 w-4" /> Browser Fetcher
				</DropdownMenuItem>
				<DropdownMenuItem
					onClick={async () => await props.onRegenerateSummary()}
				>
					<Bot className="mr-2 h-4 w-4" /> Genrate Summary
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

export const ArticleActions: React.FC<ArticleActionsProps> = ({
	onDelete,
	onRefetch,
	onRegenerateSummary,
	isLoading,
}) => {
	return (
		<div className="flex justify-end flex-wrap items-center gap-1 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 py-1">
			<RefreshDropdownMenu
				isLoading={isLoading}
				onRefetch={onRefetch}
				onRegenerateSummary={onRegenerateSummary}
			/>
			<TooltipProvider>
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
