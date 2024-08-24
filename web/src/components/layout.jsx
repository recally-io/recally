import { AppShell, createTheme, Input, MantineProvider } from "@mantine/core";
import "@mantine/core/styles.css";
import { ModalsProvider } from "@mantine/modals";
import { Notifications } from "@mantine/notifications";
import "@mantine/notifications/styles.css";
import { QueryClientProvider } from "@tanstack/react-query";

import React from "react";
import { queryClient } from "../libs/api";
import useStore from "../libs/store";
import Header from "./header";

const theme = createTheme({
  components: {
    Input: Input.extend({
      styles: {
        input: {
          fontSize: "16px",
        },
      },
    }),
  },
});

export default function Layout({ main, nav = null, header = null }) {
  let haveNav = nav !== null;
  const isSidebarOpen = useStore((state) => state.isSidebarOpen);

  return (
    <QueryClientProvider client={queryClient}>
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <ModalsProvider>
          <Notifications />
          <AppShell
            header={{ height: 40 }}
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
              {header ? header : <Header showNavBurger={haveNav} />}
            </AppShell.Header>
            <AppShell.Navbar p="md">{nav}</AppShell.Navbar>
            <AppShell.Main>{main}</AppShell.Main>
          </AppShell>
        </ModalsProvider>
      </MantineProvider>
    </QueryClientProvider>
  );
}
