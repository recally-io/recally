import { Badge, Box, Group } from "@mantine/core";

import React from "react";
import Markdown from "react-markdown";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { oneDark } from "react-syntax-highlighter/dist/esm/styles/prism";
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
          return match ? (
            <Box pos="relative">
              <SyntaxHighlighter
                {...props}
                PreTag="div"
                children={code}
                language={match[1]}
                style={oneDark}
              />
              <Group pos="absolute" right="10px" top="10px" size="xs">
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
