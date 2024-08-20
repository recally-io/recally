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
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState, useRef } from "react";
import Markdown from "react-markdown";
import avatarImgUrl from "../assets/avatar-1.png";
import useStore from "../libs/store";
import { AssistantsApi } from "../sdk/index";

const url = new URL(window.location.href);
const api = new AssistantsApi();

export default function ChatWindowsComponent() {
  const isLogin = useStore((state) => state.isLogin);
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

  const assistantId = url.searchParams.get("assistant-id");
  const threadId = url.searchParams.get("thread-id");
  const [newText, setNewText] = useState("");

  const [messageList, setMessageList] = useState([]);
  const chatArea = useRef(null);

  const listMessages = useQuery({
    queryKey: ["list-messages", threadId],
    queryFn: async () => {
      const response =
        await api.assistantsAssistantIdThreadsThreadIdMessagesGet({
          assistantId: assistantId,
          threadId: threadId,
        });
      return response.data || [];
    },
    enabled: isLogin && !!threadId && !!assistantId,
  });

  useEffect(() => {
    if (listMessages.data) {
      setMessageList(listMessages.data);
    }
  }, [listMessages.isLoading, listMessages.data]);
  useEffect(() => {
    chatArea.current.scrollTo({
      top: chatArea.current.scrollHeight,
      behavior: "smooth",
    });
  }, [messageList]);

  const createMessage = useMutation({
    mutationFn: async (text) => {
      const response =
        await api.assistantsAssistantIdThreadsThreadIdMessagesPost({
          assistantId: assistantId,
          threadId: threadId,
          message: {
            role: "user",
            text: text,
            model: "gpt-4o",
          },
        });
      return response.data;
    },
    onSuccess: (data) => {
      setMessageList((prevMessageList) => [
        ...prevMessageList,
        {
          role: data.role,
          text: data.text,
          id: data.id,
        },
      ]);
    },
  });

  const messageS = (id, role, text) => {
    return (
      <Flex
        justify="flex-end"
        align="flex-start"
        direction="row"
        gap="sm"
        key={id}
      >
        <Paper
          shadow="sm"
          px="sm"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "blue.2"}
        >
          <ScrollArea type="auto" scrollbars="x">
            <Markdown>{text}</Markdown>
          </ScrollArea>
        </Paper>

        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
      </Flex>
    );
  };

  const messageR = (id, role, text) => {
    return (
      <Flex justify="flex-start" direction="row" gap="sm" key={id}>
        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
        <Paper
          shadow="sm"
          px="sm"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "green.2"}
        >
          <ScrollArea type="auto" scrollbars="x">
            <Markdown>{text}</Markdown>
          </ScrollArea>
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

  const sendMessage = async (text) => {
    setNewText("");
    setMessageList((prevMessageList) => [
      ...prevMessageList,
      { role: "user", text, id: Math.random() },
    ]);
    await createMessage.mutateAsync(text);
  };

  return (
    <>
      <Container size="xl">
        <Flex direction="column" justify="space-between" h="89vh">
          <ScrollArea
            viewportRef={chatArea}
            type="auto"
            offsetScrollbars
            scrollbars="y"
            style={{
              flex: 1,
            }}
          >
            <Stack spacing="md" py="lg">
              {messageList.map((item) => {
                if (item.role === "user") {
                  return messageS(item.id, item.role, item.text);
                } else {
                  return messageR(item.id, item.role, item.text);
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
              onKeyDown={async (e) => {
                if (e.key === "Enter") {
                  await sendMessage(e.currentTarget.value);
                }
              }}
              rightSection={
                <ActionIcon
                  variant="transparent"
                  aria-label="Settings"
                  disabled={createMessage.isLoading}
                  onClick={async (e) => {
                    await sendMessage(newText);
                  }}
                >
                  {createMessage.isLoading ? (
                    <Icon icon="svg-spinners:180-ring" />
                  ) : (
                    <Icon icon="tabler:send"></Icon>
                  )}
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
