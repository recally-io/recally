
-- name: CreateBookmarkTag :one
INSERT INTO bookmark_tags (name, user_id)
VALUES ($1, $2)
ON CONFLICT (name, user_id) DO UPDATE
SET name = EXCLUDED.name
RETURNING *;

-- name: DeleteBookmarkTag :exec
DELETE FROM bookmark_tags
WHERE id = $1
  AND user_id = $2;

-- name: LinkBookmarkWithTags :exec
INSERT INTO bookmark_tags_mapping (bookmark_id, tag_id)
SELECT $1, bt.id
FROM bookmark_tags bt
WHERE bt.name = ANY ($2::text[])
  AND bt.user_id = $3;

-- name: UnLinkBookmarkWithTags :exec
DELETE FROM bookmark_tags_mapping
WHERE bookmark_id = $1
  AND tag_id IN (SELECT id
                 FROM bookmark_tags
                 WHERE name = ANY ($2::text[])
                   AND user_id = $3);

-- name: ListExistingBookmarkTagsByTags :many
SELECT name
FROM bookmark_tags
WHERE name = ANY ($1::text[])
  AND user_id = $2;


-- name: ListBookmarkTagsByUser :many
SELECT name, count(*) as cnt
FROM bookmark_tags
WHERE user_id = $1
GROUP BY name
ORDER BY cnt DESC;

-- name: ListBookmarkTagsByBookmarkId :many
SELECT bt.name
FROM bookmark_tags bt
  JOIN bookmark_tags_mapping btm ON bt.id = btm.tag_id
WHERE btm.bookmark_id = $1;

-- name: OwnerTransferBookmarkTag :exec
UPDATE bookmark_tags
SET 
    user_id = sqlc.narg('new_user_id'),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.narg('user_id');
