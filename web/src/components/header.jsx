import { Icon } from "@iconify/react";
import {
  Avatar,
  Button,
  Flex,
  Menu,
  NavLink,
  useMantineColorScheme,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import Cookie from "js-cookie";
import React, { useEffect } from "react";
import { checkIsLogin } from "../libs/auth";
import useStore from "../libs/store";

const url = new URL(window.location.href);

export default function Header({ showNavBurger }) {
  const [isSidebarOpen, setIsSidebarOpen] = useStore((state) => [
    state.isSidebarOpen,
    state.setIsSidebarOpen,
  ]);
  const isLogin = useStore((state) => state.isLogin);
  const setIsLogin = useStore((state) => state.setIsLogin);
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  const authPage = "/auth.html";

  const [isDarkMode, setIsDarkMode] = useStore((state) => [
    state.isDarkMode,
    state.setIsDarkMode,
  ]);

  useEffect(() => {
    setIsDarkMode(colorScheme === "dark" ? true : false);
  }, [colorScheme]);

  const checkLogin = useQuery({
    queryKey: ["check-login"],
    queryFn: async () => {
      const isLoggedIn = await checkIsLogin();
      return isLoggedIn;
    },
  });

  useEffect(() => {
    // wait until the query is done
    if (checkLogin.isLoading) {
      return;
    }
    if (checkLogin.data) {
      setIsLogin(true);
      console.log("User is logged in");
      if (window.location.pathname === authPage) {
        const redirect = url.searchParams.get("redirect");
        console.log("Redirecting to", redirect || "/");
        window.location.href = redirect || "/";
      }
    } else {
      setIsLogin(false);
      console.log("User is not logged in");
      if (window.location.pathname !== authPage) {
        const redirect = url.pathname + url.search;
        console.log("Redirecting to login page: " + redirect);
        window.location.href = authPage + "?redirect=" + redirect;
      }
    }
  }, [checkLogin.isFetching]);

  const navHome = () => {
    return (
      <Button variant="transparent" px="0">
        <NavLink
          size="sm"
          href="/"
          label="Home"
          leftSection={<Icon icon="tabler:home" />}
        ></NavLink>
      </Button>
    );
  };

  const navAssistants = () => {
    return (
      <Button variant="transparent" px="0">
        <NavLink
          size="sm"
          href="/assistants.html"
          label="Assistants"
          leftSection={<Icon icon="tabler:augmented-reality" />}
        ></NavLink>
      </Button>
    );
  };

  const loginButton = () => {
    if (isLogin) {
      return (
        <Button
          leftSection={<Icon icon="tabler:logout" />}
          variant="transparent"
          onClick={() => {
            Cookie.remove("token");
            window.location.href = "/";
          }}
        >
          Logout
        </Button>
      );
    }
    return (
      <Button
        leftSection={<Icon icon="tabler:login" />}
        variant="transparent"
        onClick={() => {
          window.location.href = authPage;
        }}
      >
        Login
      </Button>
    );
  };

  const themeToggleButton = () => {
    return (
      <Button
        variant="transparent"
        size="sm"
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
      </Button>
    );
  };

  return (
    <>
      <Flex
        direction="row"
        justify="space-between"
        align="center"
        gap="lg"
        // p="0"
      >
        {isSidebarOpen && (
          <Button
            opened={isSidebarOpen}
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
            variant="transparent"
            size="md"
          >
            <Icon icon="tabler:layout-sidebar" />
          </Button>
        )}
        <div></div>
        <Flex visibleFrom="md" gap="1">
          {navHome()}
          {navAssistants()}
          {themeToggleButton()}
          {loginButton()}
        </Flex>
        <Menu
          shadow="xl"
          px="2"
          trigger="click"
          transition="slide-up"
          withArrow
          hiddenFrom="md"
        >
          <Menu.Target>
            <Button variant="transparent" size="sm">
              <Avatar size="sm" radius="lg" />
            </Button>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Label>Vibrain</Menu.Label>
            <Menu.Item>{navHome()}</Menu.Item>
            <Menu.Item>{navAssistants()}</Menu.Item>
            <Menu.Item>{themeToggleButton()}</Menu.Item>
            <Menu.Item>{loginButton()}</Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Flex>
    </>
  );
}
