import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Box,
  Button,
  Container,
  Divider,
  Flex,
  FocusTrap,
  Group,
  Image,
  Modal,
  NativeSelect,
  Textarea,
  Tooltip,
} from "@mantine/core";
import { getHotkeyHandler } from "@mantine/hooks";
import { useEffect, useState } from "react";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";
import { UploadButton } from "./upload-button";

export function ThreadChatInput() {
  const { sendThreadMessage, listModels, getThread } = useQueryContext();

  const [newText, setNewText] = useStore((state) => [
    state.threadNewText,
    state.setThreadNewText,
  ]);
  const [threadChatImages, setThreadChatImages] = useStore((state) => [
    state.threadChatImages,
    state.setThreadChatImages,
  ]);

  const [chatModel, setChatModel] = useState("");
  const [sendKey, setSendKey] = useState("enter");
  const [openedImage, setOpenedImage] = useState(null);

  useEffect(() => {
    if (getThread.data) {
      setChatModel(getThread.data.model);
    }
  }, [getThread.data]);

  const getSendKeyLabel = () => {
    switch (sendKey) {
      case "shift+enter":
        return (
          <Flex align="center" gap="1">
            <Icon icon="mdi:apple-keyboard-shift" />
            <Icon icon="mdi:keyboard-return" />
            <span>Send</span>
          </Flex>
        );
      case "mod+enter":
        return (
          <Flex align="center" gap="1">
            <Icon icon="mdi:apple-keyboard-command" />
            <Icon icon="mdi:keyboard-return" />
            <span>Send</span>
          </Flex>
        );
      case "enter":
        return (
          <Flex align="center" gap="1">
            <Icon icon="mdi:keyboard-return" />
            <span>Send</span>
          </Flex>
        );
      default:
        return "Send";
    }
  };

  const cycleSendKey = () => {
    const keys = ["shift+enter", "mod+enter", "enter"];
    const currentIndex = keys.indexOf(sendKey);
    const nextIndex = (currentIndex + 1) % keys.length;
    setSendKey(keys[nextIndex]);
  };

  const selectModel = () => {
    return (
      <NativeSelect
        variant="filled"
        maw="150px"
        radius="lg"
        size="xs"
        px="0"
        mx="0"
        value={chatModel}
        onChange={(e) => {
          setChatModel(e.target.value);
        }}
        data={listModels.data}
      />
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
        <Divider />
        <Group px="xs" justify="space-between">
          <Group align="center" gap="3">
            <Tooltip label="Upload image">
              <UploadButton useButton={true} />
            </Tooltip>
            {selectModel()}
          </Group>

          <Tooltip label="Click to change hotkey">
            <Button
              variant="transparent"
              size="xs"
              align="center"
              onClick={cycleSendKey}
            >
              {getSendKeyLabel()}
            </Button>
          </Tooltip>
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
              getHotkeyHandler([
                [
                  sendKey,
                  async () => {
                    if (newText.trim() !== "") {
                      const localText = newText.trim();
                      setNewText("");
                      await sendThreadMessage.mutateAsync({
                        model: chatModel,
                        text: localText,
                        images: threadChatImages,
                      });
                    }
                  },
                ],
              ])(e);
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
                disabled={newText === "" || sendThreadMessage.isPending}
                onClick={async () => {
                  await sendThreadMessage.mutateAsync();
                }}
              >
                {sendThreadMessage.isPending ? (
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
