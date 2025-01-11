-- name: CreateAPIKey :one
INSERT INTO auth_api_keys (
    user_id, name, key_prefix, key_hash, scopes, 
    expires_at
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListAPIKeys :many
SELECT * FROM auth_api_keys 
WHERE user_id = $1
    AND (
        (sqlc.narg('prefix')::text IS NULL OR key_prefix = sqlc.narg('prefix')::text)
        AND (sqlc.narg('is_active')::bool IS NULL OR 
            CASE 
                WHEN sqlc.narg('is_active')::bool = true THEN (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
                WHEN sqlc.narg('is_active')::bool = false THEN true
            END
        )
    )
ORDER BY created_at DESC;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE auth_api_keys 
SET last_used_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM auth_api_keys 
WHERE id = $1;


-- name: GetUserByApiKey :one
SELECT u.* FROM users u 
JOIN auth_api_keys ak ON u.uuid = ak.user_id
WHERE ak.key_hash = $1
    AND ak.expires_at > CURRENT_TIMESTAMP;
