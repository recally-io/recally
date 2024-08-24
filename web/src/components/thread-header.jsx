import { Icon } from "@iconify/react";
import { ActionIcon, Flex, Group, Tooltip } from "@mantine/core";
import React from "react";
import useStore from "../libs/store";
import { HeaderMenu } from "./header-menu";
import { ThreadAddButton } from "./thread-add-button";

export default function ThreadHeader() {
  const [isSidebarOpen, toggleSidebar] = useStore((state) => [
    state.isSidebarOpen,
    state.toggleSidebar,
  ]);
  const setThreadIsOpenSettings = useStore(
    (state) => state.setThreadIsOpenSettings,
  );

  return (
    <>
      <Flex direction="row" justify="space-between" align="center" gap="lg">
        <Group gap="xs">
          <Tooltip label="Toggle Sidebar">
            <ActionIcon
              onClick={toggleSidebar}
              variant="subtle"
              radius="lg"
              size="lg"
            >
              {isSidebarOpen ? (
                <Icon icon="tabler:chevron-right" />
              ) : (
                <Icon icon="tabler:chevron-left" />
              )}
            </ActionIcon>
          </Tooltip>
          <ThreadAddButton />
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
        </Group>
        <HeaderMenu />
      </Flex>
    </>
  );
}
