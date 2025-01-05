-- name: ListContents :many
WITH total AS (SELECT COUNT(*) AS total_count
               FROM content AS tc
                        LEFT JOIN content_tags_mapping AS tctm ON tc.id = tctm.content_id
                        LEFT JOIN content_tags AS tct ON tctm.tag_id = tct.id
               WHERE tc.user_id = $1
                 AND (
                   sqlc.narg('domains') :: text[] IS NULL
                       OR tc.domain = ANY (sqlc.narg('domains') :: text[])
                   )
                 AND (
                   sqlc.narg('types') :: text[] IS NULL
                       OR tc.type = ANY (sqlc.narg('types') :: text[])
                   )
                 AND (
                   sqlc.narg('tags') :: text[] IS NULL
                       OR tct.name = ANY (sqlc.narg('tags') :: text[])
                   ))
SELECT c.*,
       t.total_count,
       COALESCE(
                       array_agg(ct.name) FILTER (
                   WHERE
                   ct.name IS NOT NULL
                   ),
                       ARRAY [] :: VARCHAR[]
       ) AS tags
FROM content AS c
         CROSS JOIN total AS t
         LEFT JOIN content_tags_mapping AS ctm ON c.id = ctm.content_id
         LEFT JOIN content_tags AS ct ON ctm.tag_id = ct.id
WHERE c.user_id = $1
  AND (
    sqlc.narg('domains') :: text[] IS NULL
        OR c.domain = ANY (sqlc.narg('domains') :: text[])
    )
  AND (
    sqlc.narg('types') :: text[] IS NULL
        OR c.type = ANY (sqlc.narg('types') :: text[])
    )
  AND (
    sqlc.narg('tags') :: text[] IS NULL
        OR ct.name = ANY (sqlc.narg('tags') :: text[])
    )
GROUP BY c.id,
         t.total_count
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetContent :one
SELECT c.*,
       COALESCE(
                       array_agg(ct.name) FILTER (
                   WHERE
                   ct.name IS NOT NULL
                   ),
                       ARRAY [] :: VARCHAR[]
       ) as tags
FROM content c
         LEFT JOIN content_tags_mapping ctm ON c.id = ctm.content_id
         LEFT JOIN content_tags ct ON ctm.tag_id = ct.id
WHERE c.id = $1
  AND c.user_id = $2
GROUP BY c.id
LIMIT 1;

-- name: IsContentExistWithURL :one
SELECT EXISTS (SELECT 1
               FROM content
               WHERE url = $1
                 AND user_id = $2);

-- name: CreateContent :one
INSERT INTO content (user_id,
                     type,
                     title,
                     description,
                     url,
                     domain,
                     s3_key,
                     summary,
                     content,
                     html,
                     metadata,
                     is_favorite)
VALUES ($1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12)
RETURNING *;

-- name: UpdateContent :one
UPDATE
    content
SET title       = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    url         = COALESCE(sqlc.narg('url'), url),
    domain      = COALESCE(sqlc.narg('domain'), domain),
    s3_key      = COALESCE(sqlc.narg('s3_key'), s3_key),
    summary     = COALESCE(sqlc.narg('summary'), summary),
    content     = COALESCE(sqlc.narg('content'), content),
    html        = COALESCE(sqlc.narg('html'), html),
    metadata    = COALESCE(sqlc.narg('metadata'), metadata),
    is_favorite = COALESCE(sqlc.narg('is_favorite'), is_favorite)
WHERE id = $1
  AND user_id = $2
RETURNING *;

-- name: DeleteContent :exec
DELETE
FROM content
WHERE id = $1
  AND user_id = $2;

-- name: DeleteContentsByUser :exec
DELETE
FROM content
WHERE user_id = $1;

-- name: OwnerTransferContent :exec
UPDATE
    content
SET user_id = $2
WHERE id = $1
  AND user_id = $3;

-- name: ListTagsByUser :many
SELECT ct.name, count(*) as count
FROM content_tags ct
WHERE ct.user_id = $1
GROUP BY ct.name
ORDER BY count DESC;

-- name: ListContentTags :many
SELECT ct.name
FROM content_tags ct
         JOIN content_tags_mapping ctm ON ct.id = ctm.tag_id
WHERE ctm.content_id = $1
  AND ct.user_id = $2;

-- name: ListContentDomains :many
SELECT domain, count(*) as count
FROM content
WHERE user_id = $1 
AND domain IS NOT NULL
GROUP BY domain
ORDER BY count DESC, domain ASC;

-- name: CreateContentTag :one
INSERT INTO content_tags (name, user_id)
VALUES ($1, $2)
ON CONFLICT (name, user_id) DO UPDATE
    SET usage_count = content_tags.usage_count + 1
RETURNING *;

-- name: DeleteContentTag :exec
DELETE
FROM content_tags
WHERE id = $1
  AND user_id = $2;

-- name: LinkContentWithTags :exec
-- $1: content_id, $2: text[], $3: user_id
INSERT INTO content_tags_mapping (content_id, tag_id)
SELECT $1,
       ct.id
FROM content_tags ct
WHERE ct.name = ANY ($2 :: text[])
  AND ct.user_id = $3;

-- name: UnLinkContentWithTags :exec
-- $1: content_id, $2: text[], $3: user_id
DELETE FROM content_tags_mapping
WHERE content_id = $1
  AND tag_id IN (SELECT id
                 FROM content_tags
                 WHERE name = ANY ($2 :: text[])
                   AND user_id = $3);

-- name: ListExistingTagsByTags :many
SELECT name
FROM content_tags
WHERE name = ANY ($1 :: text[])
  AND user_id = $2;
