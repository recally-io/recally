import { Avatar, Flex, Paper, ScrollArea, Stack, Text } from "@mantine/core";

import { Icon } from "@iconify/react/dist/iconify.js";
import { useMutation } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import { post, queryClient } from "../libs/api";
import useStore from "../libs/store";
import { CopyBtn } from "./copy-button";
import { MarkdownRenderer } from "./markdown-renderer";

export function ThreadChatWindows({ settingsForm }) {
  const threadId = useStore((state) => state.threadId);
  const assistantId = useStore((state) => state.assistantId);

  const isDarkMode = useStore((state) => state.isDarkMode);
  const messageList = useStore((state) => state.threadMessageList);
  const [isTitleGenerated, setIsTitleGenerated] = useStore((state) => [
    state.threadIsTitleGenerated,
    state.setThreadIsTitleGenerated,
  ]);
  const chatArea = useRef(null);

  useEffect(() => {
    chatArea.current.scrollTo({
      top: chatArea.current.scrollHeight,
      behavior: "smooth",
    });

    const generate = async () => {
      await generateTitle.mutateAsync();
      if (generateTitle.isSuccess) {
        setIsTitleGenerated(true);
      }
    };

    if (messageList.length >= 4 && !isTitleGenerated) {
      console.log("Generate title");
      generate();
    }
  }, [messageList]);

  const generateTitle = useMutation({
    mutationFn: async () => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/generate-title`,
        null,
        {},
      );
      return res.data;
    },
    onSuccess: (data) => {
      settingsForm.setFieldValue("name", data.name);
      queryClient.invalidateQueries(["get-thread", threadId]);
    },
  });

  const messagePaper = (message) => {
    const isSender = message.role === "user";
    const bgColor = isSender ? "blue.2" : "green.1";
    return (
      <Flex
        justify={isSender ? "flex-end" : "flex-start"}
        align="flex-start"
        direction="row"
        gap="2"
        key={message.id}
      >
        {!isSender && (
          <Avatar size="sm" radius="lg" color="cyan" variant="filled">
            <Icon icon="tabler:robot" />
          </Avatar>
        )}
        <Flex align={isSender ? "flex-end" : "flex-start"} direction="column">
          <Text size="lg" variant="gradient">
            {isSender ? "You" : message.model}
          </Text>

          <Paper
            shadow="sm"
            px="xs"
            radius="lg"
            withBorder
            maw={{ base: "85vw", lg: "60vw" }}
            bg={isDarkMode ? "" : bgColor}
          >
            <MarkdownRenderer content={message.text} />
          </Paper>

          {CopyBtn({ data: message.text })}
        </Flex>
        {isSender && (
          <Avatar size="sm" radius="lg" color="cyan" variant="filled">
            <Icon icon="tabler:mood-crazy-happy" />
          </Avatar>
        )}
      </Flex>
    );
  };

  return (
    <>
      <ScrollArea
        viewportRef={chatArea}
        type="scroll"
        offsetScrollbars
        scrollbarSize="4"
        scrollbars="y"
        py="xs"
      >
        <Stack gap="md" align="stretch" justify="flex-start">
          {Array.isArray(messageList) &&
            messageList.map((item) => {
              return messagePaper(item);
            })}
        </Stack>
      </ScrollArea>
    </>
  );
}
