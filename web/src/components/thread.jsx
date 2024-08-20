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
  Textarea,
  Tooltip,
  useComputedColorScheme,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useRef, useState } from "react";
import Markdown from "react-markdown";
import avatarImgUrl from "../assets/avatar-1.png";
import { toastError } from "../libs/alert";
import { get, post } from "../libs/api";
import useStore from "../libs/store";

const url = new URL(window.location.href);

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
      const res = await get(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/messages`,
      );
      return res.data || [];
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
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/messages`,
        null,
        {
          role: "user",
          text: text,
          model: "gpt-4o",
        },
      );
      return res.data;
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
    onError: (error) => {
      toastError("Failed to send message: " + error.message);
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
            <Textarea
              placeholder="Send a message, use Shift + Enter to send."
              radius="lg"
              leftSection={menu()}
              leftSectionWidth={42}
              minRows={1}
              maxRows={5}
              autosize
              disabled={createMessage.isPending}
              onKeyDown={async (e) => {
                // Shift + Enter to send
                if (e.key === "Enter" && e.shiftKey === true) {
                  await sendMessage(e.currentTarget.value);
                }
              }}
              rightSection={
                <ActionIcon
                  variant="transparent"
                  aria-label="Settings"
                  onClick={async () => {
                    const text = newText;
                    setNewText("");
                    await sendMessage(text);
                  }}
                >
                  {createMessage.isPending ? (
                    <Icon icon="svg-spinners:180-ring" />
                  ) : (
                    <Icon icon="tabler:send"></Icon>
                  )}
                </ActionIcon>
              }
              value={newText}
              onChange={(e) => setNewText(e.currentTarget.value)}
            ></Textarea>
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
