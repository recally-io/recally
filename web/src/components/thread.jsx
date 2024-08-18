import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Avatar,
  Button,
  Container,
  FileButton,
  Flex,
  Group,
  Menu,
  Modal,
  Paper,
  ScrollArea,
  Slider,
  Stack,
  Text,
  TextInput,
  Tooltip,
  useComputedColorScheme,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import Markdown from "react-markdown";
import avatarImgUrl from "../assets/avatar-1.png";
import useStore from "../libs/store";
import { AssistantsApi } from "../sdk/index";

const api = new AssistantsApi();

export default function ChatWindowsComponent() {
  const queryClient = useQueryClient();
  const [settingsOpened, { open: openSettings, close: closeSettings }] =
    useDisclosure(false);
  const colorScheme = useComputedColorScheme("light");
  const settingsForm = useForm({
    initialValues: {
      temperature: 0.7,
      maxToken: 4096,
      model: "gpt-4o",
    },
  });

  const assistantId = useStore((state) => state.activateAssistantId);
  const threadId = useStore((state) => state.activateThreadId);
  const [newText, setNewText] = useState("");

  const getThread = useQuery({
    queryKey: ["get-thread", threadId],
    queryFn: async () => {
      console.log(
        `getThread: threadId ${threadId}, assistantId ${assistantId}`,
      );
      const response =
        await api.AssistantsAssistantIdThreadsThreadIdMessagesGetRequest({
          assistantId: assistantId,
          threadId: threadId,
        });
      console.log(JSON.stringify(response));
      return response.data;
    },
  });

  const listMessages = useQuery({
    queryKey: ["list-messages", threadId],
    queryFn: async () => {
      console.log(
        `listMessages: threadId ${threadId}, assistantId ${assistantId}`,
      );
      const response =
        await api.assistantsAssistantIdThreadsThreadIdMessagesGet({
          assistantId: assistantId,
          threadId: threadId,
        });
      return response.data;
    },
  });

  const createMessage = useMutation({
    mutationFn: async () => {
      const response =
        await api.assistantsAssistantIdThreadsThreadIdMessagesPost({
          assistantId: assistantId,
          threadId: threadId,
          message: {
            role: "user",
            text: newText,
            model: "gpt-4o",
          },
        });
      return response.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["list-messages"]);
    },
  });

  const messageS = (role, text) => {
    return (
      <Flex justify="flex-end" align="flex-start" direction="row" gap="sm">
        <Paper
          shadow="sm"
          p="md"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "blue.2"}
        >
          <Markdown>{text}</Markdown>
        </Paper>
        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
      </Flex>
    );
  };

  const messageR = (role, text) => {
    return (
      <Flex justify="flex-start" direction="row" gap="sm">
        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
        <Paper
          shadow="sm"
          p="md"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "green.2"}
        >
          <Markdown>{text}</Markdown>
        </Paper>
      </Flex>
    );
  };

  const menu = () => {
    return (
      <Menu shadow="md" position="top" withArrow>
        <Menu.Target>
          <Button size="compact-lg" variant="subtle" radius="lg">
            <Icon icon="tabler:plus"></Icon>
          </Button>
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Item>
            <Tooltip label="thread settings">
              <Button
                variant="transparent"
                size="sm"
                onClick={openSettings}
                leftSection={<Icon icon="tabler:settings"></Icon>}
              >
                Settings
              </Button>
            </Tooltip>
          </Menu.Item>
          <Menu.Item>
            <FileButton
              size="sm"
              variant="transparent"
              multiple
              leftSection={<Icon icon="tabler:upload"></Icon>}
            >
              {(props) => <Button {...props}>Upload image</Button>}
            </FileButton>
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    );
  };

  return (
    <>
      <Container size="xl">
        <Flex direction="column" justify="space-between" h="89vh">
          <ScrollArea
            style={{
              flex: 1,
            }}
          >
            <Stack spacing="md" py="lg">
              {listMessages.data &&
                listMessages.data.map((item) => {
                  if (item.role === "user") {
                    return messageS(item.role, item.text);
                  } else {
                    return messageR(item.role, item.text);
                  }
                })}
            </Stack>
          </ScrollArea>
          <Container
            w="100%"
            style={{
              position: "sticky",
              bottom: 0,
            }}
          >
            <TextInput
              placeholder="Send a message"
              variant="filled"
              radius="lg"
              leftSection={menu()}
              leftSectionWidth={42}
              disabled={createMessage.isLoading}
              rightSection={
                <ActionIcon
                  variant="transparent"
                  aria-label="Settings"
                  disabled={createMessage.isLoading}
                  onClick={async () => {
                    await createMessage.mutateAsync();
                    setNewText("");
                  }}
                >
                  <Icon icon="tabler:send"></Icon>
                </ActionIcon>
              }
              value={newText}
              onChange={(e) => setNewText(e.currentTarget.value)}
            ></TextInput>
          </Container>
        </Flex>

        {/* settings modal */}
        <Modal
          opened={settingsOpened}
          onClose={closeSettings}
          title="Advance Settings"
        >
          <form
            onSubmit={settingsForm.onSubmit((values) => console.log(values))}
            mode=""
          >
            <Stack spacing="md">
              <Stack spacing="xs">
                <Text size="sm">Temperature</Text>
                <Slider
                  min={0}
                  max={1}
                  step={0.1}
                  key={settingsForm.key("temperature")}
                  {...settingsForm.getInputProps("temperature")}
                  labelAlwaysOn
                />
              </Stack>
              <Stack spacing="xs">
                <Text size="sm">Max Tokens</Text>
                <Slider
                  min={0}
                  max={4096}
                  step={1}
                  key={settingsForm.key("maxToken")}
                  {...settingsForm.getInputProps("maxToken")}
                  labelAlwaysOn
                />
              </Stack>
            </Stack>
            <Group justify="flex-end" mt="md">
              <Button type="submit">Submit</Button>
              <Button type="reset">Reset</Button>
            </Group>
          </form>
        </Modal>
      </Container>
    </>
  );
}
