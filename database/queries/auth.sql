-- name: CreateUser :one
INSERT INTO users (username, email, phone, password_hash, activate_assistant_id, activate_thread_id, status, settings)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM users WHERE uuid = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: GetUserByUsername :one 
SELECT * FROM users WHERE username = $1;

-- name: UpdateUserById :one
UPDATE users SET username = $2, email = $3, phone = $4, password_hash = $5,
  activate_assistant_id=$6, activate_thread_id=$7, status = $8, settings = $9
WHERE uuid = $1
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE uuid = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;

-- name: ListUsersByStatus :many
SELECT * FROM users WHERE status = $1 ORDER BY created_at DESC;

-- name: CreateOAuthConnection :one
INSERT INTO auth_user_oauth_connections (
    user_id, provider, provider_user_id, provider_email, 
    access_token, refresh_token, token_expires_at, provider_data
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetOAuthConnectionByUserAndProvider :one
SELECT * FROM auth_user_oauth_connections 
WHERE user_id = $1 AND provider = $2;

-- name: GetOAuthConnectionByProviderAndProviderID :one
SELECT * FROM auth_user_oauth_connections 
WHERE provider = $1 AND provider_user_id = $2;

-- name: GetUserByOAuthProviderId :one
SELECT * FROM users 
WHERE uuid = (SELECT user_id FROM auth_user_oauth_connections WHERE provider = $1 AND provider_user_id = $2);

-- name: ListOAuthConnectionsByUser :many
SELECT * FROM auth_user_oauth_connections 
WHERE user_id = $1 
ORDER BY created_at DESC;

-- name: UpdateOAuthConnection :one
UPDATE auth_user_oauth_connections SET 
    provider_email = $3,
    access_token = $4,
    refresh_token = $5,
    token_expires_at = $6,
    provider_data = $7,
    user_id = $8
WHERE provider_user_id = $1 AND provider = $2
RETURNING *;

-- name: DeleteOAuthConnection :exec
DELETE FROM auth_user_oauth_connections 
WHERE user_id = $1 AND provider = $2;

-- name: RevokeToken :one
INSERT INTO auth_revoked_tokens (
    jti, user_id, expires_at, reason
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: IsTokenRevoked :one
SELECT EXISTS (
    SELECT 1 FROM auth_revoked_tokens 
    WHERE jti = $1 AND user_id = $2
) AS is_revoked;

-- name: ListRevokedTokensByUser :many
SELECT * FROM auth_revoked_tokens 
WHERE user_id = $1 
ORDER BY revoked_at DESC;

-- name: DeleteExpiredRevokedTokens :exec
DELETE FROM auth_revoked_tokens 
WHERE expires_at < CURRENT_TIMESTAMP;

-- name: DeleteRevokedToken :exec
DELETE FROM auth_revoked_tokens 
WHERE jti = $1 AND user_id = $2;
