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
import { useMutation, useQuery } from "@tanstack/react-query";
import { request } from "../libs/api";

import useStore from "../libs/store";

const url = new URL(window.location.href);

export default function Sidebar() {
  const isLogin = useStore((state) => state.isLogin);
  const assistantId = url.searchParams.get("assistant-id");
  if (!assistantId) {
    window.location.href = "/assistants.html";
  }
  const threadId = url.searchParams.get("thread-id");

  const setThreadId = (id) => {
    url.searchParams.set("thread-id", id);
    window.location.href = url;
  };

  const listThreads = useQuery({
      queryKey: ["list-threads", assistantId],
      queryFn: async () => {
          const res = await request(
              `/api/v1/assistants/${assistantId}/threads`
          );
          const data = await res.json().data;
          data.map((item) => {
              item["value"] =
                  item["name"] +
                  " - " +
                  item["description"] +
                  " - " +
                  item["id"];
          });
          return data;
      },
      enabled: isLogin,
  });

  const getAssistant = useQuery({
      queryKey: ["get-assistant", assistantId],
      queryFn: async () => {
          const res = await request(`/api/v1/assistants/${assistantId}`);
          const data = await res.json().data;
          return data;
      },
      enabled: isLogin,
  });

  const createThread = useMutation({
    mutationFn: async (data) => {
      const res = await request(
          `/api/v1/assistants/${assistantId}/threads`,
          (method = "POST"),
          (body = JSON.stringify(data))
      );
      const response = await res.json();
      return response.data;
    },
    onSuccess: (data) => {
      setThreadId(data.id);
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
        <ScrollArea>
          <LoadingOverlay visible={listThreads.isLoading} />
          <Stack align="stretch" justify="start" gap="sm">
            {listThreads.data &&
              listThreads.data.map((item) => (
                <Button
                  radius="md"
                  w="100%"
                  key={item.id}
                  variant={threadId == item.id ? "filled" : "subtle"}
                  onClick={() => {
                    setThreadId(item.id);
                  }}
                >
                  {item.name}
                </Button>
              ))}
          </Stack>
        </ScrollArea>
      </Flex>
    </>
  );
}
