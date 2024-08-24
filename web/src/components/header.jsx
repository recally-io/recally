import { Icon } from "@iconify/react";
import { Button, Flex } from "@mantine/core";
import React from "react";
import useStore from "../libs/store";
import { HeaderMenu } from "./header-menu";

export default function Header({ showNavBurger }) {
  const toggleSidebar = useStore((state) => state.toggleSidebar);

  return (
    <>
      <Flex direction="row" justify="space-between" align="center" gap="lg">
        {showNavBurger ? (
          <Button onClick={toggleSidebar} variant="transparent" size="md">
            <Icon icon="tabler:layout-sidebar" />
          </Button>
        ) : (
          <div></div>
        )}

        <HeaderMenu />
      </Flex>
    </>
  );
}
