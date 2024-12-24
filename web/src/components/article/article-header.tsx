interface ArticleHeaderProps {
	title: string;
	url: string;
	author?: string;
	publishedAt?: string;
	readingTime?: string;
	tags?: string[];
}

export const ArticleHeader: React.FC<ArticleHeaderProps> = ({
	title,
	url,
	publishedAt,
	readingTime,
	tags,
}) => {
	const domain = new URL(url).hostname;

	return (
		<div className="space-y-6 border-b pb-8">
			{/* Title */}
			<h1 className="text-4xl font-serif font-bold leading-tight tracking-tighter md:text-5xl">
				{title}
			</h1>

			{/* Metadata Bar */}
			<div className="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-muted-foreground">
				<a
					href={url}
					target="_blank"
					rel="noopener noreferrer"
					className="flex items-center gap-2 hover:text-primary transition-colors"
				>
					<img
						src={`https://www.google.com/s2/favicons?domain=${domain}&sz=32`}
						alt=""
						className="h-4 w-4 rounded-sm"
					/>
					<span>{domain}</span>
				</a>
				{readingTime && (
					<>
						<span className="text-muted-foreground/40">•</span>
						<span>{readingTime}</span>
					</>
				)}
				{publishedAt && (
					<>
						<span className="text-muted-foreground/40">•</span>
						<span>{publishedAt}</span>
					</>
				)}
			</div>

			{/* Tags */}
			{tags && tags.length > 0 && (
				<div className="flex flex-wrap gap-2">
					{tags.map((tag) => (
						<span
							key={tag}
							className="inline-flex items-center rounded-full bg-secondary px-2.5 py-0.5 text-xs font-medium text-secondary-foreground"
						>
							{tag}
						</span>
					))}
				</div>
			)}
		</div>
	);
};
