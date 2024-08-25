import { Icon } from "@iconify/react";
import { ActionIcon, Tooltip } from "@mantine/core";
import React from "react";
import useStore from "../libs/store";

export function ShowNavbarButton({ hasNavBar }) {
  const toggleMobileSidebar = useStore((state) => state.toggleMobileSidebar);
  const toggleDesktopSidebar = useStore((state) => state.toggleDesktopSidebar);

  return (
    <>
      {hasNavBar ? (
        <>
          <Tooltip label="Toggle Sidebar" hiddenFrom="sm">
            <ActionIcon
              onClick={toggleMobileSidebar}
              variant="subtle"
              radius="lg"
              size="lg"
              hiddenFrom="sm"
            >
              <Icon icon="tabler:menu-3" />
            </ActionIcon>
          </Tooltip>
          <Tooltip label="Toggle Sidebar" visibleFrom="sm">
            <ActionIcon
              onClick={toggleDesktopSidebar}
              variant="subtle"
              size="lg"
              visibleFrom="sm"
            >
              <Icon icon="tabler:menu-3" />
            </ActionIcon>
          </Tooltip>
        </>
      ) : (
        <div></div>
      )}
    </>
  );
}
