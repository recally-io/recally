import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Anchor,
  Autocomplete,
  Button,
  Divider,
  Flex,
  Group,
  LoadingOverlay,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
} from "@mantine/core";
import { useQueryContext } from "../libs/query-context";

import useStore from "../libs/store";
import { ThreadAddButton } from "./thread-add-button";

export default function Sidebar() {
  const { listThreads, deleteThread, updateThreadId } = useQueryContext();

  const isDarkMode = useStore((state) => state.isDarkMode);
  const assistantId = useStore((state) => state.assistantId);
  const threadId = useStore((state) => state.threadId);
  const toggleMobileSidebar = useStore((state) => state.toggleMobileSidebar);

  return (
    <>
      <Flex
        direction="column"
        align="stretch"
        justify="start"
        gap="md"
        p="sm"
        h="100%"
        m="0"
        radius="md"
        bg={isDarkMode ? "dark.6" : "gray.2"}
      >
        <Stack align="stretch" justify="start" gap="md">
          <Flex justify="center" align="center" gap="md">
            <Button variant="outline" radius="lg" size="sm">
              <Anchor
                href={`/assistants.html?id=${assistantId}`}
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
                toggleMobileSidebar();
                updateThreadId(filteredItems[0].id);
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
                <Group
                  justify="space-between"
                  grow
                  preventGrowOverflow={false}
                  gap="2"
                  key={item.id}
                >
                  <Button
                    radius="md"
                    color={threadId == item.id ? "accent" : "default"}
                    justify="space-between"
                    variant={threadId == item.id ? "filled" : "subtle"}
                    onClick={() => {
                      toggleMobileSidebar();
                      updateThreadId(item.id);
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
                  {threadId == item.id ? (
                    <ActionIcon
                      variant="subtle"
                      color="danger"
                      onClick={async () => {
                        await deleteThread.mutateAsync(item.id);
                      }}
                    >
                      <Icon icon="tabler:trash" />
                    </ActionIcon>
                  ) : null}
                </Group>
              ))}
          </Stack>
        </ScrollArea>
      </Flex>
    </>
  );
}
