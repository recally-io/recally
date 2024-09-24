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
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";
import { UploadButton } from "./upload-button";

export function ThreadSettingsModal() {
  const params = useParams();
  const threadId = params.threadId;
  const assistantId = params.assistantId;

  const {
    listTools,
    listModels,
    updateThread,
    getAssistant,
    getThread,
    upsertAssistant,
    listAttachmentsByAssistant,
    listAttachmentsByThread,
  } = useQueryContext();

  const [isOpen, toggleThreadIsOpenSettings] = useStore((state) => [
    state.threadIsOpenSettings,
    state.toggleThreadIsOpenSettings,
  ]);

  const threadSettings = useStore((state) => state.threadSettings);
  const setThreadSettings = useStore((state) => state.setThreadSettings);

  const [attachments, setAttachments] = useState([]);

  const form = useForm({
    initialValues: threadSettings,
  });

  useEffect(() => {
    form.setValues(threadSettings);
  }, [threadSettings]);

  useEffect(() => {
    if (assistantId && !getAssistant.isLoading && getAssistant.data) {
      const assistant = getAssistant.data;
      setThreadSettings(assistant);
      form.setValues(assistant);
    }

    if (threadId && !getThread.isLoading && getThread.data) {
      const assistant = getAssistant.data;
      const thread = getThread.data;

      const settings = {
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
      };

      setThreadSettings(settings);
      form.setValues(settings);
    }
  }, [getAssistant.data, getThread.data]);

  useEffect(() => {
    if (listAttachmentsByAssistant.data) {
      setAttachments(listAttachmentsByAssistant.data);
    }
    if (listAttachmentsByThread.data) {
      setAttachments(listAttachmentsByThread.data);
    }
  }, [listAttachmentsByAssistant.data, listAttachmentsByThread.data]);

  return (
    <Modal
      opened={isOpen}
      onClose={toggleThreadIsOpenSettings}
      title="Settings"
    >
      <UploadButton />
      <Stack spacing="sm" py="xs">
        {attachments.map((attachment) => (
          <Group
            key={attachment.id}
            position="apart"
            style={{ alignItems: "center" }}
          >
            <div style={{ flex: 1 }}>
              <a
                href={attachment.url}
                target="_blank"
                rel="noopener noreferrer"
              >
                <strong>{attachment.name}</strong>
              </a>
              {attachment.size > 0 && (
                <div style={{ fontSize: "0.875rem", color: "#6b7280" }}>
                  {attachment.size} KB
                </div>
              )}
            </div>
            {/* <ActionIcon variant="outline" color="red" size="xs">
              <Icon icon="tabler:x"></Icon>
            </ActionIcon> */}
          </Group>
        ))}
      </Stack>
      <Divider my="sm" variant="dashed" />
      <form
        onSubmit={form.onSubmit(async (values) => {
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
          <Textarea
            withAsterisk
            label="Description"
            minRows={3}
            maxRows={5}
            placeholder="Description of the thread"
            key={form.key("description")}
            {...form.getInputProps("description")}
          />
          <Textarea
            withAsterisk
            minRows={3}
            maxRows={5}
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
          <Button type="submit" onClick={toggleThreadIsOpenSettings}>
            Submit
          </Button>
          <Button type="button" onClick={toggleThreadIsOpenSettings}>
            Cancel
          </Button>
        </Group>
      </form>
    </Modal>
  );
}
