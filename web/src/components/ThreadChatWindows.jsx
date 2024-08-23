import { Icon } from "@iconify/react/dist/iconify.js";
import { Avatar, Flex, Paper, ScrollArea, Stack, Text } from "@mantine/core";

import { useMutation } from "@tanstack/react-query";
import { useEffect, useRef } from "react";
import avatarImgUrl from "../assets/avatar-1.png";
import { post, queryClient } from "../libs/api";
import useStore from "../libs/store";
import { CopyBtn } from "./CopyButton";
import { MarkdownRenderer } from "./MarkdownRenderer";

const url = new URL(window.location.href);

export function ThreadChatWindows({ settingsForm }) {
  const assistantId = url.searchParams.get("assistant-id");
  let threadId = url.searchParams.get("thread-id");

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

  const messageS = (message) => {
    return (
      <Flex
        justify="flex-end"
        align="flex-start"
        direction="row"
        gap="sm"
        key={message.id}
      >
        <Flex align="flex-end" direction="column">
          <Text size="lg" variant="gradient">
            You
          </Text>
          <Paper
            shadow="sm"
            px="sm"
            w="100%"
            radius="lg"
            bg={isDarkMode ? "" : "blue.2"}
          >
            <ScrollArea type="auto" scrollbars="x">
              <MarkdownRenderer content={message.text} />
            </ScrollArea>
          </Paper>
          {CopyBtn({ data: message.text })}
        </Flex>

        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
      </Flex>
    );
  };

  const messageR = (message) => {
    return (
      <Flex justify="flex-start" direction="row" gap="sm" key={message.id}>
        <Avatar size="sm" radius="lg" color="cyan" variant="filled">
          <Icon icon="tabler:robot" />
        </Avatar>
        <Flex align="flex-start" direction="column" maw="90%">
          <Text size="lg" variant="gradient">
            {message.model}
          </Text>
          <Paper
            shadow="sm"
            px="sm"
            w="100%"
            radius="lg"
            bg={isDarkMode ? "" : "green.1"}
          >
            <ScrollArea type="auto" scrollbars="x">
              <MarkdownRenderer content={message.text} />
            </ScrollArea>
          </Paper>
          {CopyBtn({ data: message.text })}
        </Flex>
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
      >
        <Stack spacing="md" py="lg">
          {messageList.map((item) => {
            if (item.role === "user") {
              return messageS(item);
            } else {
              return messageR(item);
            }
          })}
        </Stack>
      </ScrollArea>
    </>
  );
}
