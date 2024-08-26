import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Anchor,
  Autocomplete,
  Button,
  Divider,
  Flex,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
} from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { del, get, queryClient } from "../libs/api";

import { useEffect } from "react";
import useStore from "../libs/store";
import { ThreadAddButton } from "./thread-add-button";

export default function Sidebar() {
  const isLogin = useStore((state) => state.isLogin);
  const assistant = useStore((state) => state.assistant);
  const thread = useStore((state) => state.thread);
  const setThread = useStore((state) => state.setThread);

  const [threadId, setThreadId] = useStore((state) => [
    state.threadId,
    state.setThreadId,
  ]);

  const toggleMobileSidebar = useStore((state) => state.toggleMobileSidebar);
  const setMessageList = useStore((state) => state.setThreadMessageList);

  useEffect(() => {
    if (threadId) {
      const url = new URL(window.location.href);
      url.searchParams.set("thread-id", threadId);
      window.history.pushState({}, "", url);
    }
  }, [threadId]);

  const listThreads = useQuery({
    queryKey: ["list-threads", assistant.id],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistant.id}/threads`);
      const data = res.data;
      data.map((item) => {
        item["value"] =
          item["name"] + " - " + item["description"] + " - " + item["id"];
      });
      return data;
    },
    enabled: isLogin && !!assistant,
  });

  const deleteThread = useMutation({
    mutationFn: async () => {
      await del(`/api/v1/assistants/${assistant.id}/threads/${thread.id}`);
      console.log("delete thread success");
    },
    onSuccess: () => {
      console.log("onSuccess: delete thread success");
      queryClient.invalidateQueries({
        queryKey: ["list-threads", assistant.id],
      });
      setThread(null);
      setMessageList([]);
      const url = new URL(window.location.href);
      url.searchParams.delete("thread-id");
      window.history.pushState({}, "", url);
    },
  });

  return (
    <>
      <Flex
        direction="column"
        align="stretch"
        justify="start"
        gap="md"
        p="sm"
        h="100%"
        radius="md"
      >
        <Stack align="stretch" justify="start" gap="md">
          <Flex justify="center" align="center" gap="md">
            <Button variant="outline" radius="lg" size="sm">
              <Anchor
                href={`/assistants.html?id=${assistant.id}`}
                variant="gradient"
                gradient={{ from: "pink", to: "yellow" }}
                underline="always"
              >
                <Text size="sm">Assistant</Text>
              </Anchor>
            </Button>

            <ThreadAddButton />
            <Tooltip label="Toggle Sidebar" hiddenFrom="sm">
              <ActionIcon
                onClick={toggleMobileSidebar}
                variant="subtle"
                radius="lg"
                size="lg"
                hiddenFrom="sm"
              >
                <Icon icon="tabler:menu-3" />
              </ActionIcon>
            </Tooltip>
          </Flex>

          <Autocomplete
            placeholder="Search Threads ... "
            limit={10}
            leftSection={<Icon icon="tabler:search" />}
            radius="lg"
            data={listThreads.data}
            onChange={(item) => {
              var filteredItems = listThreads.data.filter(
                (i) => i.value == item,
              );
              if (filteredItems.length > 0) {
                setThreadId(filteredItems[0].id);
                toggleMobileSidebar();
              }
            }}
          />
        </Stack>
        <Divider />
        <ScrollArea scrollbarSize="4" scrollbars="y">
          <LoadingOverlay visible={listThreads.isLoading} />
          <Stack align="stretch" justify="start" gap="sm">
            {listThreads.data &&
              listThreads.data.map((item) => (
                <Flex
                  key={item.id}
                  align="center"
                  w="95%"
                  justify="space-between"
                >
                  <Button
                    radius="md"
                    color={thread?.id == item.id ? "accent" : "default"}
                    variant={thread?.id == item.id ? "filled" : "subtle"}
                    onClick={() => {
                      setThreadId(item.id);
                      toggleMobileSidebar();
                    }}
                    styles={{
                      inner: {
                        display: "flex",
                        flexDirection: "column",
                        alignItems: "flex-start",
                      },
                    }}
                  >
                    <Text
                      size="sm"
                      lineClamp={2}
                      style={{
                        whiteSpace: "normal",
                        textAlign: "left",
                      }}
                    >
                      {item.name}
                    </Text>
                  </Button>
                  {thread?.id == item.id && (
                    <ActionIcon
                      variant="subtle"
                      color="danger"
                      onClick={async () => {
                        await deleteThread.mutateAsync(item.id);
                      }}
                    >
                      <Icon icon="tabler:trash" />
                    </ActionIcon>
                  )}
                </Flex>
              ))}
          </Stack>
        </ScrollArea>
      </Flex>
    </>
  );
}
