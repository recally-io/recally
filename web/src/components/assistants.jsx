import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Button,
  Card,
  Container,
  Flex,
  Grid,
  Group,
  LoadingOverlay,
  Modal,
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
import { useMutation, useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toastError } from "../libs/alert";
import { get, post, put, queryClient } from "../libs/api";

export default function Assistants() {
  const [assistantId, setAssistantId] = useState("");
  const [filteredAssistants, setFilteredAssistants] = useState([]);
  const [searchValue, setSearchValue] = useState("");

  const listAssistants = useQuery({
    queryKey: ["list-asstants"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants");
      return res.data || [];
    },
  });

  useEffect(() => {
    if (listAssistants.data) {
      setFilteredAssistants(listAssistants.data);
    }
  }, [listAssistants.data]);

  const listModels = useQuery({
    queryKey: ["list-assistants-models"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants/models");
      return res.data || [];
    },
  });

  const upsertAssistant = useMutation({
    mutationFn: async (data) => {
      if (assistantId) {
        console.log(
          `update assistant ${assistantId}, data: ${JSON.stringify(data)}`,
        );
        const res = await put(`/api/v1/assistants/${assistantId}`, null, data);
        return res.data;
      } else {
        console.log(`create assistant: ${JSON.stringify(data)}`);
        const res = await post("/api/v1/assistants", null, data);
        return res.data;
      }
    },
    onSuccess: async () => {
      queryClient.invalidateQueries("list-asstants");
    },
    onError: (error) => {
      toastError(
        `Failed to upsert assistant ${assistantId} : ${error.message}`,
      );
    },
  });

  const [opened, { open, close }] = useDisclosure(false);
  const form = useForm({
    initialValues: {
      name: "Assistant name",
      description: "Assistant description",
      systemPrompt: "You are a helpful assistant.",
      model: "gpt-4o",
    },

    validate: {},
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
            key={form.key("systemPrompt")}
            {...form.getInputProps("systemPrompt")}
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

          <Group justify="flex-end" mt="md">
            <Button type="summit" onClick={close}>
              Submit
            </Button>
            <Button type="reset" onClick={close}>
              Cancel
            </Button>
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
                    (assistant.name + assistant.description)
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
                      <Button variant="outline" size="xs" w={60}>
                        <a href={`/threads.html?assistant-id=${assistant.id}`}>
                          <Icon icon="tabler:message-2" />
                        </a>
                      </Button>
                    </Tooltip>
                    <Tooltip label="Edit">
                      <Button
                        variant="outline"
                        size="xs"
                        w={60}
                        onClick={(e) => {
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
