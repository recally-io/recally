import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/Layout";
import Sidebar from "../components/ThreadSidebar";

import { Container, Flex, LoadingOverlay } from "@mantine/core";

import { useForm } from "@mantine/form";
import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { ThreadChatWindows } from "../components/ThreadChatWindows";
import { ThreadChatInput } from "../components/ThreadInput";
import { ThreadSettingsModal } from "../components/ThreadSettings";
import { get } from "../libs/api";
import useStore from "../libs/store";

function ThreadApp() {
  const [threadId, setThreadId] = useStore((state) => [
    state.threadId,
    state.setThreadId,
  ]);
  const [assistantId, setAssistantId] = useStore((state) => [
    state.assistantId,
    state.setAssistantId,
  ]);

  useEffect(() => {
    const url = new URL(window.location.href);
    if (!url.searchParams.get("assistant-id")) {
      window.location.href = "/assistants.html";
    }
    setAssistantId(url.searchParams.get("assistant-id"));
    setThreadId(url.searchParams.get("thread-id"));
  }, []);

  const isLogin = useStore((state) => state.isLogin);
  const [chatModel, setChatModel] = useStore((state) => [
    state.setThreadChatModel,
    state.setThreadChatModel,
  ]);
  const setIsTitleGenerated = useStore(
    (state) => state.setThreadIsTitleGenerated,
  );
  const setMessageList = useStore((state) => state.setThreadMessageList);

  const settingsForm = useForm({
    initialValues: {
      name: "New Thread",
      description: "",
      systemPrompt: "",
      temperature: 0.7,
      maxToken: 4096,
      model: chatModel,
    },
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
      setMessageList(getThread.data.messages || []);
      if (getThread.data.model != "") {
        setChatModel(getThread.data.model);
      }
      setIsTitleGenerated(!!getThread.data.metadata.is_generated_title);
      window.document.title = getThread.data.name;
    }
  }, [getThread.data]);
  return (
    <>
      <Container px="xs" h="95svh" fluid>
        <LoadingOverlay visible={getThread.isLoading} />
        <Flex direction="column" justify="space-between" h="100%">
          <ThreadChatWindows settingsForm={settingsForm} />
          <ThreadChatInput />
        </Flex>
        <ThreadSettingsModal settingsForm={settingsForm} />
      </Container>
    </>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <Layout main={<ThreadApp />} nav={<Sidebar />} showNavBurger={true} />
  </React.StrictMode>,
);
