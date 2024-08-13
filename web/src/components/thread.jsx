import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Avatar,
  Button,
  Container,
  FileButton,
  Flex,
  Group,
  Menu,
  Modal,
  Paper,
  ScrollArea,
  Slider,
  Stack,
  Text,
  TextInput,
  Tooltip,
  useComputedColorScheme,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import Markdown from "react-markdown";
import avatarImgUrl from "../assets/avatar-1.png";

const content = `
# react-markdown
React component to render markdown.
 
## Feature highlights

* [x] **[safe][section-security] by default**
  (no \'dangerouslySetInnerHTML\' or XSS attacks)
* [x] **[components][section-components]**
  (pass your own component to use instead of \'<h2>\' for \'## hi\')
* [x] **[plugins][section-plugins]**
  (many plugins you can pick and choose from)
* [x] **[compliant][section-syntax]**
  (100% to CommonMark, 100% to GFM with a plugin)

## Contents

* [What is this?](#what-is-this)
* [When should I use this?](#when-should-i-use-this)
* [Install](#install)
* [Use](#use)

## What is this?

This package is a [React][] component that can be given a string of markdown
that it’ll safely render to React elements.
You can pass plugins to change how markdown is transformed and pass components
that will be used instead of normal HTML elements.

* to learn markdown, see this [cheatsheet and tutorial][commonmark-help]
* to try out \'react-markdown\', see [our demo][demo]

`;

export default function ChatWindowsComponent() {
  const [settingsOpened, { open: openSettings, close: closeSettings }] =
    useDisclosure(false);
  const colorScheme = useComputedColorScheme("light");
  const settingsForm = useForm({
    initialValues: {
      temperature: 0.7,
      maxToken: 4096,
    },
  });

  let url = new URL(window.location.href);
  let params = new URLSearchParams(url.search);
  let threadId = params.get("threadId");
  const messageS = () => {
    return (
      <Flex justify="flex-end" align="flex-start" direction="row" gap="sm">
        <Paper
          shadow="sm"
          p="md"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "blue.2"}
        >
          <Markdown>{content}</Markdown>
        </Paper>
        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
      </Flex>
    );
  };

  const messageR = () => {
    return (
      <Flex justify="flex-start" direction="row" gap="sm">
        <Avatar size="sm" radius="lg" src={avatarImgUrl} />
        <Paper
          shadow="sm"
          p="md"
          maw="90%"
          radius="lg"
          bg={colorScheme === "dark" ? "" : "green.2"}
        >
          <Markdown>{content}</Markdown>
        </Paper>
      </Flex>
    );
  };

  const menu = () => {
    return (
      <Menu shadow="md" position="top" withArrow>
        <Menu.Target>
          <Button size="compact-lg" variant="subtle" radius="lg">
            <Icon icon="tabler:plus"></Icon>
          </Button>
        </Menu.Target>
        <Menu.Dropdown>
          <Menu.Item>
            <Tooltip label="thread settings">
              <Button
                variant="transparent"
                size="sm"
                onClick={openSettings}
                leftSection={<Icon icon="tabler:settings"></Icon>}
              >
                Settings
              </Button>
            </Tooltip>
          </Menu.Item>
          <Menu.Item>
            <FileButton
              size="sm"
              variant="transparent"
              multiple
              leftSection={<Icon icon="tabler:upload"></Icon>}
            >
              {(props) => <Button {...props}>Upload image</Button>}
            </FileButton>
          </Menu.Item>
        </Menu.Dropdown>
      </Menu>
    );
  };

  return (
    <>
      <Container size="xl">
        <Flex direction="column" justify="space-between" h="89vh">
          <ScrollArea
            style={{
              flex: 1,
            }}
          >
            <Stack spacing="md" py="lg">
              {messageS()}
              {messageR()}
            </Stack>
          </ScrollArea>
          <Container
            w="100%"
            style={{
              position: "sticky",
              bottom: 0,
            }}
          >
            <TextInput
              placeholder="Send a message"
              variant="filled"
              radius="lg"
              leftSection={menu()}
              leftSectionWidth={42}
            ></TextInput>
          </Container>
        </Flex>

        {/* settings modal */}
        <Modal
          opened={settingsOpened}
          onClose={closeSettings}
          title="Advance Settings"
        >
          <form
            onSubmit={settingsForm.onSubmit((values) => console.log(values))}
            mode=""
          >
            <Stack spacing="md">
              <Stack spacing="xs">
                <Text size="sm">Temperature</Text>
                <Slider
                  min={0}
                  max={1}
                  step={0.1}
                  key={settingsForm.key("temperature")}
                  {...settingsForm.getInputProps("temperature")}
                  labelAlwaysOn
                />
              </Stack>
              <Stack spacing="xs">
                <Text size="sm">Max Tokens</Text>
                <Slider
                  min={0}
                  max={4096}
                  step={1}
                  key={settingsForm.key("maxToken")}
                  {...settingsForm.getInputProps("maxToken")}
                  labelAlwaysOn
                />
              </Stack>
            </Stack>
            <Group justify="flex-end" mt="md">
              <Button type="submit">Submit</Button>
              <Button type="reset">Reset</Button>
            </Group>
          </form>
        </Modal>
      </Container>
    </>
  );
}
