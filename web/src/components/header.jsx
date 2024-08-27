import { Flex } from "@mantine/core";
import React from "react";
import { HeaderMenu } from "./header-menu";
import { ShowNavbarButton } from "./header-show-navbar-button";
import useStore from "../libs/store";

export default function Header({ hasNavBar }) {
  const isDarkMode = useStore((state) => state.isDarkMode);
  return (
    <>
      <Flex
        direction="row"
        justify="space-between"
        align="center"
        gap="lg"
        px="sm"
        h="100%"
        bg={isDarkMode ? "dark.8" : "gray.3"}
      >
        <ShowNavbarButton hasNavBar={hasNavBar} />

        <HeaderMenu />
      </Flex>
    </>
  );
}
