-- name: ListBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND (sqlc.narg('domains')::text[] IS NULL OR bc.domain = ANY(sqlc.narg('domains')::text[]))
    AND (sqlc.narg('types')::text[] IS NULL OR bc.type = ANY(sqlc.narg('types')::text[]))
    AND (sqlc.narg('tags')::text[] IS NULL OR bct.name = ANY(sqlc.narg('tags')::text[]))
)
SELECT b.*,
       bc.*,
       t.total_count,
       COALESCE(
         array_agg(bct.name) FILTER (WHERE bct.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) AS tags
FROM bookmarks AS b
         JOIN bookmark_content AS bc ON b.content_id = bc.id
         CROSS JOIN total AS t
         LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
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
           LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
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
SELECT b.*,
       bc.*,
       t.total_count,
       COALESCE(
         array_agg(bct.name) FILTER (WHERE bct.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) AS tags
FROM bookmarks AS b
         JOIN bookmark_content AS bc ON b.content_id = bc.id
         CROSS JOIN total AS t
         LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
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

-- name: GetBookmark :one
SELECT b.*,
       bc.*,
       COALESCE(
         array_agg(bct.name) FILTER (WHERE bct.name IS NOT NULL),
         ARRAY[]::VARCHAR[]
       ) as tags
FROM bookmarks b
         JOIN bookmark_content bc ON b.content_id = bc.id
         LEFT JOIN bookmark_tags_mapping bctm ON bc.id = bctm.content_id
         LEFT JOIN bookmark_tags bct ON bctm.tag_id = bct.id
WHERE b.id = $1
  AND b.user_id = $2
GROUP BY b.id, bc.id
LIMIT 1;

-- name: IsBookmarkExistWithURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmarks b
           JOIN bookmark_content bc ON b.content_id = bc.id
  WHERE bc.url = $1
    AND b.user_id = $2
);

-- name: IsBookmarkContentExistWithURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmark_content bc
  WHERE bc.url = $1
);


-- name: CreateBookmarkContent :one
INSERT INTO bookmark_content (
  type, title, description, user_id, url, domain, s3_key,
  summary, content, html, metadata
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: CreateBookmark :one
INSERT INTO bookmarks (
  user_id, content_id, is_favorite, is_archive,
  is_public, reading_progress, metadata
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateBookmark :one
UPDATE bookmarks
SET is_favorite = COALESCE(sqlc.narg('is_favorite'), is_favorite),
    is_archive = COALESCE(sqlc.narg('is_archive'), is_archive),
    is_public = COALESCE(sqlc.narg('is_public'), is_public),
    reading_progress = COALESCE(sqlc.narg('reading_progress'), reading_progress),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: UpdateBookmarkContent :one
UPDATE bookmark_content
SET title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    domain = COALESCE(sqlc.narg('domain'), domain),
    s3_key = COALESCE(sqlc.narg('s3_key'), s3_key),
    summary = COALESCE(sqlc.narg('summary'), summary),
    content = COALESCE(sqlc.narg('content'), content),
    html = COALESCE(sqlc.narg('html'), html),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = $1
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
    user_id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: ListBookmarkTagsByUser :many
SELECT name FROM bookmark_tags
WHERE user_id = $1;

-- name: ListBookmarkTagsByBookmarkId :many
SELECT bt.name
FROM bookmark_tags bt
  JOIN bookmark_tags_mapping btm ON bt.id = btm.tag_id
WHERE btm.bookmark_id = $1;

-- name: ListBookmarkDomains :many
SELECT bc.domain, count(*) as cnt
FROM bookmarks b
  JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.user_id = $1 
AND bc.domain IS NOT NULL
GROUP BY bc.domain
ORDER BY cnt DESC, domain ASC;

-- name: CreateBookmarkTag :one
INSERT INTO bookmark_tags (name, user_id)
VALUES ($1, $2)
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
