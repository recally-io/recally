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
  NavLink,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useDisclosure } from "@mantine/hooks";
import { useMutation, useQuery } from "@tanstack/react-query";
import { get, post, queryClient } from "../libs/api";

export default function Assistants() {
  const listAssistants = useQuery({
    queryKey: ["list-asstants"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants");
      return res.data || [];
    },
  });

  const createAssistant = useMutation({
    mutationFn: async (data) => {
      const res = await post("/api/v1/assistants", null, data);
      return res.data;
    },
    onSuccess: async () => {
      queryClient.invalidateQueries("list-asstants");
    },
  });

  const [opened, { open, close }] = useDisclosure(false);
  const form = useForm({
    initialValues: {
      name: "Assistant name",
      description: "Assistant description",
      systemPrompt: "You are a helpful assistant.",
    },

    validate: {},
  });

  if (listAssistants.error) {
    return <div>Error: {listAssistants.error.message}</div>;
  }

  return (
    <>
      <Modal opened={opened} onClose={close} title="Assistant details" centered>
        <form
          onSubmit={form.onSubmit(async (values) => {
            console.log(
              `start createAssistant.mutate: ${JSON.stringify(values)}`,
            );
            await createAssistant.mutateAsync(values);
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
          <Title order={1}>All assistants</Title>
          <Stack justify="space-between" align="center">
            <TextInput
              size="md"
              w="100%"
              radius="md"
              label="Search for assistants"
              description="search assistants by name or description"
              placeholder="Type to search"
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
            {listAssistants.data &&
              listAssistants.data.map((assistant) => (
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

                    <Group mt="xs" justify="flex-end">
                      <Button
                        variant="outline"
                        size="xs"
                        leftSection={<Icon icon="tabler:message-2" />}
                      >
                        <NavLink
                          href={`/threads.html?assistant-id=${assistant.id}`}
                          label="Chat"
                          p="0"
                          size="xs"
                        ></NavLink>
                      </Button>
                      <Button
                        variant="outline"
                        size="xs"
                        leftSection={<Icon icon="tabler:edit" />}
                        onClick={(e) => {
                          form.initialize(assistant);
                          open();
                        }}
                      >
                        Edit
                      </Button>
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
