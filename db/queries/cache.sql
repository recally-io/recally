-- name: GetCacheByKey :one
SELECT * FROM cache WHERE key = $1;

-- name: CreateCache :exec
INSERT INTO cache (key, value, expires_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: UpdateCache :exec
UPDATE cache SET value = $2, expires_at = $3, updated_at = $4
WHERE key = $1;

-- name: DeleteCacheByKey :exec
DELETE FROM cache WHERE key = $1;

-- name: DeleteExpiredCache :exec
DELETE FROM cache WHERE expires_at < $1;
