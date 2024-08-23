import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Autocomplete,
  Container,
  FileButton,
  Flex,
  FocusTrap,
  Textarea,
  Tooltip,
} from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toastError } from "../libs/alert";
import { post, get } from "../libs/api";
import useStore from "../libs/store";

export function ThreadChatInput({}) {
  const isLogin = useStore((state) => state.isLogin);
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

  const setModels = useStore((state) => state.setThreadModels);

  const setThreadIsOpenSettings = useStore(
    (state) => state.setThreadIsOpenSettings,
  );

  const [modelSelecterValue, setModelSelecterValue] = useState("");

  useEffect(() => {
    if (newText === "@") {
      setIsShowModelSelecter(true);
    }
  }, [newText]);

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
    if (listModels.data) {
      setModels(listModels.data);
    }
  }, [listModels.data]);

  const sendMessage = useMutation({
    mutationFn: async () => {
      let text = newText;
      if (text.startsWith("@")) {
        text = text.replace(/^@[^ ]+\s*/, "");
      }
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
          model: chatModel,
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
        disabled={sendMessage.isPending}
      >
        {(props) => (
          <ActionIcon {...props} variant="subtle" radius="lg">
            <Icon icon="tabler:file-upload"></Icon>
          </ActionIcon>
        )}
      </FileButton>
    );
  };

  return (
    <Flex align="center">
      <Tooltip label="Settings">
        <ActionIcon
          size="lg"
          variant="subtle"
          radius="lg"
          onClick={() => setThreadIsOpenSettings(true)}
        >
          <Icon icon="tabler:settings"></Icon>
        </ActionIcon>
      </Tooltip>

      <Container
        w="100%"
        style={{
          position: "sticky",
          bottom: 0,
        }}
      >
        {isShowModelSelecter && (
          <FocusTrap active={isShowModelSelecter}>
            <Autocomplete
              label="Talk to model"
              placeholder="Type to select model"
              data={[...new Set(listModels.data)]}
              dropdownOpened={isShowModelSelecter}
              radius="lg"
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
            leftSection={fileInputButton()}
            minRows={1}
            maxRows={5}
            autosize
            disabled={sendMessage.isPending}
            onKeyDown={async (e) => {
              // Shift + Enter to send
              if (e.key === "Enter" && e.shiftKey === true) {
                await sendMessage.mutateAsync();
              }
            }}
            rightSection={
              <ActionIcon
                variant="transparent"
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
          ></Textarea>
        </FocusTrap>
      </Container>
    </Flex>
  );
}
