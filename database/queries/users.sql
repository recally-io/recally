
-- name: InserUser :one
INSERT INTO users (username, telegram, activate_assistant_id, activate_thread_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTelegramUser :one
SELECT * FROM users WHERE telegram = $1;

-- name: UpdateTelegramUser :one
UPDATE users SET activate_assistant_id = $1, activate_thread_id = $2 WHERE telegram = $3
RETURNING *;

-- name: DeleteTelegramUser :exec
DELETE FROM users WHERE telegram = $1;
