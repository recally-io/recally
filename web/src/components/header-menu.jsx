import { Icon } from "@iconify/react";
import { Avatar, Button, Menu, useMantineColorScheme } from "@mantine/core";
import Cookie from "js-cookie";
import React, { useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuthContext } from "../libs/auth-context";
import useStore from "../libs/store";

export function HeaderMenu() {
  const navigate = useNavigate();
  const { isLogin } = useAuthContext();
  const { colorScheme, setColorScheme } = useMantineColorScheme();

  const [isDarkMode, setIsDarkMode] = useStore((state) => [
    state.isDarkMode,
    state.setIsDarkMode,
  ]);

  useEffect(() => {
    setIsDarkMode(colorScheme === "dark" ? true : false);
  }, [colorScheme]);

  const onAuthClick = () => {
    if (isLogin) {
      Cookie.remove("token");
      navigate("/");
    } else {
      navigate("/auth");
    }
  };

  return (
    <>
      <Menu shadow="xl" px="2" trigger="click" transition="slide-up" withArrow>
        <Menu.Target>
          <Button variant="subtle" size="sm">
            <Avatar size="sm" radius="lg" />
          </Button>
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Label>Vibrain</Menu.Label>
          <Menu.Item leftSection={<Icon icon="tabler:home" />}>
            <Link to="/">Home</Link>
          </Menu.Item>
          <Menu.Item leftSection={<Icon icon="tabler:augmented-reality" />}>
            <Link to="/assistants">Assistants</Link>
          </Menu.Item>
          <Menu.Item
            variant="transparent"
            // size="sm"
            onClick={() => {
              setColorScheme(isDarkMode ? "light" : "dark");
            }}
            leftSection={
              isDarkMode ? (
                <Icon icon="tabler:sun" color="white" />
              ) : (
                <Icon icon="tabler:moon-filled" color="black" />
              )
            }
          >
            {isDarkMode ? "Light" : "Dark"}
          </Menu.Item>
          <Menu.Item
            leftSection={
              isLogin ? (
                <Icon icon="tabler:logout" />
              ) : (
                <Icon icon="tabler:login" />
              )
            }
            onClick={onAuthClick}
          >
            {isLogin ? "Logout" : "Login"}
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    </>
  );
}
