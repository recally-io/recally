import { useMutation, useQuery } from "@tanstack/react-query";
import { createContext, useContext, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { toastError, toastInfo } from "../libs/alert";
import {
  del,
  get,
  getPresignedUrl,
  post,
  postAttachment,
  put,
  queryClient,
  uploadFile,
} from "../libs/api";
import { useAuthContext } from "./auth-context";
import useStore from "./store";
export const QueryContext = createContext();

export function useQueryContext() {
  return useContext(QueryContext);
}

export function QueryContextProvider({ children }) {
  const { isLogin } = useAuthContext();
  const navigate = useNavigate();
  const params = useParams();
  const threadId = params.threadId;
  const assistantId = params.assistantId;

  const setMessageList = useStore((state) => state.setThreadMessageList);
  const addThreadMessage = useStore((state) => state.addThreadMessage);
  const updateLastThreadMessage = useStore(
    (state) => state.updateLastThreadMessage,
  );

  useEffect(() => {
    if (!threadId) {
      setMessageList([]);
    }
  }, [threadId]);

  const listModels = useQuery({
    queryKey: ["list-assistants-models"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants/models");
      return res.data || [];
    },
    enabled: isLogin,
  });

  const listTools = useQuery({
    queryKey: ["list-assistants-tools"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants/tools");
      let data = res.data || [];
      data = data.map((tool) => tool.name);
      return data;
    },
    enabled: isLogin,
  });

  const listAssistants = useQuery({
    queryKey: ["list-asstants"],
    queryFn: async () => {
      const res = await get("/api/v1/assistants");
      return res.data || [];
    },
    enabled: isLogin,
  });

  const getAssistant = useQuery({
    queryKey: ["get-assistant", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}`);
      return res.data;
    },
    enabled: isLogin && !!assistantId,
  });

  const upsertAssistant = useMutation({
    mutationFn: async (data) => {
      if (assistantId) {
        const res = await put(`/api/v1/assistants/${assistantId}`, null, data);
        return res.data;
      } else {
        const res = await post("/api/v1/assistants", null, data);
        return res.data;
      }
    },
    onSuccess: async () => {
      queryClient.invalidateQueries("list-asstants");
    },
    onError: (error) => {
      toastError(
        `Failed to upsert assistant ${assistantId} : ${error.message}`,
      );
    },
  });

  const deleteAssistant = useMutation({
    mutationFn: async (id) => {
      const res = await del(`/api/v1/assistants/${id}`);
      return res.data;
    },
    onSuccess: async () => {
      queryClient.invalidateQueries("list-asstants");
    },
    onError: (error) => {
      toastError(
        `Failed to delete assistant ${assistantId} : ${error.message}`,
      );
    },
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
      queryClient.invalidateQueries({
        queryKey: ["list-threads", assistantId],
      });
      if (threadId) {
        navigate(`/assistants/${assistantId}/threads/${data.id}`, {
          replace: true,
        });
      }
    },
    enabled: isLogin && !!assistantId,
  });

  const generateThreadTitle = useMutation({
    mutationFn: async () => {
      const res = await post(
        `/api/v1/assistants/${assistantId}/threads/${threadId}/generate-title`,
        null,
        {},
      );
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["get-thread", threadId]);
      queryClient.invalidateQueries(["list-threads", assistantId]);
    },
  });

  const sendThreadMessage = useMutation({
    mutationFn: async ({ model, text, images }) => {
      addThreadMessage({
        role: "user",
        text,
        id: Math.random(),
        metadata: { images: images },
      });

      const isNewThread = !threadId;
      let newThreadId = threadId;
      if (isNewThread) {
        newThreadId = crypto.randomUUID();
        const assistant = getAssistant.data;
        await createThread.mutateAsync({
          id: newThreadId,
          name: "New Thread",
          description: assistant.description,
          system_prompt: assistant.systemPrompt,
          model: assistant.model,
          metadata: {
            is_generated_title: false,
            tools: assistant.metadata.tools,
          },
        });
      }
      let payload = {
        role: "user",
        text: text,
        model: model,
      };

      if (images.length > 0) {
        payload["metadata"] = { images: images };
      }

      const uri = `/api/v1/assistants/${assistantId}/threads/${newThreadId}/messages`;
      const fetchSSE = async () => {
        const response = await fetch(uri, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(payload),
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";
        let msg = null;

        while (true) {
          const { value, done } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split("\n");

          buffer = lines.pop();

          for (const line of lines) {
            if (!line.startsWith("data: ")) continue;
            try {
              const event = JSON.parse(line.substring(6));
              if (!msg) {
                msg = event;
                addThreadMessage(msg);
              } else {
                msg.text += event.text;
                updateLastThreadMessage(msg);
              }
            } catch (error) {
              console.error("Error parsing SSE:", error);
            }
          }
        }
      };

      fetchSSE().catch((error) => {
        console.error("SSE fetch failed:", error);
        toastError("Failed to receive message stream: " + error.message);
      });

      if (isNewThread) {
        navigate(`/assistants/${assistantId}/threads/${newThreadId}`, {
          replace: true,
        });
      }
    },
    onError: (error) => {
      toastError("Failed to send message: " + error.message);
    },
  });

  const listThreads = useQuery({
    queryKey: ["list-threads", assistantId],
    queryFn: async () => {
      const res = await get(`/api/v1/assistants/${assistantId}/threads`);
      const data = res.data;
      data.map((item) => {
        item["value"] =
          item["name"] + " - " + item["description"] + " - " + item["id"];
      });
      return data;
    },
    enabled: isLogin && !!assistantId,
  });

  const getThread = useQuery({
    queryKey: ["get-thread", threadId],
    queryFn: async () => {
      const res = await get(
        `/api/v1/assistants/${assistantId}/threads/${threadId}`,
      );
      setMessageList(res.data.messages || []);
      return res.data || {};
    },
    enabled: isLogin && !!threadId && !!assistantId,
  });

  const deleteThread = useMutation({
    mutationFn: async () => {
      await del(`/api/v1/assistants/${assistantId}/threads/${threadId}`);
      console.log("delete thread success");
    },
    onSuccess: () => {
      console.log("onSuccess: delete thread success");
      queryClient.invalidateQueries({
        queryKey: ["list-threads", assistantId],
      });
      setMessageList([]);
      navigate(`/assistants/${assistantId}/threads`);
    },
  });

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
      toastInfo("Thread updated");
    },
    onError: (error) => {
      toastError("Failed to update thread: " + error.message);
    },
  });

  const getPresignedUrlMutation = useMutation({
    mutationFn: getPresignedUrl,
    onError: (error) => {
      toastError("Failed to get upload URL: " + error.message);
    },
  });

  const uploadFileMutation = useMutation({
    mutationFn: uploadFile,
    onSuccess: (data) => {
      toastInfo("File uploaded");
    },
    onError: (error) => {
      toastError("Failed to upload file: " + error.message);
    },
  });

  const postAttachmentMutation = useMutation({
    mutationFn: postAttachment,
    onSuccess: (data) => {
      toastInfo("Attachment added to knowledge base: " + data.name);
    },
    onError: (error) => {
      toastError("Failed to post attachment: " + error.message);
    },
    enabled: isLogin && !!assistantId,
  });

  return (
    <QueryContext.Provider
      value={{
        listModels,
        listTools,

        listAssistants,
        getAssistant,
        upsertAssistant,
        deleteAssistant,

        listThreads,
        getThread,
        createThread,
        updateThread,
        deleteThread,
        generateThreadTitle,

        sendThreadMessage,

        getPresignedUrlMutation,
        uploadFileMutation,
        postAttachmentMutation,
      }}
    >
      {children}
    </QueryContext.Provider>
  );
}
