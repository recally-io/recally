import {
  AppShell,
  Container,
  MantineProvider,
  Text,
  createTheme,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { useDisclosure } from "@mantine/hooks";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React, { useEffect } from "react";
import { checkIsLogin } from "../libs/auth";
import useStore from "../libs/store";
import Header from "./header";

const theme = createTheme({});
const queryClient = new QueryClient();

export default function Layout({ main, nav = null }) {
  const [opened, { toggle }] = useDisclosure(true);
  let haveNav = nav !== null;

  const setIsLogin = useStore((state) => state.setIsLogin);
  const authPage = "/auth.html";
  useEffect(() => {
    const checkLoginStatus = async () => {
      const isLoggedIn = await checkIsLogin();
      console.log("Checking login status: ", isLoggedIn);
      if (isLoggedIn) {
        setIsLogin(true);
        console.log("User is logged in");
        if (window.location.pathname === authPage) {
          console.log("Redirecting to home page");
          window.location.href = "/";
        }
      } else {
        setIsLogin(false);
        console.log("User is not logged in");
        if (window.location.pathname !== authPage) {
          console.log("Redirecting to login page");
          window.location.href = "/auth.html";
        }
      }
    };

    checkLoginStatus();
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <Notifications />
        <AppShell
          header={{ height: "36" }}
          footer={{ height: "36" }}
          navbar={{
            width: "300",
            breakpoint: "sm",
            collapsed: {
              mobile: !haveNav || !opened,
              desktop: !haveNav || !opened,
            },
          }}
          padding="md"
          withBorder={false}
        >
          <AppShell.Header>
            <Header opened={opened} toggle={toggle} showNavBurger={haveNav} />
          </AppShell.Header>
          <AppShell.Navbar p="md">{nav}</AppShell.Navbar>
          <AppShell.Main>{main}</AppShell.Main>
          <AppShell.Footer>
            <Container py="sm">
              <Text align="center" size="xs">
                © 2024 Vibrain Inc.
              </Text>
            </Container>
          </AppShell.Footer>
        </AppShell>
      </MantineProvider>
    </QueryClientProvider>
  );
}
