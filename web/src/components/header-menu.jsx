import { Icon } from "@iconify/react";
import { Avatar, Button, Menu, useMantineColorScheme } from "@mantine/core";
import Cookie from "js-cookie";
import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuthContext } from "../libs/auth-context";

export function HeaderMenu() {
  const navigate = useNavigate();
  const { isLogin } = useAuthContext();
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  const onAuthClick = () => {
    if (isLogin) {
      Cookie.remove("token");
      navigate("/");
    } else {
      navigate("/auth");
    }
  };

  const getColorSchemeIcon = () => {
    if (colorScheme === "light") {
      return <Icon icon="tabler:sun-filled" color="dark" />;
    } else if (colorScheme === "dark") {
      return <Icon icon="tabler:moon-filled" color="dark" />;
    } else {
      return <Icon icon="tabler:device-desktop" />;
    }
  };

  const toggleColorScheme = () => {
    const nextScheme =
      colorScheme === "light"
        ? "dark"
        : colorScheme === "dark"
          ? "auto"
          : "light";
    setColorScheme(nextScheme);
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
            onClick={toggleColorScheme}
            leftSection={getColorSchemeIcon()}
          >
            {colorScheme.charAt(0).toUpperCase() + colorScheme.slice(1)}
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
