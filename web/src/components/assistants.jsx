import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Anchor,
  Button,
  Card,
  Container,
  Flex,
  Grid,
  Group,
  LoadingOverlay,
  Modal,
  MultiSelect,
  NativeSelect,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { modals } from "@mantine/modals";
import { useEffect, useState } from "react";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";
const url = new URL(window.location.href);

export default function Assistants() {
  const {
    listAssistants,
    listModels,
    listTools,
    upsertAssistant,
    deleteAssistant,
  } = useQueryContext();
  const assistantId = useStore((state) => state.assistantId);
  const setAssistantId = useStore((state) => state.setAssistantId);
  const [filteredAssistants, setFilteredAssistants] = useState([]);
  const [searchValue, setSearchValue] = useState(url.searchParams.get("id"));

  useEffect(() => {
    if (listAssistants.data) {
      setFilteredAssistants(listAssistants.data);
    }
  }, [listAssistants.data]);

  const [opened, { open, close }] = useDisclosure(false);
  const form = useForm({
    initialValues: {
      name: "Assistant name",
      description: "Assistant description",
      system_prompt: "You are a helpful assistant.",
      model: "gpt-4o",
      metadata: {
        tools: [],
      },
    },

    validate: {},
  });

  const deleteConfirmModal = () =>
    modals.openConfirmModal({
      title: "Delete assistant",
      children: (
        <Text size="sm">
          Are you sure you want to delete assistant? It will delete all threads
          and messages associated with it.
        </Text>
      ),
      labels: { confirm: "Confirm", cancel: "Cancel" },
      onCancel: close,
      onConfirm: async () => {
        await deleteAssistant.mutateAsync(assistantId);
        close();
      },
    });
  return (
    <>
      <Modal opened={opened} onClose={close} title="Assistant details" centered>
        <form
          onSubmit={form.onSubmit(async (values) => {
            console.log(
              `start createAssistant.mutate: ${JSON.stringify(values)}`,
            );
            await upsertAssistant.mutateAsync(values);
          })}
        >
          <TextInput
            withAsterisk
            label="Name"
            placeholder="your@email.com"
            key={form.key("name")}
            {...form.getInputProps("name")}
          />
          <TextInput
            withAsterisk
            label="Description"
            placeholder="your@email.com"
            key={form.key("description")}
            {...form.getInputProps("description")}
          />
          <Textarea
            withAsterisk
            label="SystemPrompt"
            placeholder="your@email.com"
            key={form.key("system_prompt")}
            {...form.getInputProps("system_prompt")}
          />
          <NativeSelect
            label="Model"
            key={form.key("model")}
            {...form.getInputProps("model")}
            onChange={(e) => {
              form.setFieldValue("model", e.target.value);
            }}
            data={listModels.data}
          />
          <MultiSelect
            label="Tools"
            key={form.key("metadata.tools")}
            defaultValue={form.values.metadata.tools}
            {...form.getInputProps("metadata.tools", {
              type: "checkbox",
            })}
            data={listTools.data}
            searchable
          />

          <Group justify="space-between" mt="md">
            <Button
              type="button"
              onClick={deleteConfirmModal}
              color="red"
              variant="filled"
            >
              Delete
            </Button>
            <Group>
              <Button type="summit" onClick={close}>
                Submit
              </Button>
              <Button type="reset" onClick={close}>
                Cancel
              </Button>
            </Group>
          </Group>
        </form>
      </Modal>
      <Container size="xl" mih="100vh" py="md">
        <Flex justify="center" align="center" direction="column" gap="lg">
          <Title order={1}>Assistants Hub</Title>
          <Stack justify="space-between" align="center">
            <TextInput
              size="md"
              w="100%"
              radius="md"
              // label="Search for assistants"
              description="search assistants by name or description"
              placeholder="Type to search"
              value={searchValue}
              onChange={(e) => {
                setSearchValue(e.currentTarget.value);
                setFilteredAssistants(
                  listAssistants.data.filter((assistant) =>
                    (assistant.name + assistant.description + assistant.id)
                      .toLowerCase()
                      .includes(e.currentTarget.value.toLowerCase()),
                  ),
                );
              }}
            />
            <Button
              w="100%"
              onClick={() => {
                form.reset();
                open();
              }}
            >
              Add assistant
            </Button>
          </Stack>
          <LoadingOverlay visible={listAssistants.isLoading} />
          <Grid gutter="lg" justify="center" align="center" w="100%">
            {filteredAssistants.map((assistant) => (
              <Grid.Col
                key={assistant.id}
                span={{ base: 12, md: 6, lg: 3 }}
                m="md"
              >
                <Card shadow="sm" padding="lg" radius="md" withBorder>
                  <Title order={3} c="cyan">
                    {assistant.name}
                  </Title>
                  <Text order={4} c="grape">
                    {assistant.description}{" "}
                  </Text>

                  <Group mt="xs" mb="1" justify="flex-end">
                    <Tooltip label="Chat">
                      <Anchor
                        href={`/threads.html?assistant-id=${assistant.id}`}
                      >
                        <Button variant="outline" size="xs" w={60}>
                          <Icon icon="tabler:message-2" />
                        </Button>
                      </Anchor>
                    </Tooltip>
                    <Tooltip label="Edit">
                      <Button
                        variant="outline"
                        size="xs"
                        w={60}
                        onClick={() => {
                          setAssistantId(assistant.id);
                          form.setValues(assistant);
                          open();
                        }}
                      >
                        <Icon icon="tabler:edit" />
                      </Button>
                    </Tooltip>
                  </Group>
                </Card>
              </Grid.Col>
            ))}
          </Grid>
        </Flex>
      </Container>
    </>
  );
}
