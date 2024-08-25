import { Flex } from "@mantine/core";
import React from "react";
import { HeaderMenu } from "./header-menu";
import { ShowNavbarButton } from "./header-show-navbar-button";

export default function Header({ hasNavBar }) {
  return (
    <>
      <Flex
        direction="row"
        justify="space-between"
        align="center"
        gap="lg"
        px="sm"
        m="0"
      >
        <ShowNavbarButton hasNavBar={hasNavBar} />

        <HeaderMenu />
      </Flex>
    </>
  );
}
