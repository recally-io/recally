import { Icon } from "@iconify/react/dist/iconify.js";
import { CodeHighlight } from "@mantine/code-highlight";
import {
  Accordion,
  Avatar,
  Badge,
  Box,
  Button,
  Collapse,
  Container,
  Divider,
  em,
  Flex,
  Group,
  Image,
  LoadingOverlay,
  Modal,
  Paper,
  ScrollArea,
  SimpleGrid,
  Stack,
  Text,
  useComputedColorScheme,
} from "@mantine/core";
import { useDisclosure, useMediaQuery } from "@mantine/hooks";
import { useEffect, useRef, useState } from "react";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";
import { CopyBtn } from "./copy-button";
import { MarkdownRenderer } from "./markdown-renderer";

export function ThreadChatWindows() {
  const computedColorScheme = useComputedColorScheme("light");
  const desktopSidebarOpen = useStore((state) => state.desktopSidebarOpen);
  const isMobile = useMediaQuery(`(max-width: ${em(750)})`);
  const messageList = useStore((state) => state.threadMessageList);
  const { getThread } = useQueryContext();
  const chatArea = useRef(null);
  const [imageOpened, setImageOpened] = useState(null);
  const [stepsOpened, { toggle: toggleSteps }] = useDisclosure(false);
  const [messageId, setMessageId] = useState(null);

  useEffect(() => {
    chatArea.current.scrollTo({
      top: chatArea.current.scrollHeight,
      behavior: "smooth",
    });
  }, [messageList]);

  const avater = (children) => {
    return (
      <Avatar
        size={{ base: "sm", md: "md" }}
        radius="md"
        color="cyan"
        variant="filled"
      >
        {children}
      </Avatar>
    );
  };

  const IntermediateStep = ({ step, index }) => {
    return (
      <Accordion.Item key={index} value={step.name + "-" + `${index}`}>
        <Accordion.Control
          icon={
            <Badge size="xs" color="pink">
              {step.type}
            </Badge>
          }
        >
          {step.name}
        </Accordion.Control>
        <Accordion.Panel>
          <Paper
            shadow="md"
            p="xs"
            radius="md"
            withBorder
            style={{
              wordBreak: "break-word",
            }}
          >
            <Text size="sm" c="dimmed">
              Input:
            </Text>
            <CodeHighlight
              language="json"
              code={
                typeof step.input === "string"
                  ? JSON.stringify(JSON.parse(step.input), null, 2)
                  : JSON.stringify(step.input, null, 2)
              }
            />
            <Text size="sm" c="dimmed" mt="xs">
              Output:
            </Text>
            <CodeHighlight
              code={
                typeof step.output === "string"
                  ? JSON.stringify(JSON.parse(step.output), null, 4)
                  : JSON.stringify(step.output, null, 4)
              }
              language="json"
            />
          </Paper>
        </Accordion.Panel>
      </Accordion.Item>
    );
  };

  const messagePaper = (message) => {
    const isSender = message.role === "user";
    const bgColor = isSender ? "primary.2" : "secondary.2";
    return (
      <Flex
        justify={isSender ? "flex-end" : "flex-start"}
        align="flex-start"
        direction="row"
        gap="2"
        key={message.id}
      >
        {!isSender && avater(<Icon icon="tabler:brand-android" />)}
        <Flex align={isSender ? "flex-end" : "flex-start"} direction="column">
          <Paper
            shadow="md"
            px="xs"
            radius="lg"
            withBorder
            style={{
              wordBreak: "break-word",
            }}
            bg={computedColorScheme === "dark" ? "dark.6" : bgColor}
            maw={{
              base: "calc(100vw - 80px)",
              sm:
                !isMobile && desktopSidebarOpen
                  ? "calc(100vw - 380px)"
                  : "calc(100vw - 80px)",
            }}
          >
            {message.metadata?.intermediate_steps &&
              message.metadata.intermediate_steps.length > 0 && (
                <Container size="md">
                  <Flex direction="column">
                    <Button
                      variant="subtle"
                      justify="space-between"
                      radius="md"
                      onClick={() => {
                        toggleSteps();
                        setMessageId(message.id);
                      }}
                      rightSection={
                        <Icon
                          icon={
                            stepsOpened
                              ? "tabler:chevron-up"
                              : "tabler:chevron-down"
                          }
                        />
                      }
                    >
                      <Text>Intermediate Steps</Text>
                    </Button>
                    <Collapse in={stepsOpened && message.id == messageId}>
                      <Accordion>
                        {message.metadata.intermediate_steps.map(
                          (step, index) => IntermediateStep({ step, index }),
                        )}
                      </Accordion>
                    </Collapse>
                    <Divider />
                  </Flex>
                </Container>
              )}
            {message.metadata?.images && message.metadata.images.length > 0 && (
              <SimpleGrid cols={{ base: 1, md: 2 }} spacing="xs" mb="xs">
                {message.metadata.images.map((imgurl, index) => (
                  <Box key={index}>
                    <Image
                      src={imgurl}
                      radius="md"
                      alt={`Image ${imgurl.split("/").pop()}`}
                      fit="contain"
                      style={{
                        cursor: "pointer",
                      }}
                      onClick={() => setImageOpened(imgurl)}
                    />
                    <Modal
                      opened={imageOpened === imgurl}
                      onClose={() => setImageOpened(null)}
                      size="xl"
                    >
                      <Image
                        src={imgurl}
                        alt={`Full size image ${index + 1}`}
                        fit="contain"
                      />
                    </Modal>
                  </Box>
                ))}
              </SimpleGrid>
            )}
            <MarkdownRenderer content={message.text} />
          </Paper>
          <Group preventGrowOverflow={false} grow gap="2">
            {!isSender && (
              <Badge
                variant="gradient"
                gradient={{ from: "blue", to: "cyan", deg: 90 }}
                size="xs"
              >
                {message.model}
              </Badge>
            )}
            {CopyBtn({ data: message.text })}
          </Group>
        </Flex>
        {isSender && avater(<Icon icon="tabler:user" />)}
      </Flex>
    );
  };

  return (
    <>
      <LoadingOverlay visible={getThread.isLoading} />
      <ScrollArea
        viewportRef={chatArea}
        type="scroll"
        offsetScrollbars
        scrollbarSize="4"
        scrollbars="y"
        py="xs"
      >
        <Container size="md" p="0">
          <Stack gap="md" align="stretch" justify="flex-start">
            {Array.isArray(messageList) &&
              messageList.map((item) => {
                return messagePaper(item);
              })}
          </Stack>
        </Container>
      </ScrollArea>
    </>
  );
}
