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
import { useQuery } from "@tanstack/react-query";
import Cookie from "js-cookie";
import React, { useEffect } from "react";
import avatarImgUrl from "../assets/avatar-1.png";
import { checkIsLogin } from "../libs/auth";
import useStore from "../libs/store";

export default function Header({ opened, toggle, showNavBurger }) {
  const isLogin = useStore((state) => state.isLogin);
  const setIsLogin = useStore((state) => state.setIsLogin);
  const { colorScheme, setColorScheme } = useMantineColorScheme();
  const authPage = "/auth.html";
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
      console.log(`current path: ${window.location}`);
      setIsLogin(true);
      console.log("User is logged in");
      if (window.location.pathname === authPage) {
        const url = new URL(window.location.href);
        const redirect = url.searchParams.get("redirect");
        console.log("Redirecting to", redirect || "/");
        window.location.href = redirect || "/";
      }
    } else {
      setIsLogin(false);
      console.log("User is not logged in");
      if (window.location.pathname !== authPage) {
        console.log("Redirecting to login page");
        window.location.href = authPage + "?redirect=" + window.location;
      }
    }
  }, [checkLogin.isFetching]);

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
                    // setIsLogin(false);
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
