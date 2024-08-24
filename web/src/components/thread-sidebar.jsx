import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Autocomplete,
  Button,
  Divider,
  Flex,
  LoadingOverlay,
  ScrollArea,
  Space,
  Stack,
  useMantineTheme,
} from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { del, get, post, queryClient } from "../libs/api";

import { useEffect } from "react";
import useStore from "../libs/store";

export default function Sidebar() {
  const isLogin = useStore((state) => state.isLogin);
  const [threadId, setThreadId] = useStore((state) => [
    state.threadId,
    state.setThreadId,
  ]);
  const toggleSidebar = useStore((state) => state.toggleSidebar);
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

  const getAssistant = useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}`);
      return res.data;
    },
    enabled: isLogin,
  });

  const createThread = useMutation({
    mutationFn: async (data) => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads`,
        null,
        data,
      );
      return res.data;
    },
    onSuccess: (data) => {
      setThreadId(data.id);
    },
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
      toggleSidebar();
      // reload the page
      setMessageList([]);
      const url = new URL(window.location.href);
      url.searchParams.delete("thread-id");
      window.history.pushState({}, "", url);
    },
  });

  const theme = useMantineTheme();

  const addNewThread = async () => {
    await createThread.mutateAsync({
      name: "Thread name",
      description: "Thread description",
      systemPrompt: getAssistant.data.systemPrompt,
      model: getAssistant.data.model,
    });
  };

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
            <Button
              variant="subtle"
              radius="lg"
              color={theme.primaryColor}
              onClick={addNewThread}
            >
              <Icon icon="tabler:message-circle" width={18} height={18} />
              <Space w="xs" />
              <span>New Thread</span>
            </Button>
            <Button onClick={toggleSidebar} variant="transparent" size="lg">
              <Icon icon="tabler:layout-sidebar" />
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
