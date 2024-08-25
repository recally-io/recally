import { Icon } from "@iconify/react";
import { ActionIcon, Flex, Group, Tooltip } from "@mantine/core";
import React from "react";
import useStore from "../libs/store";
import { HeaderMenu } from "./header-menu";
import { ShowNavbarButton } from "./header-show-navbar-button";
import { ThreadAddButton } from "./thread-add-button";

export default function ThreadHeader() {
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
        mx="sm"
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
