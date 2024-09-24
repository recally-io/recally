import { Icon } from "@iconify/react/dist/iconify.js";
import {
  ActionIcon,
  Button,
  Card,
  Container,
  Flex,
  Grid,
  Group,
  LoadingOverlay,
  Stack,
  Text,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { modals } from "@mantine/modals";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQueryContext } from "../libs/query-context";
import useStore, { defaultThreadSettings } from "../libs/store";
import { ThreadSettingsModal } from "./thread-settings";

export default function Assistants() {
  const navigate = useNavigate();
  const params = useParams();
  const { listAssistants, deleteAssistant } = useQueryContext();
  const toggleThreadIsOpenSettings = useStore(
    (state) => state.toggleThreadIsOpenSettings,
  );
  const setThreadSettings = useStore((state) => state.setThreadSettings);

  const [filteredAssistants, setFilteredAssistants] = useState([]);
  const [searchValue, setSearchValue] = useState(params.assistantId);

  useEffect(() => {
    if (listAssistants.data) {
      setFilteredAssistants(listAssistants.data);
    }
  }, [listAssistants.data]);

  const setAssistantId = function (assistantId) {
    if (assistantId) {
      navigate(`/assistants/${assistantId}`);
    } else {
      navigate(`/assistants`);
    }
  };

  const deleteConfirmModal = (assistantId) =>
    modals.openConfirmModal({
      title: "Delete assistant",
      children: (
        <Text size="sm">
          Are you sure you want to delete assistant? It will delete all threads
          and messages associated with it.
        </Text>
      ),
      labels: { confirm: "Confirm", cancel: "Cancel" },
      onConfirm: async () => {
        await deleteAssistant.mutateAsync(assistantId);
      },
    });
  return (
    <>
      <ThreadSettingsModal />
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
                setAssistantId(null);
                toggleThreadIsOpenSettings();
                setThreadSettings(defaultThreadSettings);
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
                p="md"
              >
                <Card shadow="sm" padding="xs" radius="md" withBorder>
                  <Title order={3} c="cyan">
                    {assistant.name}
                  </Title>
                  <Text order={4} c="grape">
                    {assistant.description}{" "}
                  </Text>

                  <Group justify="flex-end" gap="xs">
                    <Tooltip label="Chat">
                      <ActionIcon
                        variant="filled"
                        size="xs"
                        onClick={() => {
                          navigate(`/assistants/${assistant.id}/threads`);
                        }}
                      >
                        <Icon icon="tabler:message-2" />
                      </ActionIcon>
                    </Tooltip>
                    <Tooltip label="Edit">
                      <ActionIcon
                        variant="filled"
                        size="xs"
                        onClick={() => {
                          setAssistantId(assistant.id);
                          toggleThreadIsOpenSettings();
                          setThreadSettings(assistant);
                        }}
                      >
                        <Icon icon="tabler:edit" />
                      </ActionIcon>
                    </Tooltip>
                    <Tooltip label="Edit">
                      <ActionIcon
                        variant="filled"
                        size="xs"
                        color="red"
                        onClick={() => {
                          deleteConfirmModal(assistant.id);
                        }}
                      >
                        <Icon icon="tabler:trash" />
                      </ActionIcon>
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
