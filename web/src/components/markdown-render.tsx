import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import rehypeRaw from "rehype-raw";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import {
	oneDark,
	oneLight,
} from "react-syntax-highlighter/dist/esm/styles/prism";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { useTheme } from "next-themes";

interface MarkdownRendererProps {
	content: string;
	className?: string;
}

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
	content,
	className,
}) => {
	const { theme } = useTheme();

	// Custom components for markdown elements
	const components = {
		// Heading components
		h1: ({ node, ...props }) => (
			<h1
				{...props}
				className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4"
			/>
		),
		h2: ({ node, ...props }) => (
			<h2
				{...props}
				className="scroll-m-20 text-3xl font-semibold tracking-tight mb-3"
			/>
		),
		h3: ({ node, ...props }) => (
			<h3
				{...props}
				className="scroll-m-20 text-2xl font-semibold tracking-tight mb-2"
			/>
		),

		// Paragraph and text elements
		p: ({ node, ...props }) => (
			<p {...props} className="leading-7 [&:not(:first-child)]:mt-6" />
		),
		a: ({ node, ...props }) => (
			<a
				{...props}
				className="font-medium text-primary underline underline-offset-4 hover:text-primary/80"
				target="_blank"
				rel="noopener noreferrer"
			/>
		),

		// List elements
		ul: ({ node, ...props }) => (
			<ul {...props} className="my-6 ml-6 list-disc [&>li]:mt-2" />
		),
		ol: ({ node, ...props }) => (
			<ol {...props} className="my-6 ml-6 list-decimal [&>li]:mt-2" />
		),

		// Code blocks
		code: ({ node, inline, className, children, ...props }) => {
			const match = /language-(\w+)/.exec(className || "");
			return !inline && match ? (
				<Card className="my-4">
					<CardContent className="p-0">
						<SyntaxHighlighter
							style={theme === "dark" ? oneDark : oneLight}
							language={match[1]}
							PreTag="div"
							customStyle={{
								margin: 0,
							}}
							{...props}
						>
							{String(children).replace(/\n$/, "")}
						</SyntaxHighlighter>
					</CardContent>
				</Card>
			) : (
				<code
					className="relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold"
					{...props}
				>
					{children}
				</code>
			);
		},

		// Blockquote
		blockquote: ({ node, ...props }) => (
			<blockquote
				{...props}
				className="mt-6 border-l-2 border-primary pl-6 italic"
			/>
		),

		// Table elements
		table: ({ node, ...props }) => (
			<div className="my-6 w-full overflow-y-auto">
				<table {...props} className="w-full" />
			</div>
		),
		th: ({ node, ...props }) => (
			<th
				{...props}
				className="border px-4 py-2 text-left font-bold [&[align=center]]:text-center [&[align=right]]:text-right"
			/>
		),
		td: ({ node, ...props }) => (
			<td
				{...props}
				className="border px-4 py-2 text-left [&[align=center]]:text-center [&[align=right]]:text-right"
			/>
		),
	};

	return (
		<div className={cn("markdown-content", className)}>
			<ReactMarkdown
				components={components}
				remarkPlugins={[remarkGfm]}
				rehypePlugins={[rehypeRaw]}
			>
				{content}
			</ReactMarkdown>
		</div>
	);
};

export default MarkdownRenderer;
