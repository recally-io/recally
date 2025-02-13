interface ArticleHeaderProps {
	title: string;
	url: string;
	publishedAt?: string;
}

export const ArticleHeader: React.FC<ArticleHeaderProps> = ({
	title,
	url,
	publishedAt,
}) => {
	const domain = url.startsWith("http")
		? new URL(url).hostname
		: window.location.hostname;

	return (
		<div className="space-y-4">
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

				{publishedAt && (
					<>
						<span className="text-muted-foreground/40">•</span>
						<span>
							{new Date(publishedAt).toLocaleDateString("en-US", {
								year: "numeric",
								month: "long",
								day: "numeric",
							})}
						</span>
					</>
				)}
			</div>
		</div>
	);
};
