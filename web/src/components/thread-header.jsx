import { Icon } from "@iconify/react";
import { ActionIcon, Flex, Group, Tooltip } from "@mantine/core";
import React from "react";
import useStore from "../libs/store";
import { HeaderMenu } from "./header-menu";
import { ShowNavbarButton } from "./header-show-navbar-button";
import { ThreadAddButton } from "./thread-add-button";

export default function ThreadHeader() {
  const isDarkMode = useStore((state) => state.isDarkMode);
  const setThreadIsOpenSettings = useStore(
    (state) => state.setThreadIsOpenSettings,
  );

  return (
    <>
      <Flex
        direction="row"
        justify="space-between"
        align="center"
        gap="lg"
        px="sm"
        h="100%"
        bg={isDarkMode ? "dark.8" : "gray.4"}
      >
        <Group gap="2" align="center">
          <ShowNavbarButton hasNavBar={true} />
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