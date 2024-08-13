import { Icon } from "@iconify/react";
import {
  Avatar,
  Burger,
  Button,
  Flex,
  Menu,
  NavLink,
  useMantineColorScheme,
} from "@mantine/core";
import avatarImgUrl from "../assets/avatar-1.png";

export default function Header({ opened, toggle, showNavBurger }) {
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  return (
    <>
      <Flex direction="row" justify="space-between" align="center" px="md">
        {showNavBurger && (
          <Burger
            opened={opened}
            onClick={toggle}
            size="sm"
            aria-label="Toggle navigation"
          />
        )}
        <NavLink href="/" label="Vibrain"></NavLink>
        <NavLink href="/assistants.html" label="Assistants"></NavLink>
        <Flex direction="row" justify="space-around" align="center">
          <Button
            variant="transparent"
            color="white"
            onClick={() => {
              setColorScheme(colorScheme === "dark" ? "light" : "dark");
            }}
          >
            {colorScheme === "dark" ? (
              <Icon icon="tabler:sun" color="white" width={18} height={18} />
            ) : (
              <Icon
                icon="tabler:moon-filled"
                color="black"
                width={18}
                height={18}
              />
            )}
          </Button>
          <Menu
            opened={opened}
            onClose={toggle}
            position="right"
            shadow="xl"
            transition="slide-up"
            withArrow
            control={<Burger opened={opened} size="sm" />}
          >
            <Avatar size="sm" radius="lg" src={avatarImgUrl} />
          </Menu>
        </Flex>
      </Flex>
    </>
  );
}
