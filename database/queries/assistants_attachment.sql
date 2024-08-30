-- CRUD for assistant_attachments
-- name: CreateAssistantAttachment :one
INSERT INTO assistant_attachments (user_id, assistant_id, thread_id, name, type, url, size, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAssistantAttachmentById :one
SELECT * FROM assistant_attachments WHERE uuid = $1;

-- name: ListAssistantAttachmentsByUserId :many
SELECT * FROM assistant_attachments WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListAssistantAttachmentsByAssistantId :many
SELECT * FROM assistant_attachments WHERE assistant_id = $1 ORDER BY created_at DESC;

-- name: ListAssistantAttachmentsByThreadId :many
SELECT * FROM assistant_attachments WHERE thread_id = $1 ORDER BY created_at DESC;

-- name: UpdateAssistantAttachment :exec
UPDATE assistant_attachments SET name = $2, type = $3, url = $4, size = $5, metadata = $6 WHERE uuid = $1;

-- name: DeleteAssistantAttachment :exec
DELETE FROM assistant_attachments WHERE uuid = $1;

-- name: DeleteAssistantAttachmentsByAssistantId :exec
DELETE FROM assistant_attachments WHERE assistant_id = $1;

-- name: DeleteAssistantAttachmentsByThreadId :exec
DELETE FROM assistant_attachments WHERE thread_id = $1;

-- name: DeleteAssistantAttachmentsByUserId :exec
DELETE FROM assistant_attachments WHERE user_id = $1;
