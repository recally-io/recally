import { Icon } from "@iconify/react";
import { Avatar, Button, Menu, useMantineColorScheme } from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import Cookie from "js-cookie";
import React, { useEffect } from "react";
import { checkIsLogin } from "../libs/auth";
import useStore from "../libs/store";

const url = new URL(window.location.href);

export function HeaderMenu() {
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

  const onAuthClick = () => {
    if (isLogin) {
      Cookie.remove("token");
      window.location.href = "/";
    } else {
      window.location.href = authPage;
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
          <Menu.Item
            leftSection={<Icon icon="tabler:home" />}
            component="a"
            href="/"
            target="_blank"
          >
            Home
          </Menu.Item>
          <Menu.Item
            leftSection={<Icon icon="tabler:augmented-reality" />}
            component="a"
            href="/assistants.html"
            target="_blank"
          >
            Assistants
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
