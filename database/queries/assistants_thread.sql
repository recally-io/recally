-- CRUD for assistant_threads
-- name: CreateAssistantThread :one
INSERT INTO assistant_threads (uuid, user_id, assistant_id, name, description, system_prompt, model, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAssistantThread :one
SELECT * FROM assistant_threads WHERE uuid = $1;

-- name: UpdateAssistantThread :one
UPDATE assistant_threads SET name = $2, description = $3, model = $4, metadata = $5, system_prompt = $6
WHERE uuid = $1
RETURNING *;

-- name: DeleteAssistantThread :exec
DELETE FROM assistant_threads WHERE uuid = $1;

-- name: DeleteAssistantThreadsByAssistant :exec
DELETE FROM assistant_threads WHERE assistant_id = $1;

-- name: ListAssistantThreadsByUser :many
SELECT * FROM assistant_threads WHERE user_id = $1 ORDER BY created_at DESC;

-- name: ListAssistantThreads :many
SELECT * FROM assistant_threads WHERE assistant_id = $1 ORDER BY created_at DESC;
