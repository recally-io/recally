import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Avatar,
  Button,
  Container,
  FileButton,
  Flex,
  Group,
  LoadingOverlay,
  Modal,
  NativeSelect,
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
import { get, post, queryClient } from "../libs/api";
import useStore from "../libs/store";

const url = new URL(window.location.href);

export default function ChatWindowsComponent() {
  const isLogin = useStore((state) => state.isLogin);
  const colorScheme = useComputedColorScheme("light");

  const [settingsOpened, { open: openSettings, close: closeSettings }] =
    useDisclosure(false);
  const settingsForm = useForm({
    initialValues: {
      name: "New Thread",
      description: "",
      systemPrompt: "",
      temperature: 0.7,
      maxToken: 4096,
      model: "gpt-4o",
    },
  });

  const assistantId = url.searchParams.get("assistant-id");
  let threadId = url.searchParams.get("thread-id");
  const [newText, setNewText] = useState("");

  const [messageList, setMessageList] = useState([]);
  const chatArea = useRef(null);

  const [fileContent, setFileContent] = useState("");

  const getAssistant = useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}`);
      return res.data;
    },
    enabled: isLogin && !!assistantId,
  });

  const getThread = useQuery({
    queryKey: ["get-thread", threadId],
    queryFn: async () => {
      const res = await get(
        `/api/v1/assistants/${assistantId}/threads/${threadId}`,
      );
      return res.data || {};
    },
    enabled: isLogin && !!threadId && !!assistantId,
  });

  useEffect(() => {
    if (getThread.data) {
      settingsForm.setValues(getThread.data);
      setMessageList(getThread.data.messages);
      window.document.title = `Chat with ${getThread.data.name}`;
    }
  }, [getThread.data]);

  const createThread = useMutation({
    mutationFn: async (data) => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads`,
        null,
        data,
      );
      return res.data;
    },
  });

  const listModels = useQuery({
    queryKey: ["list-assistants-models"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants/models");
      return res.data || [];
    },
    enabled: isLogin,
  });

  useEffect(() => {
    chatArea.current.scrollTo({
      top: chatArea.current.scrollHeight,
      behavior: "smooth",
    });
    console.log("message list: ", messageList);
    console.log(getThread.data);

    const generate = async () => {
      await generateTitle.mutateAsync();
    };

    if (getThread.data) {
      if (
        messageList.length >= 4 &&
        !getThread.data.metadata.is_generated_title
      ) {
        console.log("Generate title");
        generate();
      }
    }
  }, [messageList]);

  const createMessage = useMutation({
    mutationFn: async () => {
      const text = newText;
      setNewText("");
      setMessageList((prevMessageList) => [
        ...prevMessageList,
        { role: "user", text, id: Math.random() },
      ]);

      const isNewThread = threadId === null;
      if (isNewThread) {
        threadId = crypto.randomUUID();
        let data = settingsForm.getValues();
        data.id = threadId;
        await createThread.mutateAsync(data);
      }

      const res = await post(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/messages`,
        null,
        {
          role: "user",
          text: text,
          model: "gpt-4o",
        },
      );

      if (isNewThread) {
        window.location.href = `/threads.html?assistant-id=${assistantId}&thread-id=${threadId}`;
      }

      return res.data;
    },
    onSuccess: (data) => {
      setMessageList((prevMessageList) => [...prevMessageList, data]);
    },
    onError: (error) => {
      toastError("Failed to send message: " + error.message);
    },
  });

  const generateTitle = useMutation({
    mutationFn: async () => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/generate-title`,
        null,
        {},
      );
      return res.data;
    },
    onSuccess: (data) => {
      settingsForm.setFieldValue("name", data.name);
      queryClient.invalidateQueries(["get-thread", threadId]);
    },
  });

  const messageS = (message) => {
    return (
      <Flex
        justify="flex-end"
        align="flex-start"
        direction="row"
        gap="sm"
        key={message.id}
      >
        <Flex align="flex-end" direction="column" maw="90%">
          <Text size="lg" variant="gradient">
            You
          </Text>
          <Paper
            shadow="sm"
            px="sm"
            w="100%"
            radius="lg"
            bg={colorScheme === "dark" ? "" : "blue.2"}
          >
            <ScrollArea type="auto" scrollbars="x">
              <Markdown>{message.text}</Markdown>
            </ScrollArea>
          </Paper>
        </Flex>

        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
      </Flex>
    );
  };

  const messageR = (message) => {
    return (
      <Flex justify="flex-start" direction="row" gap="sm" key={message.id}>
        <Avatar size="sm" radius="lg" color="cyan" variant="filled">
          <Icon icon="tabler:robot" />
        </Avatar>
        <Flex align="flex-start" direction="column" maw="90%">
          <Text size="lg" variant="gradient">
            {message.model}
          </Text>
          <Paper
            shadow="sm"
            px="sm"
            w="100%"
            radius="lg"
            bg={colorScheme === "dark" ? "" : "green.2"}
          >
            <ScrollArea type="auto" scrollbars="x">
              <Markdown>{message.text}</Markdown>
            </ScrollArea>
          </Paper>
        </Flex>
      </Flex>
    );
  };

  const modalSettings = () => {
    return (
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
            <NativeSelect
              label="Model"
              key={settingsForm.key("model")}
              {...settingsForm.getInputProps("model")}
              onChange={(e) => {
                settingsForm.setFieldValue("model", e.target.value);
              }}
              data={listModels.data}
            />
          </Stack>
          <FileButton
            size="sm"
            variant="transparent"
            multiple
            leftSection={<Icon icon="tabler:upload"></Icon>}
          >
            {(props) => <Button {...props}>Upload image</Button>}
          </FileButton>
          <Group justify="flex-end" mt="md">
            <Button type="submit" onClick={closeSettings}>
              Submit
            </Button>
            <Button type="button" onClick={closeSettings}>
              Cancel
            </Button>
          </Group>
        </form>
      </Modal>
    );
  };

  const handleFileChange = async (file) => {
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        setFileContent(e.target.result);
      };
      reader.readAsText(file);
    }
  };

  const fileInputButton = () => {
    return (
      <FileButton
        onChange={handleFileChange}
        // accept="image/png,image/jpeg"
        // multiple
        disabled={createMessage.isPending}
      >
        {(props) => (
          <ActionIcon {...props} variant="subtle" radius="lg">
            <Icon icon="tabler:file-upload"></Icon>
          </ActionIcon>
        )}
      </FileButton>
    );
  };

  const textInput = () => {
    return (
      <Container
        w="100%"
        style={{
          position: "sticky",
          bottom: 0,
        }}
      >
        <Textarea
          placeholder="Shift + Enter to send"
          radius="lg"
          leftSection={fileInputButton()}
          minRows={1}
          maxRows={5}
          autosize
          disabled={createMessage.isPending}
          onKeyDown={async (e) => {
            // Shift + Enter to send
            if (e.key === "Enter" && e.shiftKey === true) {
              await createMessage.mutateAsync();
            }
          }}
          rightSection={
            <ActionIcon
              variant="transparent"
              aria-label="Settings"
              disabled={newText === "" || createMessage.isPending}
              onClick={async () => {
                await createMessage.mutateAsync();
              }}
            >
              {createMessage.isPending ? (
                <Icon icon="svg-spinners:180-ring" />
              ) : (
                <Icon icon="tabler:arrow-up"></Icon>
              )}
            </ActionIcon>
          }
          value={newText}
          onChange={(e) => setNewText(e.currentTarget.value)}
        ></Textarea>
      </Container>
    );
  };

  return (
    <>
      <Container px="xs" h="95svh">
        <LoadingOverlay visible={getThread.isLoading} />
        <Flex direction="column" justify="space-between" h="100%">
          <ScrollArea
            viewportRef={chatArea}
            type="scroll"
            offsetScrollbars
            scrollbarSize="4"
            scrollbars="y"
          >
            <Stack spacing="md" py="lg">
              {messageList.map((item) => {
                if (item.role === "user") {
                  return messageS(item);
                } else {
                  return messageR(item);
                }
              })}
            </Stack>
          </ScrollArea>
          <Flex align="center">
            <Tooltip label="Settings">
              <ActionIcon
                size="lg"
                variant="subtle"
                radius="lg"
                onClick={openSettings}
              >
                <Icon icon="tabler:settings"></Icon>
              </ActionIcon>
            </Tooltip>
            {textInput()}
          </Flex>
        </Flex>
        {modalSettings()}
      </Container>
    </>
  );
}
