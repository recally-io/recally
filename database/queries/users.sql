
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

-- name: GetUserById :one
SELECT * FROM users WHERE uuid = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserById :one
UPDATE users SET username = $2, email = $3, github = $4,
  google = $5, telegram = $6, 
  activate_assistant_id=$7, activate_thread_id=$8, status = $9
WHERE uuid = $1
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, github, google, telegram, activate_assistant_id, activate_thread_id, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;
