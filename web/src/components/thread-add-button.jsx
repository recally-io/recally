import { Icon } from "@iconify/react/dist/iconify.js";
import { ActionIcon, Tooltip } from "@mantine/core";
import { useQueryContext } from "../libs/query-context";

export function ThreadAddButton() {
  const { createThread, getAssistant } = useQueryContext();

  const addNewThread = async () => {
    const assistant = getAssistant.data;
    const data = {
      id: crypto.randomUUID(),
      name: "New Thread",
      description: assistant.description,
      system_prompt: assistant.system_prompt,
      model: assistant.model,
      metadata: {
        is_generated_title: false,
        tools: assistant.metadata.tools,
      },
    };
    await createThread.mutateAsync(data);
  };

  return (
    <Tooltip label="New Thread">
      <ActionIcon size="lg" variant="subtle" radius="lg" onClick={addNewThread}>
        <Icon icon="tabler:message-circle-plus" />
      </ActionIcon>
    </Tooltip>
  );
}
