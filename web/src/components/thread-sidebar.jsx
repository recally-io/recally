import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Autocomplete,
  Button,
  Divider,
  Flex,
  LoadingOverlay,
  ScrollArea,
  Stack,
} from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { del, get, queryClient } from "../libs/api";

import { useEffect } from "react";
import useStore from "../libs/store";
import { ThreadAddButton } from "./thread-add-button";

export default function Sidebar() {
  const isLogin = useStore((state) => state.isLogin);
  const [threadId, setThreadId] = useStore((state) => [
    state.threadId,
    state.setThreadId,
  ]);
  const [isSidebarOpen, toggleSidebar] = useStore((state) => [
    state.isSidebarOpen,
    state.toggleSidebar,
  ]);
  const assistantId = useStore((state) => state.assistantId);
  const setMessageList = useStore((state) => state.setThreadMessageList);

  useEffect(() => {
    if (threadId) {
      const url = new URL(window.location.href);
      url.searchParams.set("thread-id", threadId);
      window.history.pushState({}, "", url);
    }
  }, [threadId]);

  const listThreads = useQuery({
    queryKey: ["list-threads", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}/threads`);
      const data = res.data;
      data.map((item) => {
        item["value"] =
          item["name"] + " - " + item["description"] + " - " + item["id"];
      });
      return data;
    },
    enabled: isLogin,
  });

  const deleteThread = useMutation({
    mutationFn: async (threadId) => {
      await del(`/api/v1/assistants/${assistantId}/threads/${threadId}`);
      console.log("delete thread success");
    },
    onSuccess: () => {
      console.log("onSuccess: delete thread success");
      queryClient.invalidateQueries({
        queryKey: ["list-threads", assistantId],
      });
      setThreadId(null);
      // toggleSidebar();
      // reload the page
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
          <Flex justify="space-evenly" align="center">
            <ThreadAddButton />
            <Button
              onClick={toggleSidebar}
              variant="subtle"
              size="lg"
              hiddenFrom="sm"
            >
              {isSidebarOpen ? (
                <Icon icon="tabler:chevron-right" />
              ) : (
                <Icon icon="tabler:chevron-left" />
              )}
            </Button>
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
                  w={{ xs: "95%", sm: "90%" }}
                  justify="space-between"
                >
                  <Button
                    radius="md"
                    variant={threadId == item.id ? "filled" : "subtle"}
                    onClick={() => {
                      setThreadId(item.id);
                    }}
                  >
                    {item.name}
                  </Button>
                  {threadId == item.id && (
                    <ActionIcon
                      variant="subtle"
                      color="red"
                      onClick={async () => {
                        await deleteThread.mutateAsync(item.id);
                      }}
                    >
                      <Icon icon="tabler:trash" width={18} height={18} />
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
