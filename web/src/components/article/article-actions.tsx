import { Button } from "@/components/ui/button";
import { Calendar as CalendarComponent } from "@/components/ui/calendar";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Label } from "@/components/ui/label";
import {
	Popover,
	PopoverContent,
	PopoverTrigger,
} from "@/components/ui/popover";
import { Switch } from "@/components/ui/switch";
import {
	Tooltip,
	TooltipContent,
	TooltipProvider,
	TooltipTrigger,
} from "@/components/ui/tooltip";
import {
	SiGooglechrome,
	SiJinja,
	SiOpenai,
} from "@icons-pack/react-simple-icons";
import {
	Clock,
	Copy,
	Globe,
	Link2,
	Link2Off,
	RefreshCw,
	Share,
	Trash2,
} from "lucide-react";
import { useState } from "react";

interface ArticleActionsProps {
	onDelete: () => Promise<void>;
	onRefetch: (type: string, isProxyImage: boolean) => Promise<void>;
	onRegenerateSummary: () => Promise<void>;
	onShare: () => Promise<void>;
	onUnshare: () => Promise<void>;
	isLoading: boolean;
	shareStatus?: {
		isShared: boolean;
		isExpired: boolean;
	};
	copyLink?: () => Promise<void>;
	shareExpireTime?: Date;
	onUpdateExpiration: (date: Date) => Promise<void>;
}

interface RefreshDropdownMenuProps {
	isLoading: boolean;
	onRefetch: (type: string, isProxyImage: boolean) => Promise<void>;
	onRegenerateSummary: () => Promise<void>;
}

const RefreshDropdownMenu = (props: RefreshDropdownMenuProps) => {
	const [isProxyImage, setProxyImage] = useState(false);
	return (
		<DropdownMenu>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger asChild>
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
					</TooltipTrigger>
					<TooltipContent>Refresh article content</TooltipContent>
				</Tooltip>
			</TooltipProvider>
			<DropdownMenuContent align="end">
				<DropdownMenuItem>
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<div className="flex items-center space-x-2">
									<Switch
										id="is-proxy-image"
										onClick={(e) => e.stopPropagation()}
										checked={isProxyImage}
										onCheckedChange={setProxyImage}
									/>
									<Label htmlFor="is-proxy-image">Proxy Image</Label>
								</div>
							</TooltipTrigger>
							<TooltipContent>
								Enable this will upload images to S3 to prevent CORS issues
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				</DropdownMenuItem>
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<DropdownMenuItem
								onClick={async () =>
									await props.onRefetch("http", isProxyImage)
								}
							>
								<Globe className="mr-2 h-4 w-4" /> HTTP Fetcher
							</DropdownMenuItem>
						</TooltipTrigger>
						<TooltipContent>
							Basic HTTP request to fetch the article
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<DropdownMenuItem
								onClick={async () =>
									await props.onRefetch("jinaReader", isProxyImage)
								}
							>
								<SiJinja className="mr-2 h-4 w-4" /> Jina Reader
							</DropdownMenuItem>
						</TooltipTrigger>
						<TooltipContent>
							Use Jina Reader to extract article content
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<DropdownMenuItem
								onClick={async () =>
									await props.onRefetch("browser", isProxyImage)
								}
							>
								<SiGooglechrome className="mr-2 h-4 w-4" /> Browser Fetcher
							</DropdownMenuItem>
						</TooltipTrigger>
						<TooltipContent>
							Use headless browser to fetch JavaScript rendered content
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<DropdownMenuItem
								onClick={async () => await props.onRegenerateSummary()}
							>
								<SiOpenai className="mr-2 h-4 w-4" /> Genrate Summary
							</DropdownMenuItem>
						</TooltipTrigger>
						<TooltipContent>Regenerate article summary using AI</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

interface ShareDropdownMenuProps {
	isShared?: boolean;
	isExpired?: boolean;
	onShare: () => Promise<void>;
	onUnshare: () => Promise<void>;
	copyLink?: () => Promise<void>;
	shareExpireTime?: Date;
	onUpdateExpiration?: (date: Date) => Promise<void>;
}

const ExpirationOptions = ({
	currentDate,
	onSelect,
}: {
	currentDate?: Date;
	onSelect: (date: Date) => void;
}) => {
	const createDate = (days: number) => {
		const date = new Date();
		date.setDate(date.getDate() + days);
		return date;
	};

	return (
		<div className="flex flex-col p-2 gap-1">
			<div className="mb-2 text-sm text-muted-foreground">
				Set expiration date
			</div>
			<Button
				variant="ghost"
				className="justify-start"
				onClick={() => onSelect(createDate(7))}
			>
				<Clock className="mr-2 h-4 w-4" />1 week
			</Button>
			<Button
				variant="ghost"
				className="justify-start"
				onClick={() => onSelect(createDate(30))}
			>
				<Clock className="mr-2 h-4 w-4" />1 month
			</Button>
			<Button
				variant="ghost"
				className="justify-start"
				onClick={() => onSelect(new Date("9999-12-31"))}
			>
				<Clock className="mr-2 h-4 w-4" />
				Never expires
			</Button>
			<div className="border-t my-2" />
			<div className="text-sm text-muted-foreground mb-1">Custom date</div>
			<CalendarComponent
				mode="single"
				selected={currentDate}
				onSelect={(date) => date && onSelect(date)}
				disabled={(date) => date < new Date()}
				initialFocus
			/>
		</div>
	);
};

const ShareDropdownMenu = (props: ShareDropdownMenuProps) => {
	const formatExpirationText = (date?: Date) => {
		if (!date) return "Set expiration";
		if (date.getFullYear() === 9999) return "Never expires";
		return `Expires: ${date.toLocaleDateString()}`;
	};

	return (
		<DropdownMenu>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger asChild>
						<DropdownMenuTrigger asChild>
							<Button
								variant="ghost"
								size="icon"
								className={`transition-all hover:scale-105 ${
									props.isShared
										? props.isExpired
											? "text-muted-foreground"
											: "text-primary"
										: ""
								}`}
							>
								{props.isShared ? (
									props.isExpired ? (
										<Link2Off className="h-4 w-4" />
									) : (
										<Link2 className="h-4 w-4" />
									)
								) : (
									<Share className="h-4 w-4" />
								)}
							</Button>
						</DropdownMenuTrigger>
					</TooltipTrigger>
					<TooltipContent>
						{props.isShared
							? props.isExpired
								? "Shared link expired"
								: "Article is shared"
							: "Share article"}
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
			<DropdownMenuContent align="end">
				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<DropdownMenuItem
								onClick={async () =>
									props.isShared
										? await props.onUnshare()
										: await props.onShare()
								}
							>
								{props.isShared ? (
									<>
										<Link2Off className="mr-2 h-4 w-4" /> Unshare
									</>
								) : (
									<>
										<Share className="mr-2 h-4 w-4" /> Share article
									</>
								)}
							</DropdownMenuItem>
						</TooltipTrigger>
						<TooltipContent>
							{props.isShared
								? "Remove shared access"
								: "Create a shareable link"}
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>

				{props.isShared && !props.isExpired && props.copyLink && (
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<DropdownMenuItem onClick={props.copyLink}>
									<Copy className="mr-2 h-4 w-4" /> Copy link
								</DropdownMenuItem>
							</TooltipTrigger>
							<TooltipContent>Copy shared link to clipboard</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				)}

				{props.isShared && props.onUpdateExpiration && (
					<Popover modal={false}>
						<PopoverTrigger asChild>
							<DropdownMenuItem
								// Prevent dropdown from auto-closing when clicked
								onSelect={(e) => {
									e.preventDefault();
									e.stopPropagation();
								}}
								className="cursor-pointer"
							>
								<Clock className="mr-2 h-4 w-4" />
								{props.isExpired
									? "Link expired"
									: formatExpirationText(props.shareExpireTime)}
							</DropdownMenuItem>
						</PopoverTrigger>
						<PopoverContent
							className="w-[300px] p-0"
							// Prevent focus from snapping back or popover from closing
							onCloseAutoFocus={(e) => e.preventDefault()}
							// Also prevent interactions outside from closing the popover
							onInteractOutside={(e) => {
								e.preventDefault();
							}}
						>
							<ExpirationOptions
								currentDate={props.shareExpireTime}
								onSelect={props.onUpdateExpiration}
							/>
						</PopoverContent>
					</Popover>
				)}
			</DropdownMenuContent>
		</DropdownMenu>
	);
};

export const ArticleActions: React.FC<ArticleActionsProps> = ({
	onDelete,
	onRefetch,
	onRegenerateSummary,
	onShare,
	onUnshare,
	isLoading,
	shareStatus,
	copyLink,
	shareExpireTime,
	onUpdateExpiration,
}) => {
	return (
		<div className="flex justify-end flex-wrap items-center gap-1 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 py-1">
			<RefreshDropdownMenu
				isLoading={isLoading}
				onRefetch={onRefetch}
				onRegenerateSummary={onRegenerateSummary}
			/>
			<ShareDropdownMenu
				isShared={shareStatus?.isShared}
				isExpired={shareStatus?.isExpired}
				onShare={onShare}
				onUnshare={onUnshare}
				copyLink={copyLink}
				shareExpireTime={shareExpireTime}
				onUpdateExpiration={onUpdateExpiration}
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
