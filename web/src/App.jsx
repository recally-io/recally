import {
  AppShell,
  Center,
  MantineProvider,
  Text,
  createTheme,
} from "@mantine/core";
import "@mantine/core/styles.css";
import { useDisclosure } from "@mantine/hooks";
import Header from "./components/header";
import Sidebar from "./components/sidebar";
import ChatWindowsComponent from "./components/thread";

export default function App() {
  const theme = createTheme({});
  return (
    <>
      <MantineProvider theme={theme} defaultColorScheme="auto">
        <Layout />
      </MantineProvider>
    </>
  );
}

function Layout() {
  const [opened, { toggle }] = useDisclosure(true);

  return (
    <AppShell
      header={{ height: "36px" }}
      navbar={{
        width: "260px",
        breakpoint: "sm",
        collapsed: { mobile: !opened, desktop: !opened },
      }}
      padding="md"
      withBorder={false}
    >
      <AppShell.Header>
        <Header opened={opened} toggle={toggle} />
      </AppShell.Header>

      <AppShell.Navbar
        p="md"
        style={{
          maxWidth: "260px",
        }}
      >
        <Sidebar />
      </AppShell.Navbar>

      <AppShell.Main>
        <ChatWindowsComponent />
      </AppShell.Main>
      <AppShell.Footer>
        <Center>
          <Text align="center" size="xs">
            © 2024 Vibrain Inc.
          </Text>
        </Center>
      </AppShell.Footer>
    </AppShell>
  );
}
