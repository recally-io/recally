-- name: CreateBookmark :one
INSERT INTO bookmarks (
    uuid,
    user_id,
    url,
    title,
    summary,
    summary_embeddings,
    content,
    content_embeddings,
    html,
    metadata,
    screenshot
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetBookmarkByUUID :one
SELECT * FROM bookmarks WHERE uuid = $1;

-- name: GetBookmarkByURL :one
SELECT * FROM bookmarks WHERE url = $1 AND user_id = $2;

-- name: ListBookmarks :many
WITH total AS (
    SELECT COUNT(*) AS total_count 
    FROM bookmarks 
    WHERE user_id = $1
)
SELECT b.*, t.total_count 
FROM bookmarks b, total t
WHERE b.user_id = $1 
ORDER BY b.updated_at DESC 
LIMIT $2 OFFSET $3;

-- name: UpdateBookmark :one
UPDATE bookmarks 
SET 
    title = COALESCE($3, title),
    summary = COALESCE($4, summary),
    summary_embeddings = COALESCE($5, summary_embeddings),
    content = COALESCE($6, content),
    content_embeddings = COALESCE($7, content_embeddings),
    html = COALESCE($8, html),
    metadata = COALESCE($9, metadata),
    screenshot = COALESCE($10, screenshot),
    updated_at = CURRENT_TIMESTAMP
WHERE uuid = $1 AND user_id = $2
RETURNING *;

-- name: DeleteBookmark :exec
DELETE FROM bookmarks WHERE uuid = $1 AND user_id = $2;

-- name: DeleteBookmarksByUser :exec
DELETE FROM bookmarks WHERE user_id = $1;

-- name: OwnerTransferBookmark :exec
UPDATE bookmarks 
SET 
    user_id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;
