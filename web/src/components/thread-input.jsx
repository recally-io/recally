import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Autocomplete,
  Box,
  Container,
  Flex,
  FocusTrap,
  Group,
  Image,
  Modal,
  Textarea,
  Tooltip,
} from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toastError } from "../libs/alert";
import {
  listModels,
  listModelsKey,
  listTools,
  listToolsKey,
  post,
} from "../libs/api";
import useStore from "../libs/store";
import { UploadButton } from "./upload-button";

export function ThreadChatInput() {
  const isLogin = useStore((state) => state.isLogin);

  const assistant = useStore((state) => state.assistant);
  const thread = useStore((state) => state.thread);
  const setThread = useStore((state) => state.setThread);

  const [isShowModelSelecter, setIsShowModelSelecter] = useStore((state) => [
    state.threadIsOpenModelSelecter,
    state.setThreadIsOpenModelSelecter,
  ]);
  const [chatModel, setChatModel] = useStore((state) => [
    state.threadChatModel,
    state.setThreadChatModel,
  ]);
  const [newText, setNewText] = useStore((state) => [
    state.threadNewText,
    state.setThreadNewText,
  ]);
  const addThreadMessage = useStore((state) => state.addThreadMessage);
  const setModels = useStore((state) => state.setThreadModels);
  const setTools = useStore((state) => state.setThreadTools);
  const [modelSelecterValue, setModelSelecterValue] = useState("");

  const threadChatImages = useStore((state) => state.threadChatImages);
  const setThreadChatImages = useStore((state) => state.setThreadChatImages);

  const [openedImage, setOpenedImage] = useState(null);

  useEffect(() => {
    if (newText === "@") {
      setIsShowModelSelecter(true);
    }
  }, [newText]);

  const createThread = useMutation({
    mutationFn: async (data) => {
      const res = await post(
        `/api/v1/assistants/${assistant.id}/threads`,
        null,
        data,
      );
      setThread(res.data);
      return res.data;
    },
  });

  const listModelsQuery = useQuery({
    queryKey: listModelsKey,
    queryFn: async () => {
      const res = await listModels();
      setModels(res);
      return res;
    },
    enabled: isLogin,
  });

  useQuery({
    queryKey: [listToolsKey],
    queryFn: async () => {
      const res = await listTools();
      setTools(res);
      return res;
    },
    enabled: isLogin,
  });

  const sendMessage = useMutation({
    mutationFn: async () => {
      let text = newText;
      if (text.startsWith("@")) {
        text = text.replace(/^@[^ ]+\s*/, "");
      }
      setNewText("");
      addThreadMessage({
        role: "user",
        text,
        id: Math.random(),
        metadata: { images: [...threadChatImages] },
      });

      const isNewThread = !thread.id;
      let newThreadId = thread.id;
      if (isNewThread) {
        newThreadId = crypto.randomUUID();
        await createThread.mutateAsync({
          id: newThreadId,
          name: "New Thread",
          description: assistant.description,
          system_prompt: assistant.systemPrompt,
          model: assistant.model,
          metadata: {
            is_generated_title: false,
            tools: assistant.metadata.tools,
          },
        });
      }
      let payload = {
        role: "user",
        text: text,
        model: chatModel,
      };

      if (threadChatImages.length > 0) {
        payload["metadata"] = { images: threadChatImages };
      }

      const res = await post(
        `/api/v1/assistants/${assistant.id}/threads/${newThreadId}/messages`,
        null,
        payload,
      );
      return res.data;
    },
    onSuccess: (data) => {
      addThreadMessage(data);
      setThreadChatImages([]);
    },
    onError: (error) => {
      toastError("Failed to send message: " + error.message);
    },
  });

  const renderAttachmentTextArea = (children) => {
    return (
      <Flex
        direction="column"
        gap="2"
        bd="1px solid primary"
        p="2"
        style={{
          borderRadius: "20px",
        }}
      >
        {threadChatImages.length > 0 && (
          <Group px="md">
            {threadChatImages.map((imgUrl, index) => (
              <Box key={index} style={{ position: "relative" }}>
                <Image
                  src={imgUrl}
                  width={30}
                  height={30}
                  fit="contain"
                  onClick={() => {
                    // Function to show large image
                    console.log("Show large image:", imgUrl);
                    setOpenedImage(imgUrl);
                  }}
                  style={{ cursor: "pointer" }}
                />
                <Modal
                  opened={openedImage === imgUrl}
                  onClose={() => setOpenedImage(null)}
                  size="xl"
                >
                  <Image
                    src={imgUrl}
                    alt={imgUrl.split("/").pop()}
                    fit="contain"
                  />
                </Modal>
                <Tooltip label="Remove image">
                  <ActionIcon
                    size="xs"
                    color="danger"
                    variant="subtle"
                    style={{
                      position: "absolute",
                      top: 0,
                      right: 0,
                    }}
                    onClick={() => {
                      // Function to remove image
                      setThreadChatImages((prevImages) =>
                        prevImages.filter((image) => image !== imgUrl),
                      );
                    }}
                  >
                    <Icon icon="tabler:x" size={10} />
                  </ActionIcon>
                </Tooltip>
              </Box>
            ))}
          </Group>
        )}
        <div>{children}</div>
      </Flex>
    );
  };

  return (
    <Container
      w="100%"
      style={{
        position: "sticky",
        bottom: 0,
      }}
    >
      <Flex align="flex-end" gap="xs">
        {isShowModelSelecter && (
          <FocusTrap active={isShowModelSelecter}>
            <Autocomplete
              label="Talk to model"
              placeholder="Type to select model"
              data={[...new Set(listModelsQuery.data)]}
              dropdownOpened={isShowModelSelecter}
              radius="lg"
              fz="16px"
              leftSectionPointerEvents="none"
              leftSection={<Icon icon="tabler:robot" />}
              value={modelSelecterValue}
              onChange={setModelSelecterValue}
              onOptionSubmit={(v) => {
                setChatModel(v);
                setNewText(`@${v} `);
                setIsShowModelSelecter(false);
              }}
            />
          </FocusTrap>
        )}
        <FocusTrap active={!isShowModelSelecter}>
          <Textarea
            placeholder="Shift + Enter to send"
            radius="lg"
            leftSection={<UploadButton useButton={true} />}
            minRows={1}
            maxRows={5}
            autosize
            w="100%"
            disabled={sendMessage.isPending}
            onKeyDown={async (e) => {
              // Shift + Enter to send
              if (e.key === "Enter" && e.shiftKey === true) {
                await sendMessage.mutateAsync();
              }
            }}
            styles={{
              input: {
                border: "none",
              },
            }}
            rightSection={
              <ActionIcon
                variant="filled"
                radius="lg"
                aria-label="Settings"
                disabled={newText === "" || sendMessage.isPending}
                onClick={async () => {
                  await sendMessage.mutateAsync();
                }}
              >
                {sendMessage.isPending ? (
                  <Icon icon="svg-spinners:180-ring" />
                ) : (
                  <Icon icon="tabler:arrow-up"></Icon>
                )}
              </ActionIcon>
            }
            value={newText}
            onChange={(e) => setNewText(e.currentTarget.value)}
            inputContainer={renderAttachmentTextArea}
          ></Textarea>
        </FocusTrap>
      </Flex>
    </Container>
  );
}
