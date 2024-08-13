import {
  AppShell,
  Container,
  MantineProvider,
  Text,
  createTheme,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { useDisclosure } from "@mantine/hooks";
import React from "react";
import Header from "./header";

const theme = createTheme({});

export default function Layout({ main, nav = null }) {
  const [opened, { toggle }] = useDisclosure(true);
  let haveNav = nav !== null;

  return (
    <MantineProvider theme={theme} defaultColorScheme="auto">
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
  );
}
