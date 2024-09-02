import {
  Button,
  Divider,
  Group,
  Modal,
  MultiSelect,
  NativeSelect,
  Stack,
  TextInput,
  Textarea,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { useEffect } from "react";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";
import { UploadButton } from "./upload-button";

export function ThreadSettingsModal() {
  const { listTools, listModels, updateThread, getAssistant, getThread } =
    useQueryContext();

  const [isOpen, setIsOpen] = useStore((state) => [
    state.threadIsOpenSettings,
    state.setThreadIsOpenSettings,
  ]);

  const form = useForm({
    initialValues: {
      name: "New Thread",
      description: "",
      system_prompt: "",
      temperature: 0.7,
      max_token: 4096,
      model: "",
      metadata: {
        tools: [],
      },
    },
  });

  useEffect(() => {
    if (getAssistant.data) {
      const assistant = getAssistant.data;
      form.setValues({
        name: assistant.name,
        description: assistant.description,
        system_prompt: assistant.system_prompt,
        model: assistant.model,
        metadata: {
          tools: assistant.metadata.tools,
        },
      });
    }
  }, [getAssistant.data]);

  useEffect(() => {
    if (getThread.data) {
      const assistant = getAssistant.data;
      const thread = getThread.data;
      form.setValues({
        name: thread.name ? thread.name : "New Thread",
        description: thread.description
          ? thread.description
          : assistant.description,
        system_prompt: thread.system_prompt
          ? thread.system_prompt
          : assistant.system_prompt,
        temperature: thread.temperature ? thread.temperature : 0.7,
        max_token: thread.max_token ? thread.max_token : 4096,
        model: thread.model ? thread.model : assistant.model,
        metadata: {
          tools: thread.metadata?.tools
            ? thread.metadata.tools
            : assistant.metadata?.tools,
        },
      });
    }
  }, [getThread.data]);

  return (
    <Modal
      opened={isOpen}
      onClose={() => setIsOpen(false)}
      title="Thread Settings"
    >
      <UploadButton />
      <Divider my="sm" variant="dashed" />
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
          <NativeSelect
            label="Model"
            key={form.key("model")}
            {...form.getInputProps("model")}
            data={listModels.data}
          />
          <MultiSelect
            label="Tools"
            key={form.key("metadata.tools")}
            {...form.getInputProps("metadata.tools", {
              type: "checkbox",
            })}
            defaultValue={form.values.metadata.tools}
            data={listTools.data}
            searchable
          />
        </Stack>
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
