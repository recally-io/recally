import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Autocomplete,
  Container,
  FileButton,
  Flex,
  FocusTrap,
  Textarea,
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

  const [images, setImages] = useState([]);

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
        metadata: { images: images },
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

      if (images.length > 0) {
        payload["metadata"] = { images: images };
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
      setImages([]);
    },
    onError: (error) => {
      toastError("Failed to send message: " + error.message);
    },
  });

  const getPresignedUrl = async ({
    assistantId,
    threadId,
    fileName,
    fileType,
    action,
    expiration,
  }) => {
    const params = new URLSearchParams({
      assistant_id: assistantId,
      thread_id: threadId,
      file_name: fileName,
      file_type: fileType,
      action: action,
      expiration: expiration,
    });
    console.log(`params: ${params}`);
    const response = await fetch(`/api/v1/files/presigned-urls?${params}`, {
      method: "GET",
      headers: { "Content-Type": "application/json" },
    });
    if (!response.ok) throw new Error("Failed to get presigned URL");
    const res = await response.json();
    console.log(`res: ${JSON.stringify(res)}`);
    return res.data;
  };

  const uploadFile = async ({ presignedUrl, file, publicUrl }) => {
    const response = await fetch(presignedUrl, {
      method: "PUT",
      body: file,
      headers: { "Content-Type": file.type },
    });
    if (!response.ok) throw new Error("Failed to upload file");
    return publicUrl;
  };

  const getPresignedUrlMutation = useMutation({
    mutationFn: getPresignedUrl,
    onSuccess: (data, variables) => handleUpload(data, variables.file),
    onError: (error) => {
      console.error("Error getting presigned URL:", error);
      toastError("Failed to get upload URL: " + error.message);
    },
  });

  const uploadFileMutation = useMutation({
    mutationFn: uploadFile,
    onSuccess: (data) => {
      setImages((prevImages) => [...prevImages, data]);
      console.log("Updated images:", [...images, data]);
    },
    onError: (error) => {
      console.error("Error uploading file:", error);
      toastError("Failed to upload file: " + error.message);
    },
  });

  const handleFileChange = async (files) => {
    if (!files) return;
    for (const file of files) {
      getPresignedUrlMutation.mutate({
        assistantId: assistant.id,
        threadId: thread.id,
        fileName: file.name,
        fileType: file.type,
        action: "put",
        expiration: 3600,
        file, // Pass the file to be used in onSuccess
      });
    }
  };

  const handleUpload = async (data, file) => {
    console.log(`uploading file: ${JSON.stringify(data)}`);
    uploadFileMutation.mutate({
      presignedUrl: data.presigned_url,
      file,
      publicUrl: data.public_url,
    });
  };

  const fileInputButton = () => {
    return (
      <FileButton
        onChange={handleFileChange}
        accept="image/*"
        multiple
        disabled={
          sendMessage.isPending ||
          getPresignedUrlMutation.isPending ||
          uploadFileMutation.isPending
        }
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
            leftSection={fileInputButton()}
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
          ></Textarea>
        </FocusTrap>
      </Flex>
    </Container>
  );
}
