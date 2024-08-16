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
import Cookie from "js-cookie";
import avatarImgUrl from "../assets/avatar-1.png";
import useStore from "../libs/store";

export default function Header({ opened, toggle, showNavBurger }) {
  const isLogin = useStore((state) => state.isLogin);
  const setIsLogin = useStore((state) => state.setIsLogin);
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
        <Flex direction="row" justify="flex-end" align="center">
          <Button
            variant="transparent"
            size="sm"
            onClick={() => {
              setColorScheme(colorScheme === "dark" ? "light" : "dark");
            }}
          >
            {colorScheme === "dark" ? (
              <Icon icon="tabler:sun" color="white" />
            ) : (
              <Icon icon="tabler:moon-filled" color="black" />
            )}
          </Button>
          <Menu shadow="xl" trigger="hover" transition="slide-up" withArrow>
            <Menu.Target>
              <Button variant="transparent" size="sm">
                <Avatar size="sm" radius="lg" src={avatarImgUrl} />
              </Button>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Label>Vibrain</Menu.Label>
              {isLogin && (
                <Menu.Item
                  color="red"
                  leftSection={<Icon icon="tabler:logout" />}
                  onClick={() => {
                    setIsLogin(false);
                    Cookie.remove("token");
                    window.location.href = "/";
                  }}
                >
                  Logout
                </Menu.Item>
              )}
            </Menu.Dropdown>
          </Menu>
        </Flex>
      </Flex>
    </>
  );
}
