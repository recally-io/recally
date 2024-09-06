import { Container, Flex } from "@mantine/core";

import { Layout } from "../components/layout";
import { ThreadChatWindows } from "../components/thread-chat-windows";
import ThreadHeader from "../components/thread-header";
import { ThreadChatInput } from "../components/thread-input";
import { ThreadSettingsModal } from "../components/thread-settings";
import ThreadSidebar from "../components/thread-sidebar";

export default function Threads() {
  const main = () => {
    return (
      <Container px="xs" h="95svh" fluid>
        <Flex direction="column" justify="space-between" h="100%">
          <ThreadChatWindows />
          <ThreadChatInput />
        </Flex>
        <ThreadSettingsModal />
      </Container>
    );
  };

  return (
    <>
      <Layout main={main()} nav={<ThreadSidebar />} header={<ThreadHeader />} />
    </>
  );
}
