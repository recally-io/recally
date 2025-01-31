-- name: CreateBookmarkShare :one
INSERT INTO bookmark_share (user_id, bookmark_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBookmarkShareContent :one
SELECT bc.*
FROM bookmark_share AS bs
  JOIN bookmarks AS b ON bs.bookmark_id = b.id
  JOIN bookmark_content AS bc ON b.content_id = bc.id
WHERE bs.id = $1
  AND (bs.expires_at is NULL OR bs.expires_at > now());

-- name: GetBookmarkShare :one
SELECT *
FROM bookmark_share
WHERE bookmark_id = $1
  AND user_id = $2;

-- name: UpdateBookmarkShareByBookmarkId :one
UPDATE bookmark_share bs
SET expires_at = $3
FROM bookmarks b
WHERE bs.bookmark_id = b.id
  AND b.id = $1
  AND b.user_id = $2
RETURNING bs.*;

-- name: DeleteBookmarkShareByBookmarkId :exec
DELETE FROM bookmark_share bs
USING bookmarks b
WHERE bs.bookmark_id = b.id
  AND b.id = $1
  AND b.user_id = $2;

-- name: DeleteExpiredBookmarkShare :exec
DELETE
FROM bookmark_share
WHERE expires_at < now();
