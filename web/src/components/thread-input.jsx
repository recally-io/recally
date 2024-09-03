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
import { useQueryContext } from "../libs/query-context";

export function ThreadChatInput() {
  const {
    sendThreadMessage,
    listModels,
    getThread,
    getAssistant,
    getPresignedUrlMutation,
    uploadFileMutation,
  } = useQueryContext();

  const [text, setText] = useState("");
  const [images, setImages] = useState([]);
  const [chatModel, setChatModel] = useState("");

  const validSendKeys = [
    {
      value: "mod+enter",
      label: "⌘ + ↵ Send",
    },
    {
      value: "shift+enter",
      label: "⇧ + ↵ Send",
    },
    {
      value: "enter",
      label: "↵ Send",
    },
  ];
  const [sendKey, setSendKey] = useState("enter");
  const [openedImage, setOpenedImage] = useState(null);

  useEffect(() => {
    if (getThread.data) {
      setChatModel(getThread.data.model);
    }
  }, [getThread.data]);

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
      const localText = text.trim();
      const localImages = images;
      setText("");
      setImages([]);
      await sendThreadMessage.mutateAsync({
        model: chatModel,
        text: localText,
        images: localImages,
      });
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
          <Tooltip label={chatModel}>
            <ActionIcon variant="subtle" size="md">
              <Icon icon="tabler:robot" />
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
      <Popover
        width="150"
        position="bottom"
        trapFocus
        withArrow
        shadow="md"
        disabled={sendThreadMessage.isPending}
      >
        <Popover.Target>
          <Tooltip
            label={validSendKeys.find((key) => key.value === sendKey).label}
          >
            <ActionIcon variant="subtle" size="md">
              <Icon icon="tabler:keyboard" />
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
          borderRadius: "20px",
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
        <Divider />
        <Group px="xs" justify="space-between">
          <Group align="center" gap="3">
            {uploadImageButton()}
            {selectModel()}
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
      <Flex align="flex-end" gap="xs">
        <FocusTrap>
          <Textarea
            placeholder="Shift + Enter to send"
            radius="lg"
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
            rightSection={
              <ActionIcon
                variant="filled"
                radius="lg"
                aria-label="Settings"
                disabled={text === "" || sendThreadMessage.isPending}
                onClick={handleSendMessage}
              >
                {sendThreadMessage.isPending ? (
                  <Icon icon="svg-spinners:180-ring" />
                ) : (
                  <Icon icon="tabler:arrow-up"></Icon>
                )}
              </ActionIcon>
            }
            value={text}
            onChange={(e) => setText(e.currentTarget.value)}
            inputContainer={renderAttachmentTextArea}
          ></Textarea>
        </FocusTrap>
      </Flex>
    </Container>
  );
}
