import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Box,
  Container,
  Divider,
  FileButton,
  Flex,
  FocusTrap,
  Group,
  Image,
  Modal,
  Popover,
  Select,
  Textarea,
  Tooltip,
} from "@mantine/core";
import { IMAGE_MIME_TYPE } from "@mantine/dropzone";
import { getHotkeyHandler } from "@mantine/hooks";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQueryContext } from "../libs/query-context";
import useStore, { defaultThreadSettings } from "../libs/store";
import { ThreadAddButton, ThreadSettingsButton } from "./thread-add-button";

export function ThreadChatInput() {
  const navigate = useNavigate();
  const params = useParams();
  const threadId = params.threadId;
  const assistantId = params.assistantId;

  const {
    sendThreadMessage,
    listModels,
    createThread,
    getThread,
    getAssistant,
    getPresignedUrlMutation,
    uploadFileMutation,
    generateThreadTitle,
  } = useQueryContext();
  const messageList = useStore((state) => state.threadMessageList);

  const [text, setText] = useState("");
  const [images, setImages] = useState([]);
  const [chatModel, setChatModel] = useState(defaultThreadSettings.model);

  const validSendKeys = [
    {
      value: "mod+enter",
      label: "⌘ + ↵",
    },
    {
      value: "shift+enter",
      label: "⇧ + ↵",
    },
    {
      value: "enter",
      label: "↵",
    },
  ];
  const [sendKey, setSendKey] = useState("enter");
  const [openedImage, setOpenedImage] = useState(null);

  const [isLoading, setLoading] = useState(false);

  useEffect(() => {
    if (getAssistant.data) {
      setChatModel(getAssistant.data.model);
    }
    if (getThread.data) {
      setChatModel(getThread.data.model);
    }
  }, [getThread.data, getAssistant.data]);

  const handleUploadImage = async (files) => {
    if (!files) return;
    for (const file of files) {
      console.log(`file: ${file.name}, type: ${file.type}`);

      // get presigned url
      const preSignedUrlRes = await getPresignedUrlMutation.mutateAsync({
        assistantId: getAssistant.data.id,
        threadId: getThread.data.id,
        fileName: file.name,
        fileType: file.type,
      });
      // upload file
      const uploadRes = await uploadFileMutation.mutateAsync({
        preSignedURL: preSignedUrlRes.preSignedURL,
        file,
        publicUrl: preSignedUrlRes.publicUrl,
      });
      setImages((prevImages) => [...prevImages, uploadRes]);
    }
  };

  const handleSendMessage = async () => {
    if (text.trim() !== "") {
      let newThreadId = threadId;
      if (!newThreadId) {
        newThreadId = crypto.randomUUID();
        const assistant = getAssistant.data;
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

      setLoading(true);
      const localText = text.trim();
      const localImages = images;
      setText("");
      setImages([]);
      await sendThreadMessage.mutateAsync({
        assistantId: assistantId,
        threadId: newThreadId,
        model: chatModel,
        text: localText,
        images: localImages,
      });
      setLoading(false);

      // if threadId is not present, navigate to new thread
      if (!threadId) {
        navigate(`/assistants/${assistantId}/threads/${newThreadId}`, {});
      }

      if (
        messageList.length >= 2 &&
        getThread.data &&
        !getThread.data.metadata.is_generated_title
      ) {
        await generateThreadTitle.mutateAsync();
      }
    }
  };

  const uploadImageButton = () => {
    return (
      <FileButton
        onChange={handleUploadImage}
        accept={[...IMAGE_MIME_TYPE]}
        multiple
        disabled={sendThreadMessage.isPending}
      >
        {(props) => (
          <Tooltip label="Upload Image">
            <ActionIcon {...props} variant="subtle" size="md">
              <Icon icon="tabler:photo" />
            </ActionIcon>
          </Tooltip>
        )}
      </FileButton>
    );
  };

  const selectModel = () => {
    return (
      <Popover
        width="200"
        position="bottom"
        trapFocus
        withArrow
        shadow="md"
        disabled={sendThreadMessage.isPending}
      >
        <Popover.Target>
          <Tooltip label={"Select Model: " + chatModel}>
            <ActionIcon variant="subtle" size="md">
              <Icon icon="tabler:at" />
            </ActionIcon>
          </Tooltip>
        </Popover.Target>
        <Popover.Dropdown p="0">
          <Select
            variant="filled"
            dropdownOpened={true}
            value={chatModel}
            onChange={setChatModel}
            data={[...new Set(listModels.data)]}
          />
        </Popover.Dropdown>
      </Popover>
    );
  };

  const selectSendHotKey = () => {
    return (
      <Popover width="100" position="bottom" trapFocus withArrow shadow="md">
        <Popover.Target>
          <Tooltip
            label={validSendKeys.find((key) => key.value === sendKey).label}
          >
            <ActionIcon
              variant="transparent"
              radius="md"
              disabled={text === "" && !isLoading}
            >
              <Icon icon="tabler:chevron-down" />
            </ActionIcon>
          </Tooltip>
        </Popover.Target>
        <Popover.Dropdown p="0">
          <Select
            variant="filled"
            radius="md"
            dropdownOpened={true}
            value={sendKey}
            onChange={setSendKey}
            data={validSendKeys}
          />
        </Popover.Dropdown>
      </Popover>
    );
  };

  const renderAttachmentTextArea = (children) => {
    return (
      <Flex
        direction="column"
        gap="2"
        bd="1px solid primary"
        p="2"
        style={{
          borderRadius: "10px",
        }}
      >
        {images.length > 0 && (
          <Group px="md">
            {images.map((imgUrl, index) => (
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
                      setImages((prevImages) =>
                        prevImages.filter((image) => image !== imgUrl)
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
        <Divider />
        <Group px="2" justify="space-between">
          <Group align="center" gap="3">
            {uploadImageButton()}
            {selectModel()}
            <Tooltip label="Select prompt">
              <ActionIcon size="md" variant="subtle" radius="lg">
                <Icon icon="tabler:slash" />
              </ActionIcon>
            </Tooltip>
            {/* <Divider orientation="vertical" /> */}
            <Tooltip label="Voice Chat">
              <ActionIcon size="md" variant="subtle" radius="lg">
                <Icon icon="tabler:microphone" />
              </ActionIcon>
            </Tooltip>
            <Divider orientation="vertical" />
            <ThreadSettingsButton />
            <ThreadAddButton />
            <Divider orientation="vertical" />
          </Group>
          <Group size="xs" p="0" gap="0">
            <ActionIcon
              variant="filled"
              radius="md"
              aria-label="Settings"
              disabled={text === "" && !isLoading}
              onClick={handleSendMessage}
            >
              {isLoading ? (
                <Icon icon="svg-spinners:180-ring" />
              ) : (
                <Icon icon="tabler:send-2" />
              )}
            </ActionIcon>
            <Divider orientation="vertical" />
            {selectSendHotKey()}
          </Group>
        </Group>
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
      <FocusTrap>
        <Textarea
          placeholder={`use ${
            validSendKeys.find((key) => key.value === sendKey).label
          } to send`}
          radius="md"
          minRows={1}
          maxRows={5}
          autosize
          w="100%"
          disabled={sendThreadMessage.isPending}
          onKeyDown={async (e) => {
            e.stopPropagation();
            getHotkeyHandler([[sendKey, handleSendMessage]])(e);
          }}
          styles={{
            input: {
              border: "none",
            },
          }}
          value={text}
          onChange={(e) => setText(e.currentTarget.value)}
          inputContainer={renderAttachmentTextArea}
        ></Textarea>
      </FocusTrap>
    </Container>
  );
}
