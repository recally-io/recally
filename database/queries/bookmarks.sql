-- name: ListBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
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
         LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
         LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
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
           LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
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
         LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
         LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
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
         LEFT JOIN bookmark_content_tags_mapping bctm ON bc.id = bctm.content_id
         LEFT JOIN bookmark_content_tags bct ON bctm.tag_id = bct.id
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

-- name: CreateBookmarkContent :one
INSERT INTO bookmark_content (
  type, title, description, url, domain, s3_key,
  summary, content, html, metadata
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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
    url = COALESCE(sqlc.narg('url'), url),
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

-- Tags related queries similar to content.sql but adapted for new schema
-- name: ListBookmarkTagsByUser :many
SELECT bct.name, count(bctm.*) as count
FROM bookmark_content_tags bct
         JOIN bookmark_content_tags_mapping bctm ON bct.id = bctm.tag_id
WHERE bct.user_id = $1
GROUP BY bct.name
ORDER BY count DESC;

-- name: ListBookmarkContentTags :many
SELECT bct.name
FROM bookmark_content_tags bct
         JOIN bookmark_content_tags_mapping bctm ON bct.id = bctm.tag_id
WHERE bctm.content_id = $1
  AND bct.user_id = $2;

-- name: ListBookmarkDomains :many
SELECT bc.domain, count(*) as count
FROM bookmarks b
         JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.user_id = $1 
AND bc.domain IS NOT NULL
GROUP BY bc.domain
ORDER BY count DESC, domain ASC;

-- name: CreateBookmarkContentTag :one
INSERT INTO bookmark_content_tags (name, user_id)
VALUES ($1, $2)
ON CONFLICT (name, user_id) DO UPDATE
    SET usage_count = bookmark_content_tags.usage_count + 1
RETURNING *;

-- name: DeleteBookmarkContentTag :exec
DELETE FROM bookmark_content_tags
WHERE id = $1
  AND user_id = $2;

-- name: LinkBookmarkContentWithTags :exec
INSERT INTO bookmark_content_tags_mapping (content_id, tag_id)
SELECT $1, bct.id
FROM bookmark_content_tags bct
WHERE bct.name = ANY ($2::text[])
  AND bct.user_id = $3;

-- name: UnLinkBookmarkContentWithTags :exec
DELETE FROM bookmark_content_tags_mapping
WHERE content_id = $1
  AND tag_id IN (SELECT id
                 FROM bookmark_content_tags
                 WHERE name = ANY ($2::text[])
                   AND user_id = $3);

-- name: ListExistingBookmarkTagsByTags :many
SELECT name
FROM bookmark_content_tags
WHERE name = ANY ($1::text[])
  AND user_id = $2;

