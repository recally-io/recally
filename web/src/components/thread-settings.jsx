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
  const {
    listTools,
    listModels,
    updateThread,
    getAssistant,
    getThread,
    upsertAssistant,
  } = useQueryContext();

  const [isOpen, setIsOpen] = useStore((state) => [
    state.threadIsOpenSettings,
    state.setThreadIsOpenSettings,
  ]);

  const assistantId = useStore((state) => state.assistantId);
  const threadId = useStore((state) => state.threadId);

  const getInitialValues = () => {
    if (threadId) {
      return {
        name: "New Thread",
        description: "",
        system_prompt: "",
        temperature: 0.7,
        max_token: 4096,
        model: "",
        metadata: {
          tools: [],
        },
      };
    }
    return {
      name: "Assistant name",
      description: "Assistant description",
      system_prompt: "You are a helpful assistant.",
      model: "gpt-4o",
      metadata: {
        tools: [],
      },
    };
  };

  const form = useForm({
    initialValues: getInitialValues(),
  });

  useEffect(() => {
    if (assistantId && getAssistant.data) {
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
    } else {
      form.reset();
    }
  }, [assistantId]);

  useEffect(() => {
    if (threadId && getThread && getThread.data) {
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
  }, [getThread]);

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
          console.log(`assistantId: ${assistantId}, threadId: ${threadId}`);
          console.log(values);
          if (threadId) {
            await updateThread.mutateAsync(values);
          } else {
            await upsertAssistant.mutateAsync(values);
          }
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
