import "@mantine/core/styles.css";
import React from "react";
import ReactDOM from "react-dom/client";
import Layout from "../components/layout";
import Sidebar from "../components/thread-sidebar";

import { Container, Flex, LoadingOverlay } from "@mantine/core";

import { useEffect } from "react";
import { ThreadChatWindows } from "../components/thread-chat-windows";
import ThreadHeader from "../components/thread-header";
import { ThreadChatInput } from "../components/thread-input";
import { ThreadSettingsModal } from "../components/thread-settings";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";

function ThreadApp() {
  const url = new URL(window.location.href);

  const urlAssistantId = url.searchParams.get("assistant-id");
  const urlThreadId = url.searchParams.get("thread-id");

  const { getThread } = useQueryContext();

  const setThreadId = useStore((state) => state.setThreadId);
  const setAssistantId = useStore((state) => state.setAssistantId);

  useEffect(() => {
    if (urlAssistantId) {
      setAssistantId(urlAssistantId);
    } else {
      window.location.href = "/assistants.html";
    }
    if (urlThreadId) {
      setThreadId(urlThreadId);
    }
  }, []);

  return (
    <>
      <Container px="xs" h="95svh" fluid>
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
