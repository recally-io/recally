import { Icon } from "@iconify/react";
import {
  ActionIcon,
  Box,
  Button,
  Card,
  Container,
  CopyButton,
  Divider,
  Group,
  LoadingOverlay,
  Menu,
  Modal,
  Paper,
  rem,
  Spoiler,
  Stack,
  Text,
  TextInput,
  Title,
  Tooltip,
} from "@mantine/core";
import { useForm } from "@mantine/form";
import { notifications } from "@mantine/notifications";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { MarkdownRenderer } from "../components/markdown-renderer";
import { del, get, post, put } from "../libs/api";
import { useAuthContext } from "../libs/auth-context";

const RefreshMenu = ({ bookmark, mutation, size = "md", iconSize = 18 }) => (
  <Menu position="bottom-end" withArrow>
    <Menu.Target>
      <ActionIcon
        variant="subtle"
        size={size}
        loading={mutation.isLoading}
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
        }}
      >
        <Icon icon="tabler:refresh" width={iconSize} />
      </ActionIcon>
    </Menu.Target>
    <Menu.Dropdown>
      <Menu.Label>Refresh Options</Menu.Label>
      <Menu.Item
        leftSection={<Icon icon="tabler:refresh" width={14} height={14} />}
        onClick={(e) => {
          // e.stopPropagation();
          mutation.mutate({ id: bookmark.id, regenerateSummary: true });
        }}
      >
        Regenerate Summary
      </Menu.Item>
      <Menu.Divider />
      <Menu.Label>Refetch Content Using</Menu.Label>
      {["http", "jina", "browser"].map((fetcher) => (
        <Menu.Item
          key={fetcher}
          leftSection={
            <Icon
              icon={
                fetcher === "http"
                  ? "tabler:world"
                  : fetcher === "jina"
                    ? "tabler:api"
                    : "tabler:browser"
              }
              width={14}
              height={14}
            />
          }
          onClick={(e) => {
            // e.stopPropagation();
            mutation.mutate({ id: bookmark.id, fetcher });
          }}
        >
          {fetcher.charAt(0).toUpperCase() + fetcher.slice(1)} Fetcher
        </Menu.Item>
      ))}
    </Menu.Dropdown>
  </Menu>
);

const BookmarkCard = ({
  bookmark,
  onEdit,
  onDelete,
  onSelect,
  refreshMutation,
}) => (
  <Card
    key={bookmark.id}
    withBorder
    padding="md"
    onClick={onSelect}
    style={{ cursor: "pointer" }}
  >
    <Group justify="space-between" align="flex-start" wrap="nowrap">
      <Group gap="md" wrap="nowrap" style={{ flex: 1 }}>
        <Box>
          <img
            src={`https://www.google.com/s2/favicons?domain=${
              new URL(bookmark.url).hostname
            }&sz=32`}
            alt=""
            style={{ width: 32, height: 32, borderRadius: 4 }}
          />
        </Box>
        <Stack gap="xs" style={{ flex: 1 }}>
          <Group justify="space-between" wrap="nowrap">
            <Text size="lg" fw={600} style={{ wordBreak: "break-word" }}>
              {bookmark.title || bookmark.url}
            </Text>
            <Text size="xs" c="dimmed" style={{ whiteSpace: "nowrap" }}>
              {new Date(bookmark.createdAt).toLocaleDateString()}
            </Text>
          </Group>
          {bookmark.metadata?.description && (
            <Text size="sm" c="dimmed" lineClamp={2}>
              {bookmark.metadata.description}
            </Text>
          )}
          <Group gap="xs">
            <Text
              size="sm"
              c="blue"
              component="a"
              href={bookmark.url}
              target="_blank"
              style={{
                textDecoration: "none",
                "&:hover": { textDecoration: "underline" },
              }}
            >
              {new URL(bookmark.url).hostname}
              <Icon
                icon="tabler:external-link"
                width={14}
                height={14}
                style={{ display: "inline", marginLeft: 4 }}
              />
            </Text>
          </Group>
        </Stack>
      </Group>
      <Group gap="xs">
        <ActionIcon
          variant="subtle"
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            onEdit(bookmark);
          }}
        >
          <Icon icon="tabler:edit" width={20} height={20} />
        </ActionIcon>
        <ActionIcon
          variant="subtle"
          color="red"
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            onDelete(bookmark.id);
          }}
        >
          <Icon icon="tabler:trash" width={20} height={20} />
        </ActionIcon>
        <RefreshMenu
          bookmark={bookmark}
          mutation={refreshMutation}
          size="md"
          iconSize={20}
        />
      </Group>
    </Group>
  </Card>
);

const BookmarkDetailModal = ({
  bookmark,
  onClose,
  onEdit,
  onDelete,
  refreshMutation,
}) => (
  <Modal
    opened={!!bookmark}
    onClose={onClose}
    size="100%"
    fullScreen
    styles={{
      header: {
        marginBottom: 0,
        padding: rem(16),
        background: "var(--mantine-color-body)",
        borderBottom: "1px solid var(--mantine-color-default-border)",
        position: "sticky",
        top: 0,
        zIndex: 100,
      },
      body: {
        padding: 0,
      },
      content: {
        background: "var(--mantine-color-gray-0)",
      },
    }}
    withCloseButton={false}
  >
    <Container size="xl" py="xl">
      <Stack gap="lg">
        {/* Quick Actions */}
        <Group gap="xs" wrap="nowrap" justify="flex-end">
          <CopyButton value={bookmark?.url || ""}>
            {({ copied, copy }) => (
              <Tooltip label={copied ? "Copied" : "Copy link"}>
                <ActionIcon
                  variant="subtle"
                  color={copied ? "teal" : "gray"}
                  onClick={copy}
                  size="md"
                >
                  <Icon
                    icon={copied ? "tabler:check" : "tabler:copy"}
                    width={18}
                  />
                </ActionIcon>
              </Tooltip>
            )}
          </CopyButton>

          <Menu position="bottom-end" withArrow>
            <Menu.Target>
              <Tooltip label="Share">
                <ActionIcon variant="subtle" size="md">
                  <Icon icon="tabler:share" width={18} />
                </ActionIcon>
              </Tooltip>
            </Menu.Target>
            <Menu.Dropdown>
              <Menu.Item
                leftSection={<Icon icon="tabler:brand-twitter" width={16} />}
              >
                Twitter
              </Menu.Item>
              <Menu.Item
                leftSection={<Icon icon="tabler:brand-linkedin" width={16} />}
              >
                LinkedIn
              </Menu.Item>
              <Menu.Item leftSection={<Icon icon="tabler:mail" width={16} />}>
                Email
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>

          <Tooltip label="Edit">
            <ActionIcon
              variant="subtle"
              size="md"
              onClick={() => onEdit(bookmark)}
            >
              <Icon icon="tabler:edit" width={18} />
            </ActionIcon>
          </Tooltip>

          <Tooltip label="Delete">
            <ActionIcon
              variant="subtle"
              color="red"
              size="md"
              onClick={() => {
                onDelete(bookmark.id);
                onClose();
              }}
            >
              <Icon icon="tabler:trash" width={18} />
            </ActionIcon>
          </Tooltip>

          <Tooltip label="Open in new tab">
            <ActionIcon
              component="a"
              href={bookmark?.url}
              target="_blank"
              variant="subtle"
              size="md"
            >
              <Icon icon="tabler:external-link" width={18} />
            </ActionIcon>
          </Tooltip>

          <RefreshMenu bookmark={bookmark} mutation={refreshMutation} />

          <Divider orientation="vertical" />
          <ActionIcon variant="subtle" onClick={onClose} size="md">
            <Icon icon="tabler:x" width={18} />
          </ActionIcon>
        </Group>

        {/* Content */}
        <Paper withBorder p="md">
          <Stack gap="md">
            <Group wrap="nowrap" justify="space-between">
              <Text size="xl" fw={700}>
                {bookmark?.title}
              </Text>
              <Text size="sm" c="dimmed">
                {new Date(bookmark?.createdAt).toLocaleDateString()}
              </Text>
            </Group>

            <Group gap="xs">
              <Text
                size="sm"
                c="blue"
                component="a"
                href={bookmark?.url}
                target="_blank"
              >
                {bookmark?.url}
              </Text>
            </Group>

            {bookmark?.summary && (
              <Spoiler maxHeight={120} showLabel="Show more" hideLabel="Hide">
                <Text size="sm">{bookmark?.summary}</Text>
              </Spoiler>
            )}

            {bookmark?.content && (
              <Box mt="md">
                <MarkdownRenderer content={bookmark?.content} />
              </Box>
            )}
          </Stack>
        </Paper>
      </Stack>
    </Container>
  </Modal>
);

export default function Bookmarks() {
  const { checkIsLogin } = useAuthContext();
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [editingBookmark, setEditingBookmark] = useState(null);
  const [selectedBookmark, setSelectedBookmark] = useState(null);
  const queryClient = useQueryClient();

  const { data: bookmarks, isLoading } = useQuery({
    queryKey: ["bookmarks"],
    queryFn: async () => {
      console.log("fetching bookmarks");
      const res = await get("/api/v1/bookmarks", {
        limit: 100,
        offset: 0,
      });
      return res.data;
    },
    enabled: !!checkIsLogin.data,
  });

  const createMutation = useMutation({
    mutationFn: async (newBookmark) =>
      await post("/api/v1/bookmarks", null, newBookmark),
    onSuccess: () => {
      queryClient.invalidateQueries(["bookmarks"]);
      setCreateModalOpen(false);
      notifications.show({
        title: "Success",
        message: "Bookmark created successfully",
        color: "green",
      });
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({ id, data }) =>
      await put(`/api/v1/bookmarks/${id}`, data),
    onSuccess: () => {
      queryClient.invalidateQueries(["bookmarks"]);
      setEditingBookmark(null);
      notifications.show({
        title: "Success",
        message: "Bookmark updated successfully",
        color: "green",
      });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id) => await del(`/api/v1/bookmarks/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries(["bookmarks"]);
      notifications.show({
        title: "Success",
        message: "Bookmark deleted successfully",
        color: "green",
      });
    },
  });

  const refreshMutation = useMutation({
    mutationFn: async ({ id, fetcher, regenerateSummary }) =>
      await post(`/api/v1/bookmarks/${id}/refresh`, null, {
        fetcher,
        regenerate_summary: regenerateSummary,
      }),
    onSuccess: () => {
      queryClient.invalidateQueries(["bookmarks"]);
      notifications.show({
        title: "Success",
        message: "Bookmark refreshed successfully",
        color: "green",
      });
    },
  });

  const form = useForm({
    initialValues: {
      url: "",
      metadata: {},
    },
    validate: {
      url: (value) => {
        try {
          new URL(value);
          return null;
        } catch {
          return "Please enter a valid URL";
        }
      },
    },
  });

  const handleSubmit = (values) => {
    if (editingBookmark) {
      updateMutation.mutate({ id: editingBookmark.id, data: values });
    } else {
      createMutation.mutate(values);
    }
    form.reset();
  };

  const handleEdit = (bookmark) => {
    form.setValues(bookmark);
    setEditingBookmark(bookmark);
    setSelectedBookmark(null);
  };

  const CreateEditModal = ({ opened, onClose }) => (
    <Modal
      opened={opened}
      onClose={onClose}
      title={editingBookmark ? "Edit Bookmark" : "Add New Bookmark"}
    >
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <TextInput
            required
            label="URL"
            placeholder="https://example.com"
            {...form.getInputProps("url")}
          />
          <TextInput
            label="Title"
            placeholder="Optional title"
            {...form.getInputProps("metadata.title")}
          />
          <TextInput
            label="Description"
            placeholder="Optional description"
            {...form.getInputProps("metadata.description")}
          />
          <Button
            type="submit"
            loading={createMutation.isLoading || updateMutation.isLoading}
          >
            {editingBookmark ? "Update" : "Create"}
          </Button>
        </Stack>
      </form>
    </Modal>
  );

  return (
    <Container size="lg" py="xl">
      <LoadingOverlay visible={isLoading} />
      <Group justify="space-between" mb="xl">
        <Title order={2}>My Bookmarks</Title>
        <Button
          leftSection={<Icon icon="tabler:plus" width={20} height={20} />}
          onClick={() => {
            form.reset();
            setCreateModalOpen(true);
          }}
        >
          Add Bookmark
        </Button>
      </Group>

      <Stack gap="md">
        {bookmarks?.map((bookmark) => (
          <BookmarkCard
            key={bookmark.id}
            bookmark={bookmark}
            onEdit={handleEdit}
            onDelete={(id) => deleteMutation.mutate(id)}
            onSelect={() => setSelectedBookmark(bookmark)}
            refreshMutation={refreshMutation}
          />
        ))}
      </Stack>

      <CreateEditModal
        opened={createModalOpen || !!editingBookmark}
        onClose={() => {
          setCreateModalOpen(false);
          setEditingBookmark(null);
          form.reset();
        }}
      />

      {selectedBookmark && (
        <BookmarkDetailModal
          bookmark={selectedBookmark}
          onClose={() => setSelectedBookmark(null)}
          onEdit={handleEdit}
          onDelete={(id) => deleteMutation.mutate(id)}
          refreshMutation={refreshMutation}
        />
      )}
    </Container>
  );
}