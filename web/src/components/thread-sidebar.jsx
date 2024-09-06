import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
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
import { Link, useNavigate, useParams } from "react-router-dom";
import { useQueryContext } from "../libs/query-context";

import useStore from "../libs/store";
import { ThreadAddButton } from "./thread-add-button";

export default function ThreadSidebar() {
  const { listThreads, deleteThread } = useQueryContext();

  const navigate = useNavigate();
  const params = useParams();
  const threadId = params.threadId;
  const assistantId = params.assistantId;

  const isDarkMode = useStore((state) => state.isDarkMode);
  const toggleMobileSidebar = useStore((state) => state.toggleMobileSidebar);
  const setMessageList = useStore((state) => state.setThreadMessageList);

  const navigateToThread = (threadId) => {
    toggleMobileSidebar();
    // setMessageList([]);
    navigate(`/assistants/${assistantId}/threads/${threadId}`, {
      replace: false,
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
        m="0"
        radius="md"
        bg={isDarkMode ? "dark.6" : "gray.2"}
      >
        <Stack align="stretch" justify="start" gap="md">
          <Flex justify="center" align="center" gap="md">
            <Button variant="outline" radius="lg" size="sm">
              <Link to={`/assistants/${assistantId}`} replace>
                <Text
                  size="sm"
                  variant="gradient"
                  gradient={{ from: "pink", to: "yellow" }}
                  sx={{ textDecoration: "underline" }}
                >
                  Assistant
                </Text>
              </Link>
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
                navigateToThread(filteredItems[0].id);
              }
            }}
          />
        </Stack>
        <Divider />
        <ScrollArea scrollbarSize="4" scrollbars="y">
          <LoadingOverlay visible={listThreads.isLoading} />
          <Stack align="start" justify="start" gap="sm">
            {listThreads.data &&
              listThreads.data.map((item) => (
                <Group
                  justify="space-between"
                  grow
                  preventGrowOverflow={false}
                  gap="xs"
                  key={item.id}
                  wrap="nowrap"
                >
                  <Button
                    radius="md"
                    color={threadId == item.id ? "accent" : "default"}
                    variant={threadId == item.id ? "filled" : "subtle"}
                    onClick={() => {
                      navigateToThread(item.id);
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
                      size="xs"
                      onClick={async () => {
                        await deleteThread.mutateAsync(item.id);
                        navigate(`/assistants/${assistantId}/threads`);
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
