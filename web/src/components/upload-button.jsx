import { Icon } from "@iconify/react/dist/iconify.js";
import { ActionIcon, FileButton } from "@mantine/core";
import { useMutation } from "@tanstack/react-query";
import { useState } from "react";
import { toastError } from "../libs/alert";
import { getPresignedUrl, uploadFile, postAttachment } from "../libs/api";
import { fileToDocs } from "../libs/rag.mjs";
import useStore from "../libs/store";

export function UploadButton() {
  const assistant = useStore((state) => state.assistant);
  const thread = useStore((state) => state.thread);

  const [files, addFile, setFiles] = useStore((state) => [
    state.files,
    state.addFile,
    state.setFiles,
  ]);
  const addThreadChatImage = useStore((state) => state.addThreadChatImage);

  const getPresignedUrlMutation = useMutation({
    mutationFn: getPresignedUrl,
    onSuccess: (data) => {
      console.log("Presigned URL:", data);
    },
    onError: (error) => {
      console.error("Error getting presigned URL:", error);
      toastError("Failed to get upload URL: " + error.message);
    },
  });

  const uploadFileMutation = useMutation({
    mutationFn: uploadFile,
    onSuccess: (data) => {
      console.log("File uploaded:", data);
    },
    onError: (error) => {
      console.error("Error uploading file:", error);
      toastError("Failed to upload file: " + error.message);
    },
  });

  const postAttachmentMutation = useMutation({
    mutationFn: postAttachment,
    onSuccess: (data) => {
      console.log("Attachment posted:", data);
    },
    onError: (error) => {
      console.error("Error posting attachment:", error);
      toastError("Failed to post attachment: " + error.message);
    },
    enabled: assistant.id,
  });

  const handleFilesChange = async (files) => {
    if (!files) return;
    for (const file of files) {
      console.log(`file: ${file.name}, type: ${file.type}`);

      // get presigned url
      const preSignedUrlRes = await getPresignedUrlMutation.mutateAsync({
        assistantId: assistant.id,
        threadId: thread.id,
        fileName: file.name,
        fileType: file.type,
      });
      // upload file
      const uploadRes = await uploadFileMutation.mutateAsync({
        preSignedURL: preSignedUrlRes.data.presigned_url,
        file,
        publicUrl: preSignedUrlRes.data.public_url,
      });
      const publicUrl = uploadRes.data;
      if (file.type.startsWith("image/")) {
        // add image to chat message
        addThreadChatImage(res.data);
      } else if (file.type.endsWith("pdf")) {
        // extract text from pdf
        const docs = await fileToDocs(file, publicUrl);
        console.log("docs", docs);

        await postAttachmentMutation.mutateAsync({
          assistantId: assistant.id,
          threadId: thread.id,
          type: file.type,
          name: file.name,
          publicUrl: publicUrl,
          docs: docs,
        });
      }

      addFile({
        type: file.type,
        name: file.name,
        url: publicUrl,
      });
    }
  };

  return (
    <FileButton
      onChange={handleFilesChange}
      accept="image/*,application/pdf"
      multiple
      disabled={
        getPresignedUrlMutation.isPending || uploadFileMutation.isPending
      }
    >
      {(props) => (
        <ActionIcon {...props} variant="subtle" radius="lg">
          <Icon icon="tabler:file-upload"></Icon>
        </ActionIcon>
      )}
    </FileButton>
  );
}
