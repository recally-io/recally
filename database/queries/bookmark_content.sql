-- name: IsBookmarkContentExistByURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmark_content
  WHERE url = $1
);

-- name: CreateBookmarkContent :one
INSERT INTO bookmark_content (
  type,
  url,
  user_id,
  title,
  description,
  domain,
  s3_key,
  summary,
  content,
  html,
  tags,
  metadata
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;


-- name: GetBookmarkContentByID :one
SELECT *
FROM bookmark_content
WHERE id = $1;

-- name: GetBookmarkContentByBookmarkID :one
SELECT bc.*
FROM bookmarks b 
  JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.id = $1;

-- name: GetBookmarkContentByURL :one
-- First try to get user specific content, then the shared content
SELECT *
FROM bookmark_content
WHERE url = $1 AND (user_id = $2 OR user_id IS NULL)
LIMIT 1;

-- name: UpdateBookmarkContent :one
UPDATE bookmark_content
SET title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    s3_key = COALESCE(sqlc.narg('s3_key'), s3_key),
    summary = COALESCE(sqlc.narg('summary'), summary),
    content = COALESCE(sqlc.narg('content'), content),
    html = COALESCE(sqlc.narg('html'), html),
    tags = COALESCE(sqlc.narg('tags'), tags),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1
RETURNING *;

-- name: OwnerTransferBookmarkContent :exec
UPDATE bookmark_content
SET 
    user_id = sqlc.narg('new_user_id'),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.narg('user_id');
