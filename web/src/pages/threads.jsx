import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";
import Sidebar from "../components/thread-sidebar";

import { Container, Flex, LoadingOverlay } from "@mantine/core";

import { useQuery } from "@tanstack/react-query";
import { useEffect } from "react";
import { ThreadChatWindows } from "../components/thread-chat-windows";
import ThreadHeader from "../components/thread-header";
import { ThreadChatInput } from "../components/thread-input";
import { ThreadSettingsModal } from "../components/thread-settings";
import { get } from "../libs/api";
import useStore from "../libs/store";

function ThreadApp() {
  const url = new URL(window.location.href);

  const assistantId = url.searchParams.get("assistant-id");
  const urlThreadId = url.searchParams.get("thread-id");
  const [threadId, setThreadId] = useStore((state) => [
    state.threadId,
    state.setThreadId,
  ]);

  useEffect(() => {
    if (!assistantId) {
      console.log("no assistant id");
      window.location.href = "/assistants.html";
    }
  }, [assistantId]);

  useEffect(() => {
    setThreadId(urlThreadId);
  }, [urlThreadId]);

  const isLogin = useStore((state) => state.isLogin);
  const setChatModel = useStore((state) => state.setThreadChatModel);
  const setIsTitleGenerated = useStore(
    (state) => state.setThreadIsTitleGenerated,
  );
  const setMessageList = useStore((state) => state.setThreadMessageList);

  const setThread = useStore((state) => state.setThread);
  const setAssistant = useStore((state) => state.setAssistant);

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
    if (getThread.data && threadId) {
      setThread(getThread.data);
      setMessageList(getThread.data.messages || []);
      setChatModel(getThread.data.model);
      setIsTitleGenerated(!!getThread.data.metadata.is_generated_title);
      window.document.title = getThread.data.name;
    }
  }, [getThread.data, threadId]);

  useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}`);
      setAssistant(res.data);
      return res.data;
    },
    enabled: isLogin && !!assistantId,
  });

  return (
    <>
      <Container px="xs" h="95svh">
        <LoadingOverlay visible={getThread.isLoading} />
        <Flex direction="column" justify="space-between" h="100%">
          <ThreadChatWindows />
          <ThreadChatInput />
        </Flex>
        <ThreadSettingsModal />
      </Container>
    </>
  );
}

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <Layout main={<ThreadApp />} nav={<Sidebar />} header={<ThreadHeader />} />
  </React.StrictMode>,
);
