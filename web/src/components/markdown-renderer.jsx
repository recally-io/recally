import { CodeHighlight } from "@mantine/code-highlight";
import "@mantine/code-highlight/styles.css";
import { Badge, Box, Code, Group } from "@mantine/core";
import React from "react";
import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { CopyBtn } from "./copy-button";

export function MarkdownRenderer({ content }) {
  return (
    <Markdown
      children={content}
      remarkPlugins={[remarkGfm]}
      components={{
        code({ children, className, node, ...rest }) {
          // Extract language from className if present
          const match = /language-(\w+)/.exec(className || "");
          // Remove trailing newline from code content
          const code = String(children).replace(/\n$/, "");
          console.log(className, match, code);
          // If no className or language match, render as inline code
          if (!className) {
            if (code.includes("\n")) {
              return (
                <CodeHighlight code={code} {...rest} withCopyButton={false} />
              );
            }
            return <Code {...rest}>{code}</Code>;
          }

          // If language match is found
          if (match) {
            return (
              <Box pos="relative">
                {/* Render code with syntax highlighting */}
                <CodeHighlight
                  code={code}
                  language={match[1]}
                  withCopyButton={false}
                  {...rest}
                />
                {/* Add language badge and copy button */}
                <Group pos="absolute" right="10px" top="10px" size="xs">
                  <Badge>{match[1]}</Badge>
                  <CopyBtn data={code} />
                </Group>
              </Box>
            );
          }

          // If no language match, render use default syntax highlighting
          return <CodeHighlight code={code} {...rest} withCopyButton={false} />;
        },
      }}
    />
  );
}
