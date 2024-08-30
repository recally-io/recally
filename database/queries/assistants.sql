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
