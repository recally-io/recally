import { Icon } from "@iconify/react/dist/iconify.js";
import { ActionIcon, Tooltip } from "@mantine/core";
import { useMutation, useQuery } from "@tanstack/react-query";
import { get, post, queryClient } from "../libs/api";
import useStore from "../libs/store";

export function ThreadAddButton() {
  const isLogin = useStore((state) => state.isLogin);
  const assistantId = useStore((state) => state.assistantId);
  const setThreadId = useStore((state) => state.setThreadId);
  const getAssistant = useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}`);
      return res.data;
    },
    enabled: isLogin && !!assistantId,
  });

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
  });

  const addNewThread = async () => {
    await createThread.mutateAsync({
      name: "Thread name",
      description: "Thread description",
      systemPrompt: getAssistant.data.systemPrompt,
      model: getAssistant.data.model,
    });
  };

  return (
    <Tooltip label="New Thread">
      <ActionIcon size="lg" variant="subtle" radius="lg" onClick={addNewThread}>
        <Icon icon="tabler:message-circle-plus" />
      </ActionIcon>
    </Tooltip>
  );
}
