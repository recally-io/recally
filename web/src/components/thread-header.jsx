import { Flex, Group, useComputedColorScheme } from "@mantine/core";
import React from "react";
import { HeaderMenu } from "./header-menu";
import { ShowNavbarButton } from "./header-show-navbar-button";

export default function ThreadHeader() {
  const computedColorScheme = useComputedColorScheme("light");

  return (
    <>
      <Flex
        direction="row"
        justify="space-between"
        align="center"
        gap="lg"
        px="sm"
        h="100%"
        bg={computedColorScheme === "dark" ? "dark.8" : "gray.4"}
      >
        <Group gap="2" align="center">
          <ShowNavbarButton hasNavBar={true} />
        </Group>
        <HeaderMenu />
      </Flex>
    </>
  );
}
