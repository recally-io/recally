-- CRUD for assistants

-- name: CreateAssistant :one
INSERT INTO assistants (user_id, name, description, system_prompt, model, metadata)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAssistant :one
SELECT * FROM assistants WHERE uuid = $1;

-- name: UpdateAssistant :one
UPDATE assistants SET name = $2, description = $3, system_prompt = $4, model = $5, metadata = $6
WHERE uuid = $1
RETURNING *;

-- name: DeleteAssistant :exec
DELETE FROM assistants WHERE uuid = $1;

-- name: ListAssistantsByUser :many
SELECT * FROM assistants WHERE user_id = $1 ORDER BY created_at DESC;

-- CRUD for assistant_threads
-- name: CreateAssistantThread :one
INSERT INTO assistant_threads (user_id, assistant_id, name, description, system_prompt, model, is_long_term_memory, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAssistantThread :one
SELECT * FROM assistant_threads WHERE uuid = $1;

-- name: UpdateAssistantThread :exec
UPDATE assistant_threads SET name = $2, description = $3, model = $4, is_long_term_memory = $5, metadata = $6, system_prompt = $7 
WHERE uuid = $1;

-- name: DeleteAssistantThread :exec
DELETE FROM assistant_threads WHERE uuid = $1;

-- name: ListAssistantThreadsByUser :many
SELECT * FROM assistant_threads WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListAssistantThreads :many
SELECT * FROM assistant_threads WHERE assistant_id = $1 ORDER BY created_at DESC;

-- CRUD for assistant_thread_messages
-- name: CreateThreadMessage :one
INSERT INTO assistant_messages (user_id, thread_id, model, token, role, text, attachments, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetThreadMessage :one
SELECT * FROM assistant_messages WHERE uuid = $1;

-- name: UpdateThreadMessage :exec
UPDATE assistant_messages SET text = $2, attachments = $3, metadata = $4 WHERE uuid = $1;

-- name: DeleteThreadMessage :exec
DELETE FROM assistant_messages WHERE uuid = $1;

-- name: ListThreadMessages :many
SELECT * FROM assistant_messages WHERE thread_id = $1 ORDER BY created_at ASC;

-- CRUD for assistant_attachments
-- name: CreateAttachment :one
INSERT INTO assistant_attachments (user_id, entity, entity_id, file_type, file_url, size, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAttachment :one
SELECT * FROM assistant_attachments WHERE uuid = $1;

-- name: UpdateAttachment :exec
UPDATE assistant_attachments SET file_type = $2, file_url = $3, size = $4, metadata = $5 WHERE uuid = $1;

-- name: DeleteAttachment :exec
DELETE FROM assistant_attachments WHERE uuid = $1;

-- name: ListAttachments :many
SELECT * FROM assistant_attachments WHERE entity = $1 AND entity_id = $2 ORDER BY created_at DESC;

-- name: ListAttachmentsByUser :many
SELECT * FROM assistant_attachments WHERE user_id = $1 ORDER BY created_at DESC;

-- CRUD for assistant_message_embedddings
-- name: CreateAssistantEmbedding :one
INSERT INTO assistant_embedddings (user_id, message_id, attachment_id, embeddings)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteAssistantEmbeddings :exec
DELETE FROM assistant_embedddings WHERE id = $1;

-- It need combine all these results to get the final result:
-- 1. assistants -> assistant_attachments -> assistant_message_embedddings 
-- 2. assistant_threads -> assistant_attachments -> assistant_message_embedddings 
-- 3. assistant_threads -> assistant_messages -> assistant_attachments -> assistant_message_embedddings
-- name: SimilaritySearchForThreadByCosineDistance :many
SELECT ae.id, ae.text, 1 - (embeddings <=> $2) AS score  
FROM assistant_embedddings ae
JOIN assistant_attachments aa ON ae.attachment_id = aa.uuid
WHERE ae.attachment_id IN (
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistant_messages am ON aa.entity_id = am.uuid
        JOIN assistant_threads at ON am.thread_id = at.uuid
        WHERE at.uuid = $1
    UNION
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistant_threads at ON aa.entity_id = at.uuid
        WHERE at.uuid = $1
    UNION
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistants a ON aa.entity_id = a.uuid
        JOIN assistant_threads at ON a.uuid = at.assistant_id
        WHERE at.uuid = $1
)
AND embeddings <=> $2
ORDER BY 1 - (embedding <=> $2) LIMIT $3;
