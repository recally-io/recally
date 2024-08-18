import { Icon } from "@iconify/react/dist/iconify.js";
import {
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
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import useStore from "../libs/store";
import { AssistantsApi } from "../sdk/index";

const api = new AssistantsApi();

export default function Sidebar() {
  const queryClient = useQueryClient();
  const [threadId, setThreadId] = useStore((state) => [
    state.activateThreadId,
    state.setActivateThreadId,
  ]);

  const [assistantId, setAssistantId] = useStore((state) => [
    state.activateAssistantId,
    state.setActivateAssistantId,
  ]);

  useEffect(() => {
    let url = new URL(window.location.href);
    let params = new URLSearchParams(url.search);
    let assistantId = params.get("assistant-id");
    console.log(`useEffect: assistantId ${assistantId}`);
    if (!assistantId) {
      setAssistantId("");
      window.location.href = "/assistants.html";
    }
    setAssistantId(assistantId);
    const threadId = params.get("thread-id");
    console.log(`useEffect: threadId ${threadId}`);
    if (threadId) {
      setThreadId(threadId);
    } else {
      setThreadId("");
    }
    console.log(`useEffect: assistantId ${assistantId} , threadId ${threadId}`);
  }, []);

  const listThreads = useQuery({
    queryKey: ["list-threads", assistantId],
    queryFn: async () => {
      console.log(`listThreads: assistantId ${assistantId}`);
      const response = await api.assistantsAssistantIdThreadsGet({
        assistantId: assistantId,
      });
      const data = response.data;
      data.map((item) => {
        console.log(`item`, JSON.stringify(item));
        item["value"] =
          item["name"] + " - " + item["description"] + " - " + item["id"];
      });
      return data;
    },
  });

  const getAssistant = useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const response = await api.assistantsAssistantIdGet({
        assistantId: assistantId,
      });
      return response.data;
    },
  });

  const createThread = useMutation({
    mutationFn: async (data) => {
      const response = await api.assistantsAssistantIdThreadsPost({
        assistantId: assistantId,
        thread: data,
      });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["list-threads"]);
    },
  });

  const theme = useMantineTheme();

  const addNewThread = async () => {
    const newThread = await createThread.mutateAsync({
      name: "Thread name",
      description: "Thread description",
      systemPrompt: getAssistant.data.systemPrompt,
      model: getAssistant.data.model,
    });
    setThreadId(newThread.id);
  };

  return (
    <>
      <Flex
        direction="column"
        align="stretch"
        justify="start"
        gap="md"
        padding="md"
        h="100%"
      >
        <Stack align="stretch" justify="start" gap="md">
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
          <Autocomplete
            placeholder="Search Threads ... "
            limit={10}
            leftSection={<Icon icon="tabler:search" />}
            radius="lg"
            data={listThreads.data}
            onChange={(item) => {
              var filteredItems = conversations.filter((i) => i.value == item);
              if (filteredItems.length > 0) {
                setThreadId(filteredItems[0].id);
              }
            }}
          />
        </Stack>
        <Divider />
        <ScrollArea>
          <LoadingOverlay visible={listThreads.isLoading} />
          <Stack align="stretch" justify="start" gap="sm">
            {listThreads.data &&
              listThreads.data.map((item) => (
                <Button
                  radius="md"
                  w="100%"
                  variant={threadId == item.id ? "filled" : "subtle"}
                  onClick={() => {
                    console.log(
                      `assistantId ${assistantId}, threadId ${item.id}`,
                    );
                    setThreadId(item.id);
                  }}
                >
                  {item.value}
                </Button>
              ))}
          </Stack>
        </ScrollArea>
      </Flex>
    </>
  );
}
