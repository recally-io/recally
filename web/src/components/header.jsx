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

const url = new URL(window.location.href);

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
          setColorScheme(colorScheme === "dark" ? "light" : "dark");
        }}
        leftSection={
          colorScheme === "dark" ? (
            <Icon icon="tabler:sun" color="white" />
          ) : (
            <Icon icon="tabler:moon-filled" color="black" />
          )
        }
      >
        {colorScheme === "dark" ? "Light" : "Dark"}
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
        bg={colorScheme === "dark" ? "gray.8" : "gray.3"}
      >
        {showNavBurger && (
          <Burger
            opened={opened}
            onClick={toggle}
            size="sm"
            pl="lg"
            aria-label="Toggle navigation"
          />
        )}
        {!showNavBurger && <div></div>}
        <Flex visibleFrom="md" gap="1">
          {navHome()}
          {navAssistants()}
          {themeToggleButton()}
          {loginButton()}
        </Flex>
        <Menu
          shadow="xl"
          trigger="click"
          transition="slide-up"
          withArrow
          hiddenFrom="md"
        >
          <Menu.Target>
            <Button variant="transparent" size="sm">
              <Avatar size="sm" radius="lg" src={avatarImgUrl} />
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
