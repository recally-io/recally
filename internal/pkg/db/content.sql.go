// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: content.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createContent = `-- name: CreateContent :one
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
RETURNING id, user_id, type, title, description, url, domain, s3_key, summary, content, html, metadata, is_favorite, created_at, updated_at
`

type CreateContentParams struct {
	UserID      uuid.UUID
	Type        string
	Title       string
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
	IsFavorite  pgtype.Bool
}

func (q *Queries) CreateContent(ctx context.Context, db DBTX, arg CreateContentParams) (Content, error) {
	row := db.QueryRow(ctx, createContent,
		arg.UserID,
		arg.Type,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.Domain,
		arg.S3Key,
		arg.Summary,
		arg.Content,
		arg.Html,
		arg.Metadata,
		arg.IsFavorite,
	)
	var i Content
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Metadata,
		&i.IsFavorite,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createContentTag = `-- name: CreateContentTag :one
INSERT INTO content_tags (name, user_id)
VALUES ($1, $2)
ON CONFLICT (name, user_id) DO UPDATE
    SET usage_count = content_tags.usage_count + 1
RETURNING id, name, user_id, usage_count, created_at, updated_at
`

type CreateContentTagParams struct {
	Name   string
	UserID uuid.UUID
}

func (q *Queries) CreateContentTag(ctx context.Context, db DBTX, arg CreateContentTagParams) (ContentTag, error) {
	row := db.QueryRow(ctx, createContentTag, arg.Name, arg.UserID)
	var i ContentTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.UsageCount,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createShareContent = `-- name: CreateShareContent :one
INSERT INTO content_share (user_id, content_id, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, content_id, expires_at, created_at, updated_at
`

type CreateShareContentParams struct {
	UserID    uuid.UUID
	ContentID pgtype.UUID
	ExpiresAt pgtype.Timestamptz
}

func (q *Queries) CreateShareContent(ctx context.Context, db DBTX, arg CreateShareContentParams) (ContentShare, error) {
	row := db.QueryRow(ctx, createShareContent, arg.UserID, arg.ContentID, arg.ExpiresAt)
	var i ContentShare
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteContent = `-- name: DeleteContent :exec
DELETE
FROM content
WHERE id = $1
  AND user_id = $2
`

type DeleteContentParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteContent(ctx context.Context, db DBTX, arg DeleteContentParams) error {
	_, err := db.Exec(ctx, deleteContent, arg.ID, arg.UserID)
	return err
}

const deleteContentTag = `-- name: DeleteContentTag :exec
DELETE
FROM content_tags
WHERE id = $1
  AND user_id = $2
`

type DeleteContentTagParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteContentTag(ctx context.Context, db DBTX, arg DeleteContentTagParams) error {
	_, err := db.Exec(ctx, deleteContentTag, arg.ID, arg.UserID)
	return err
}

const deleteContentsByUser = `-- name: DeleteContentsByUser :exec
DELETE
FROM content
WHERE user_id = $1
`

func (q *Queries) DeleteContentsByUser(ctx context.Context, db DBTX, userID uuid.UUID) error {
	_, err := db.Exec(ctx, deleteContentsByUser, userID)
	return err
}

const deleteExpiredShareContent = `-- name: DeleteExpiredShareContent :exec
DELETE
FROM content_share
WHERE expires_at < now()
`

func (q *Queries) DeleteExpiredShareContent(ctx context.Context, db DBTX) error {
	_, err := db.Exec(ctx, deleteExpiredShareContent)
	return err
}

const deleteShareContent = `-- name: DeleteShareContent :exec
DELETE FROM content_share cs
USING content c
WHERE cs.content_id = c.id
  AND c.id = $1
  AND c.user_id = $2
`

type DeleteShareContentParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteShareContent(ctx context.Context, db DBTX, arg DeleteShareContentParams) error {
	_, err := db.Exec(ctx, deleteShareContent, arg.ID, arg.UserID)
	return err
}

const getContent = `-- name: GetContent :one
SELECT c.id, c.user_id, c.type, c.title, c.description, c.url, c.domain, c.s3_key, c.summary, c.content, c.html, c.metadata, c.is_favorite, c.created_at, c.updated_at,
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
LIMIT 1
`

type GetContentParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type GetContentRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        string
	Title       string
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
	IsFavorite  pgtype.Bool
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Tags        interface{}
}

func (q *Queries) GetContent(ctx context.Context, db DBTX, arg GetContentParams) (GetContentRow, error) {
	row := db.QueryRow(ctx, getContent, arg.ID, arg.UserID)
	var i GetContentRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Metadata,
		&i.IsFavorite,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Tags,
	)
	return i, err
}

const getShareContent = `-- name: GetShareContent :one
SELECT id, user_id, content_id, expires_at, created_at, updated_at
FROM content_share
WHERE content_id = $1
  AND user_id = $2
`

type GetShareContentParams struct {
	ContentID pgtype.UUID
	UserID    uuid.UUID
}

// info about the shared content
func (q *Queries) GetShareContent(ctx context.Context, db DBTX, arg GetShareContentParams) (ContentShare, error) {
	row := db.QueryRow(ctx, getShareContent, arg.ContentID, arg.UserID)
	var i ContentShare
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSharedContent = `-- name: GetSharedContent :one
SELECT c.id, c.user_id, c.type, c.title, c.description, c.url, c.domain, c.s3_key, c.summary, c.content, c.html, c.metadata, c.is_favorite, c.created_at, c.updated_at
FROM content_share AS cs
  JOIN content AS c ON cs.content_id = c.id
WHERE cs.id = $1
  AND (cs.expires_at is NULL OR cs.expires_at > now())
`

// get the shared content from content table
func (q *Queries) GetSharedContent(ctx context.Context, db DBTX, id uuid.UUID) (Content, error) {
	row := db.QueryRow(ctx, getSharedContent, id)
	var i Content
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Metadata,
		&i.IsFavorite,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const isContentExistWithURL = `-- name: IsContentExistWithURL :one
SELECT EXISTS (SELECT 1
               FROM content
               WHERE url = $1
                 AND user_id = $2)
`

type IsContentExistWithURLParams struct {
	Url    pgtype.Text
	UserID uuid.UUID
}

func (q *Queries) IsContentExistWithURL(ctx context.Context, db DBTX, arg IsContentExistWithURLParams) (bool, error) {
	row := db.QueryRow(ctx, isContentExistWithURL, arg.Url, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const linkContentWithTags = `-- name: LinkContentWithTags :exec
INSERT INTO content_tags_mapping (content_id, tag_id)
SELECT $1,
       ct.id
FROM content_tags ct
WHERE ct.name = ANY ($2 :: text[])
  AND ct.user_id = $3
`

type LinkContentWithTagsParams struct {
	ContentID uuid.UUID
	Column2   []string
	UserID    uuid.UUID
}

// $1: content_id, $2: text[], $3: user_id
func (q *Queries) LinkContentWithTags(ctx context.Context, db DBTX, arg LinkContentWithTagsParams) error {
	_, err := db.Exec(ctx, linkContentWithTags, arg.ContentID, arg.Column2, arg.UserID)
	return err
}

const listContentDomains = `-- name: ListContentDomains :many
SELECT domain, count(*) as count
FROM content
WHERE user_id = $1 
AND domain IS NOT NULL
GROUP BY domain
ORDER BY count DESC, domain ASC
`

type ListContentDomainsRow struct {
	Domain pgtype.Text
	Count  int64
}

func (q *Queries) ListContentDomains(ctx context.Context, db DBTX, userID uuid.UUID) ([]ListContentDomainsRow, error) {
	rows, err := db.Query(ctx, listContentDomains, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListContentDomainsRow
	for rows.Next() {
		var i ListContentDomainsRow
		if err := rows.Scan(&i.Domain, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listContentTags = `-- name: ListContentTags :many
SELECT ct.name
FROM content_tags ct
         JOIN content_tags_mapping ctm ON ct.id = ctm.tag_id
WHERE ctm.content_id = $1
  AND ct.user_id = $2
`

type ListContentTagsParams struct {
	ContentID uuid.UUID
	UserID    uuid.UUID
}

func (q *Queries) ListContentTags(ctx context.Context, db DBTX, arg ListContentTagsParams) ([]string, error) {
	rows, err := db.Query(ctx, listContentTags, arg.ContentID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listContents = `-- name: ListContents :many
WITH total AS (
  SELECT COUNT( DISTINCT tc.*) AS total_count
               FROM content AS tc
                        LEFT JOIN content_tags_mapping AS tctm ON tc.id = tctm.content_id
                        LEFT JOIN content_tags AS tct ON tctm.tag_id = tct.id
               WHERE tc.user_id = $1
                 AND (
                   $4 :: text[] IS NULL
                       OR tc.domain = ANY ($4 :: text[])
                   )
                 AND (
                   $5 :: text[] IS NULL
                       OR tc.type = ANY ($5 :: text[])
                   )
                 AND (
                   $6 :: text[] IS NULL
                       OR tct.name = ANY ($6 :: text[])
                   )
)
SELECT c.id, c.user_id, c.type, c.title, c.description, c.url, c.domain, c.s3_key, c.summary, c.content, c.html, c.metadata, c.is_favorite, c.created_at, c.updated_at,
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
    $4 :: text[] IS NULL
        OR c.domain = ANY ($4 :: text[])
    )
  AND (
    $5 :: text[] IS NULL
        OR c.type = ANY ($5 :: text[])
    )
  AND (
    $6 :: text[] IS NULL
        OR ct.name = ANY ($6 :: text[])
    )
GROUP BY c.id,
         t.total_count
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3
`

type ListContentsParams struct {
	UserID  uuid.UUID
	Limit   int32
	Offset  int32
	Domains []string
	Types   []string
	Tags    []string
}

type ListContentsRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        string
	Title       string
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
	IsFavorite  pgtype.Bool
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	TotalCount  int64
	Tags        interface{}
}

func (q *Queries) ListContents(ctx context.Context, db DBTX, arg ListContentsParams) ([]ListContentsRow, error) {
	rows, err := db.Query(ctx, listContents,
		arg.UserID,
		arg.Limit,
		arg.Offset,
		arg.Domains,
		arg.Types,
		arg.Tags,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListContentsRow
	for rows.Next() {
		var i ListContentsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.Domain,
			&i.S3Key,
			&i.Summary,
			&i.Content,
			&i.Html,
			&i.Metadata,
			&i.IsFavorite,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TotalCount,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listExistingTagsByTags = `-- name: ListExistingTagsByTags :many
SELECT name
FROM content_tags
WHERE name = ANY ($1 :: text[])
  AND user_id = $2
`

type ListExistingTagsByTagsParams struct {
	Column1 []string
	UserID  uuid.UUID
}

func (q *Queries) ListExistingTagsByTags(ctx context.Context, db DBTX, arg ListExistingTagsByTagsParams) ([]string, error) {
	rows, err := db.Query(ctx, listExistingTagsByTags, arg.Column1, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listShareContent = `-- name: ListShareContent :many
SELECT c.id, c.user_id, c.type, c.title, c.description, c.url, c.domain, c.s3_key, c.summary, c.content, c.html, c.metadata, c.is_favorite, c.created_at, c.updated_at
FROM content_share AS cs
  JOIN content AS c ON cs.content_id = c.id
WHERE cs.user_id = $1
  AND cs.expires_at is NULL OR cs.expires_at > now()
ORDER BY cs.created_at DESC
LIMIT $2 OFFSET $3
`

type ListShareContentParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) ListShareContent(ctx context.Context, db DBTX, arg ListShareContentParams) ([]Content, error) {
	rows, err := db.Query(ctx, listShareContent, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Content
	for rows.Next() {
		var i Content
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.Domain,
			&i.S3Key,
			&i.Summary,
			&i.Content,
			&i.Html,
			&i.Metadata,
			&i.IsFavorite,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTagsByUser = `-- name: ListTagsByUser :many
SELECT ct.name, count(ctm.*) as count
FROM content_tags ct
  JOIN content_tags_mapping ctm ON ct.id = ctm.tag_id
WHERE ct.user_id = $1
GROUP BY ct.name
ORDER BY count DESC
`

type ListTagsByUserRow struct {
	Name  string
	Count int64
}

func (q *Queries) ListTagsByUser(ctx context.Context, db DBTX, userID uuid.UUID) ([]ListTagsByUserRow, error) {
	rows, err := db.Query(ctx, listTagsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTagsByUserRow
	for rows.Next() {
		var i ListTagsByUserRow
		if err := rows.Scan(&i.Name, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ownerTransferContent = `-- name: OwnerTransferContent :exec
UPDATE
    content
SET user_id = $2
WHERE id = $1
  AND user_id = $3
`

type OwnerTransferContentParams struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	UserID_2 uuid.UUID
}

func (q *Queries) OwnerTransferContent(ctx context.Context, db DBTX, arg OwnerTransferContentParams) error {
	_, err := db.Exec(ctx, ownerTransferContent, arg.ID, arg.UserID, arg.UserID_2)
	return err
}

const searchContentsWithFilter = `-- name: SearchContentsWithFilter :many
WITH total AS (
  SELECT COUNT( DISTINCT tc.*) AS total_count
               FROM content AS tc
                        LEFT JOIN content_tags_mapping AS tctm ON tc.id = tctm.content_id
                        LEFT JOIN content_tags AS tct ON tctm.tag_id = tct.id
               WHERE tc.user_id = $1
                 AND (
                   $4 :: text[] IS NULL
                       OR tc.domain = ANY ($4 :: text[])
                   )
                 AND (
                   $5 :: text[] IS NULL
                       OR tc.type = ANY ($5 :: text[])
                   )
                 AND (
                   $6 :: text[] IS NULL
                       OR tct.name = ANY ($6 :: text[])
                   )
                AND (
                  $7 :: text IS NULL
                      OR tc.title @@@ $7
                      OR tc.description @@@ $7
                      OR tc.summary @@@ $7
                      OR tc.content @@@ $7
                      OR tc.metadata @@@ $7
                    )
)
SELECT c.id, c.user_id, c.type, c.title, c.description, c.url, c.domain, c.s3_key, c.summary, c.content, c.html, c.metadata, c.is_favorite, c.created_at, c.updated_at,
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
    $4 :: text[] IS NULL
        OR c.domain = ANY ($4 :: text[])
    )
  AND (
    $5 :: text[] IS NULL
        OR c.type = ANY ($5 :: text[])
    )
  AND (
    $6 :: text[] IS NULL
        OR ct.name = ANY ($6 :: text[])
    )
  AND (
    $7 :: text IS NULL
        OR c.title @@@ $7
        OR c.description @@@ $7
        OR c.summary @@@ $7
        OR c.content @@@ $7
        OR c.metadata @@@ $7
    )
GROUP BY c.id,
         t.total_count
ORDER BY c.created_at DESC
LIMIT $2 OFFSET $3
`

type SearchContentsWithFilterParams struct {
	UserID  uuid.UUID
	Limit   int32
	Offset  int32
	Domains []string
	Types   []string
	Tags    []string
	Query   pgtype.Text
}

type SearchContentsWithFilterRow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        string
	Title       string
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
	IsFavorite  pgtype.Bool
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	TotalCount  int64
	Tags        interface{}
}

func (q *Queries) SearchContentsWithFilter(ctx context.Context, db DBTX, arg SearchContentsWithFilterParams) ([]SearchContentsWithFilterRow, error) {
	rows, err := db.Query(ctx, searchContentsWithFilter,
		arg.UserID,
		arg.Limit,
		arg.Offset,
		arg.Domains,
		arg.Types,
		arg.Tags,
		arg.Query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchContentsWithFilterRow
	for rows.Next() {
		var i SearchContentsWithFilterRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.Domain,
			&i.S3Key,
			&i.Summary,
			&i.Content,
			&i.Html,
			&i.Metadata,
			&i.IsFavorite,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TotalCount,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unLinkContentWithTags = `-- name: UnLinkContentWithTags :exec
DELETE FROM content_tags_mapping
WHERE content_id = $1
  AND tag_id IN (SELECT id
                 FROM content_tags
                 WHERE name = ANY ($2 :: text[])
                   AND user_id = $3)
`

type UnLinkContentWithTagsParams struct {
	ContentID uuid.UUID
	Column2   []string
	UserID    uuid.UUID
}

// $1: content_id, $2: text[], $3: user_id
func (q *Queries) UnLinkContentWithTags(ctx context.Context, db DBTX, arg UnLinkContentWithTagsParams) error {
	_, err := db.Exec(ctx, unLinkContentWithTags, arg.ContentID, arg.Column2, arg.UserID)
	return err
}

const updateContent = `-- name: UpdateContent :one
UPDATE
    content
SET title       = COALESCE($3, title),
    description = COALESCE($4, description),
    url         = COALESCE($5, url),
    domain      = COALESCE($6, domain),
    s3_key      = COALESCE($7, s3_key),
    summary     = COALESCE($8, summary),
    content     = COALESCE($9, content),
    html        = COALESCE($10, html),
    metadata    = COALESCE($11, metadata),
    is_favorite = COALESCE($12, is_favorite)
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, type, title, description, url, domain, s3_key, summary, content, html, metadata, is_favorite, created_at, updated_at
`

type UpdateContentParams struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       pgtype.Text
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
	IsFavorite  pgtype.Bool
}

func (q *Queries) UpdateContent(ctx context.Context, db DBTX, arg UpdateContentParams) (Content, error) {
	row := db.QueryRow(ctx, updateContent,
		arg.ID,
		arg.UserID,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.Domain,
		arg.S3Key,
		arg.Summary,
		arg.Content,
		arg.Html,
		arg.Metadata,
		arg.IsFavorite,
	)
	var i Content
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Metadata,
		&i.IsFavorite,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateShareContent = `-- name: UpdateShareContent :one
UPDATE content_share cs
SET expires_at = $3
FROM content c
WHERE cs.content_id = c.id
  AND c.id = $1
  AND c.user_id = $2
RETURNING cs.id, cs.user_id, cs.content_id, cs.expires_at, cs.created_at, cs.updated_at
`

type UpdateShareContentParams struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt pgtype.Timestamptz
}

func (q *Queries) UpdateShareContent(ctx context.Context, db DBTX, arg UpdateShareContentParams) (ContentShare, error) {
	row := db.QueryRow(ctx, updateShareContent, arg.ID, arg.UserID, arg.ExpiresAt)
	var i ContentShare
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
