import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Autocomplete,
  Button,
  Divider,
  Flex,
  ScrollArea,
  Space,
  Stack,
  useMantineTheme,
} from "@mantine/core";
import { useState } from "react";

export default function Sidebar() {
  const theme = useMantineTheme();
  const data = [
    { id: "1", value: "Thread 1" },
    { id: "2", value: "Thread 2" },
    { id: "3", value: "Thread 3" },
    { id: "4", value: "Thread 4" },
    { id: "5", value: "Thread 5" },
    { id: "6", value: "Thread 6" },
    { id: "7", value: "Thread 7" },
    { id: "8", value: "Thread 8" },
    { id: "9", value: "Thread 9" },
    { id: "10", value: "Thread 10" },
  ];
  const [conversations, setConversations] = useState(data);
  const [activateThreadId, setActivateThreadId] = useState();

  const addNewThread = () => {
    console.log("Add new Thread", conversations);
    const newThreadId = `${conversations.length + 1}`;
    setConversations([
      ...conversations,
      { id: newThreadId, value: `Thread ${newThreadId}` },
    ]);
    setActivateThreadId(newThreadId);
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
            data={conversations}
            onChange={(item) => {
              var filteredItems = conversations.filter((i) => i.value == item);
              if (filteredItems.length > 0) {
                setActivateThreadId(filteredItems[0].id);
              }
            }}
          />
        </Stack>
        <Divider />
        <ScrollArea>
          <Stack align="stretch" justify="start" gap="sm">
            {conversations.map((item) => (
              <Button
                onClick={() => setActivateThreadId(item.id)}
                key={item.id}
                variant={activateThreadId == item.id ? "filled" : "subtle"}
                radius="md"
              >
                Thread tilte {item.value}
              </Button>
            ))}
          </Stack>
        </ScrollArea>
      </Flex>
    </>
  );
}
