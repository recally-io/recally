import { AppShell, createTheme, MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { queryClient } from "../libs/api";
import useStore from "../libs/store";
import Header from "./header";

const theme = createTheme({});

export default function Layout({ main, nav = null }) {
  let haveNav = nav !== null;
  const isSidebarOpen = useStore((state) => state.isSidebarOpen);

  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <Notifications />
        <AppShell
          header={{ height: "4dvh" }}
          navbar={{
            width: "300",
            breakpoint: "sm",
            collapsed: {
              mobile: !haveNav || isSidebarOpen,
              desktop: !haveNav || isSidebarOpen,
            },
          }}
          padding="0"
          withBorder={true}
          layout="alt"
        >
          <AppShell.Header>
            <Header showNavBurger={haveNav} />
          </AppShell.Header>
          <AppShell.Navbar p="md">{nav}</AppShell.Navbar>
          <AppShell.Main>{main}</AppShell.Main>
        </AppShell>
      </MantineProvider>
    </QueryClientProvider>
  );
}
