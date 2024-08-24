import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Button,
  FileButton,
  Group,
  Modal,
  MultiSelect,
  NativeSelect,
  Slider,
  Stack,
  Text,
  TextInput,
  Textarea,
} from "@mantine/core";
import { useForm } from "@mantine/form";

import { useMutation } from "@tanstack/react-query";
import { useEffect } from "react";
import { put, queryClient } from "../libs/api";
import useStore from "../libs/store";

export function ThreadSettingsModal() {
  const [isOpen, setIsOpen] = useStore((state) => [
    state.threadIsOpenSettings,
    state.setThreadIsOpenSettings,
  ]);
  const threadId = useStore((state) => state.threadId);
  const assistantId = useStore((state) => state.assistantId);
  const models = useStore((state) => state.threadModels);
  const tools = useStore((state) => state.threadTools);
  const threadSettings = useStore((state) => state.threadSettings);

  const form = useForm({
    initialValues: threadSettings,
  });
  useEffect(() => {
    form.setValues(threadSettings);
  }, [threadSettings]);

  const updateThread = useMutation({
    mutationFn: async (data) => {
      const res = await put(
        `/api/v1/assistants/${assistantId}/threads/${threadId}`,
        null,
        data,
      );
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["get-thread", threadId]);
    },
  });

  return (
    <Modal
      opened={isOpen}
      onClose={() => setIsOpen(false)}
      title="Thread Settings"
    >
      <form
        onSubmit={form.onSubmit(async (values) => {
          console.log(values);
          await updateThread.mutateAsync(values);
        })}
        mode=""
      >
        <Stack spacing="md">
          <TextInput
            withAsterisk
            label="Name"
            placeholder="Name of the thread"
            key={form.key("name")}
            {...form.getInputProps("name")}
          />
          <TextInput
            withAsterisk
            label="Description"
            placeholder="Description of the thread"
            key={form.key("description")}
            {...form.getInputProps("description")}
          />
          <Textarea
            withAsterisk
            label="SystemPrompt"
            placeholder="System prompt of the thread"
            key={form.key("system_prompt")}
            {...form.getInputProps("system_prompt")}
          />
          <Stack spacing="xs">
            <Text size="sm">Temperature</Text>
            <Slider
              min={0}
              max={1}
              step={0.1}
              key={form.key("temperature")}
              {...form.getInputProps("temperature")}
              labelAlwaysOn
            />
          </Stack>
          <Stack spacing="xs">
            <Text size="sm">Max Tokens</Text>
            <Slider
              min={0}
              max={4096}
              step={1}
              key={form.key("max_token")}
              {...form.getInputProps("max_token")}
              labelAlwaysOn
            />
          </Stack>
          <NativeSelect
            label="Model"
            key={form.key("model")}
            {...form.getInputProps("model")}
            data={models}
          />
          <MultiSelect
            label="Tools"
            key={form.key("metadata.tools")}
            {...form.getInputProps("metadata.tools", {
              type: "checkbox",
            })}
            defaultValue={threadSettings.metadata.tools}
            data={tools}
            searchable
          />
        </Stack>
        <FileButton
          size="sm"
          variant="transparent"
          multiple
          leftSection={<Icon icon="tabler:upload"></Icon>}
        >
          {(props) => <Button {...props}>Upload image</Button>}
        </FileButton>
        <Group justify="flex-end" mt="md">
          <Button type="submit" onClick={() => setIsOpen(false)}>
            Submit
          </Button>
          <Button type="button" onClick={() => setIsOpen(false)}>
            Cancel
          </Button>
        </Group>
      </form>
    </Modal>
  );
}
