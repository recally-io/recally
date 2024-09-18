import { Icon } from "@iconify/react/dist/iconify.js";
import { ActionIcon, Tooltip } from "@mantine/core";
import { useNavigate, useParams } from "react-router-dom";
import { useQueryContext } from "../libs/query-context";
import useStore from "../libs/store";

export function ThreadAddButton() {
  const { createThread, getAssistant } = useQueryContext();

  const navigate = useNavigate();
  const params = useParams();
  const assistantId = params.assistantId;

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
    navigate(`/assistants/${assistantId}/threads/${data.id}`, {
      replace: true,
    });
  };

  return (
    <Tooltip label="New Thread">
      <ActionIcon size="md" variant="subtle" radius="lg" onClick={addNewThread}>
        <Icon icon="tabler:message-circle-plus" />
      </ActionIcon>
    </Tooltip>
  );
}

export function ThreadSettingsButton() {
  const toggleThreadIsOpenSettings = useStore(
    (state) => state.toggleThreadIsOpenSettings,
  );
  return (
    <Tooltip label="Thread Settings">
      <ActionIcon
        size="md"
        variant="subtle"
        radius="lg"
        onClick={toggleThreadIsOpenSettings}
      >
        <Icon icon="tabler:settings"></Icon>
      </ActionIcon>
    </Tooltip>
  );
}
