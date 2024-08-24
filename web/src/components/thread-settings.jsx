import { Icon } from "@iconify/react/dist/iconify.js";
import {
  Button,
  FileButton,
  Group,
  Modal,
  NativeSelect,
  Slider,
  Stack,
  Text,
} from "@mantine/core";
import { useForm } from "@mantine/form";

import { useEffect } from "react";
import useStore from "../libs/store";

export function ThreadSettingsModal() {
  const [isOpen, setIsOpen] = useStore((state) => [
    state.threadIsOpenSettings,
    state.setThreadIsOpenSettings,
  ]);
  const models = useStore((state) => state.threadModels);
  const threadSettings = useStore((state) => state.threadSettings);

  const settingsForm = useForm({
    initialValues: threadSettings,
  });
  useEffect(() => {
    settingsForm.setValues(threadSettings);
  }, [threadSettings]);

  return (
    <Modal
      opened={isOpen}
      onClose={() => setIsOpen(false)}
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
          <NativeSelect
            label="Model"
            key={settingsForm.key("model")}
            {...settingsForm.getInputProps("model")}
            onChange={(e) => {
              settingsForm.setFieldValue("model", e.target.value);
            }}
            data={models}
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
