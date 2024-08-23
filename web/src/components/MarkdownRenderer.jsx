import { CodeHighlight } from "@mantine/code-highlight";
import "@mantine/code-highlight/styles.css";
import { Badge, Box, Group } from "@mantine/core";
import React from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CopyBtn } from "./CopyButton";

export function MarkdownRenderer({ content }) {
    return (
        <Markdown
            children={content}
            remarkPlugins={[remarkGfm]}
            components={{
                code({ node, inline, className, children, ...props }) {
                    const match = /language-(\w+)/.exec(className || "");
                    const code = String(children).replace(/\n$/, "");
                    return !inline && match ? (
                        <Box pos="relative">
                            <CodeHighlight
                                code={code}
                                language={match[1]}
                                withCopyButton={false}
                                {...props}
                            />
                            <Group
                                pos="absolute"
                                right="10px"
                                top="10px"
                                size="xs"
                            >
                                <Badge>{match[1]}</Badge>
                                <CopyBtn data={code} />
                            </Group>
                        </Box>
                    ) : (
                        <code className={className} {...props}>
                            {children}
                        </code>
                    );
                },
            }}
        />
    );
}
