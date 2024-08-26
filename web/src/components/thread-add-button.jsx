import { Icon } from "@iconify/react/dist/iconify.js";
import { ActionIcon, Tooltip } from "@mantine/core";
import { useMutation } from "@tanstack/react-query";
import { post, queryClient } from "../libs/api";
import useStore from "../libs/store";

export function ThreadAddButton() {
  const isLogin = useStore((state) => state.isLogin);
  const assistantId = useStore((state) => state.assistantId);
  const setThreadId = useStore((state) => state.setThreadId);
  const assistant = useStore((state) => state.assistant);

  const createThread = useMutation({
    mutationFn: async (data) => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads`,
        null,
        data,
      );
      return res.data;
    },
    onSuccess: (data) => {
      setThreadId(data.id);
      queryClient.invalidateQueries({
        queryKey: ["list-threads", assistantId],
      });
    },
    enabled: isLogin && !!assistantId,
  });

  const addNewThread = async () => {
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
