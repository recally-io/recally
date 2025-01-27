-- name: CreateShareContent :one
INSERT INTO content_share (user_id, content_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetSharedContent :one
-- get the shared content from content table
SELECT c.*
FROM content_share AS cs
  JOIN content AS c ON cs.content_id = c.id
WHERE cs.id = $1
  AND (cs.expires_at is NULL OR cs.expires_at > now());

-- name: GetShareContent :one
-- info about the shared content
SELECT *
FROM content_share
WHERE content_id = $1
  AND user_id = $2;

-- name: ListShareContent :many
SELECT c.*
FROM content_share AS cs
  JOIN content AS c ON cs.content_id = c.id
WHERE cs.user_id = $1
  AND cs.expires_at is NULL OR cs.expires_at > now()
ORDER BY cs.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateShareContent :one
UPDATE content_share cs
SET expires_at = $3
FROM content c
WHERE cs.content_id = c.id
  AND c.id = $1
  AND c.user_id = $2
RETURNING cs.*;

-- name: DeleteShareContent :exec
DELETE FROM content_share cs
USING content c
WHERE cs.content_id = c.id
  AND c.id = $1
  AND c.user_id = $2;

-- name: DeleteExpiredShareContent :exec
DELETE
FROM content_share
WHERE expires_at < now();
