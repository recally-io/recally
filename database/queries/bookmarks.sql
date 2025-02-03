-- name: ListBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.bookmark_id
           LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND (sqlc.narg('domains')::text[] IS NULL OR bc.domain = ANY(sqlc.narg('domains')::text[]))
    AND (sqlc.narg('types')::text[] IS NULL OR bc.type = ANY(sqlc.narg('types')::text[]))
    AND (sqlc.narg('tags')::text[] IS NULL OR bct.name = ANY(sqlc.narg('tags')::text[]))
)
SELECT sqlc.embed(b),
       sqlc.embed(bc),
       t.total_count,
       COALESCE(
         array_agg(bct.name) FILTER (WHERE bct.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) AS tags
FROM bookmarks AS b
         JOIN bookmark_content AS bc ON b.content_id = bc.id
         CROSS JOIN total AS t
         LEFT JOIN bookmark_tags_mapping AS bctm ON b.id = bctm.bookmark_id
         LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
WHERE b.user_id = $1
  AND (sqlc.narg('domains')::text[] IS NULL OR bc.domain = ANY(sqlc.narg('domains')::text[]))
  AND (sqlc.narg('types')::text[] IS NULL OR bc.type = ANY(sqlc.narg('types')::text[]))
  AND (sqlc.narg('tags')::text[] IS NULL OR bct.name = ANY(sqlc.narg('tags')::text[]))
GROUP BY b.id, bc.id, t.total_count
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_tags_mapping AS bctm ON b.id = bctm.bookmark_id
           LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND (sqlc.narg('domains')::text[] IS NULL OR bc.domain = ANY(sqlc.narg('domains')::text[]))
    AND (sqlc.narg('types')::text[] IS NULL OR bc.type = ANY(sqlc.narg('types')::text[]))
    AND (sqlc.narg('tags')::text[] IS NULL OR bct.name = ANY(sqlc.narg('tags')::text[]))
    AND (
      sqlc.narg('query')::text IS NULL
      OR bc.title @@@ sqlc.narg('query')
      OR bc.description @@@ sqlc.narg('query')
      OR bc.summary @@@ sqlc.narg('query')
      OR bc.content @@@ sqlc.narg('query')
      OR bc.metadata @@@ sqlc.narg('query')
    )
)
SELECT sqlc.embed(b),
       sqlc.embed(bc),
       t.total_count,
       COALESCE(
         array_agg(bct.name) FILTER (WHERE bct.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) AS tags
FROM bookmarks AS b
         JOIN bookmark_content AS bc ON b.content_id = bc.id
         CROSS JOIN total AS t
         LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.bookmark_id
         LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
WHERE b.user_id = $1
  AND (sqlc.narg('domains')::text[] IS NULL OR bc.domain = ANY(sqlc.narg('domains')::text[]))
  AND (sqlc.narg('types')::text[] IS NULL OR bc.type = ANY(sqlc.narg('types')::text[]))
  AND (sqlc.narg('tags')::text[] IS NULL OR bct.name = ANY(sqlc.narg('tags')::text[]))
  AND (
    sqlc.narg('query')::text IS NULL
    OR bc.title @@@ sqlc.narg('query')
    OR bc.description @@@ sqlc.narg('query')
    OR bc.summary @@@ sqlc.narg('query')
    OR bc.content @@@ sqlc.narg('query')
    OR bc.metadata @@@ sqlc.narg('query')
  )
GROUP BY b.id, bc.id, t.total_count
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetBookmarkWithContent :one
SELECT sqlc.embed(b),
       sqlc.embed(bc),
       sqlc.embed(bs),
       COALESCE(
         array_agg(bt.name) FILTER (WHERE bt.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) as tags
FROM bookmarks b
         JOIN bookmark_content bc ON b.content_id = bc.id
         LEFT JOIN bookmark_share bs ON bs.bookmark_id = b.id
         LEFT JOIN bookmark_tags_mapping btm ON btm.bookmark_id = b.id
         LEFT JOIN bookmark_tags bt ON btm.tag_id = bt.id
WHERE b.id = $1
  AND b.user_id = $2
GROUP BY b.id, bc.id, bs.id
LIMIT 1;

-- name: IsBookmarkExistWithURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmarks b
           JOIN bookmark_content bc ON b.content_id = bc.id
  WHERE bc.url = $1
    AND b.user_id = $2
);

-- name: CreateBookmark :one
INSERT INTO bookmarks (
  user_id, content_id, is_favorite, is_archive, metadata
)
VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateBookmark :one
UPDATE bookmarks
SET is_favorite = COALESCE(sqlc.narg('is_favorite'), is_favorite),
    is_archive = COALESCE(sqlc.narg('is_archive'), is_archive),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1
  AND user_id = $2
RETURNING *;


-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE id = $1 AND user_id = $2;

-- name: DeleteBookmarksByUser :exec
DELETE FROM bookmarks
WHERE user_id = $1;

-- name: OwnerTransferBookmark :exec
UPDATE bookmarks 
SET 
    user_id = sqlc.narg('new_user_id'),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = sqlc.narg('user_id');

-- name: ListBookmarkDomainsByUser :many
SELECT bc.domain, count(*) as cnt
FROM bookmarks b
  JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.user_id = $1 
AND bc.domain IS NOT NULL
GROUP BY bc.domain
ORDER BY cnt DESC, domain ASC;
