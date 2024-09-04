-- CRUD for assistant_thread_messages
-- name: CreateThreadMessage :one
INSERT INTO assistant_messages (user_id, assistant_id, thread_id, model, role, text, prompt_token, completion_token, embeddings, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetThreadMessage :one
SELECT * FROM assistant_messages WHERE uuid = $1;

-- name: UpdateThreadMessage :exec
UPDATE assistant_messages SET role = $2, text = $3, model = $4, prompt_token=$5, completion_token=$6,embeddings=$7, metadata=$8 WHERE uuid = $1;

-- name: DeleteThreadMessage :exec
DELETE FROM assistant_messages WHERE uuid = $1;

-- name: DeleteThreadMessageByThreadAndCreatedAt :exec
DELETE FROM assistant_messages WHERE thread_id = $1 AND created_at >= $2;

-- name: DeleteThreadMessagesByThread :exec
DELETE FROM assistant_messages WHERE thread_id = $1;

-- name: DeleteThreadMessagesByAssistant :exec
DELETE FROM assistant_messages
USING assistant_threads
WHERE assistant_messages.thread_id = assistant_threads.uuid
  AND assistant_threads.assistant_id = $1;

-- name: ListThreadMessages :many
SELECT * FROM assistant_messages WHERE thread_id = $1 ORDER BY created_at ASC;

-- name: ListThreadMessagesWithLimit :many
SELECT * FROM assistant_messages WHERE thread_id = $1 ORDER BY created_at DESC LIMIT $2;

-- name: SimilaritySearchMessages :many
SELECT *
FROM assistant_messages
WHERE thread_id = $1 AND (embeddings <=> $2 < 0.5)
ORDER BY 1 - (embeddings <=> $2) DESC
LIMIT $3;
