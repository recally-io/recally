import { Icon } from "@iconify/react/dist/iconify.js";
import { Group, Text } from "@mantine/core";
import { Dropzone, IMAGE_MIME_TYPE, PDF_MIME_TYPE } from "@mantine/dropzone";
import { useQueryContext } from "../libs/query-context";
import { fileToDocs } from "../libs/rag.mjs";
import useStore from "../libs/store";

export function UploadButton({ useButton = false }) {
  const {
    getAssistant,
    getThread,
    getPresignedUrlMutation,
    uploadFileMutation,
    postAttachmentMutation,
  } = useQueryContext();

  const addThreadChatImage = useStore((state) => state.addThreadChatImage);

  const handleFilesChange = async (files) => {
    if (!files) return;
    for (const file of files) {
      console.log(`file: ${file.name}, type: ${file.type}`);

      // get presigned url
      const preSignedUrlRes = await getPresignedUrlMutation.mutateAsync({
        assistantId: getAssistant.data.id,
        threadId: getThread.data.id,
        fileName: file.name,
        fileType: file.type,
      });
      // upload file
      const uploadRes = await uploadFileMutation.mutateAsync({
        preSignedURL: preSignedUrlRes.preSignedURL,
        file,
        publicUrl: preSignedUrlRes.publicUrl,
      });
      if (file.type.startsWith("image/")) {
        // add image to chat message
        addThreadChatImage(uploadRes);
      } else if (file.type.endsWith("pdf")) {
        // extract text from pdf
        const docs = await fileToDocs(file, uploadRes);

        await postAttachmentMutation.mutateAsync({
          assistantId: getAssistant.data.id,
          threadId: getThread.data.id,
          type: file.type,
          name: file.name,
          publicUrl: uploadRes,
          docs: docs,
        });
      }
    }
  };

  return (
    <Dropzone
      onDrop={handleFilesChange}
      onReject={(files) => console.log("rejected files", files)}
      maxSize={10 * 1024 ** 2}
      accept={[...IMAGE_MIME_TYPE, ...PDF_MIME_TYPE]}
    >
      <Group
        justify="center"
        gap="xl"
        mih={useButton ? 0 : 100}
        style={
          useButton
            ? {}
            : {
                pointerEvents: "none",
                border: "1px solid lightblue",
                borderRadius: "10px",
                padding: "10px",
              }
        }
      >
        <Dropzone.Accept>
          <Icon icon="tabler:file-upload" />
        </Dropzone.Accept>
        <Dropzone.Reject>
          <Icon icon="tabler:x" />
        </Dropzone.Reject>
        <Dropzone.Idle>
          <Icon icon="tabler:photo" />
        </Dropzone.Idle>
        {!useButton && (
          <div>
            <Text size="xl" inline>
              Add files to knowledge base
            </Text>
            <Text size="sm" c="dimmed" inline mt={7}>
              Add pdf, images, or any files to help AI answer questions
            </Text>
          </div>
        )}
      </Group>
    </Dropzone>
  );
}
