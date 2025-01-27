// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bookmarks.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createBookmark = `-- name: CreateBookmark :one
INSERT INTO bookmarks (
  user_id, content_id, is_favorite, is_archive,
  is_public, reading_progress, metadata
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, content_id, is_favorite, is_archive, is_public, reading_progress, metadata, created_at, updated_at
`

type CreateBookmarkParams struct {
	UserID          pgtype.UUID
	ContentID       pgtype.UUID
	IsFavorite      pgtype.Bool
	IsArchive       pgtype.Bool
	IsPublic        pgtype.Bool
	ReadingProgress pgtype.Int4
	Metadata        []byte
}

func (q *Queries) CreateBookmark(ctx context.Context, db DBTX, arg CreateBookmarkParams) (Bookmark, error) {
	row := db.QueryRow(ctx, createBookmark,
		arg.UserID,
		arg.ContentID,
		arg.IsFavorite,
		arg.IsArchive,
		arg.IsPublic,
		arg.ReadingProgress,
		arg.Metadata,
	)
	var i Bookmark
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.IsFavorite,
		&i.IsArchive,
		&i.IsPublic,
		&i.ReadingProgress,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createBookmarkContent = `-- name: CreateBookmarkContent :one
INSERT INTO bookmark_content (
  type, title, description, url, domain, s3_key,
  summary, content, html, metadata
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, type, title, description, url, domain, s3_key, summary, content, html, metadata, created_at, updated_at
`

type CreateBookmarkContentParams struct {
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
}

func (q *Queries) CreateBookmarkContent(ctx context.Context, db DBTX, arg CreateBookmarkContentParams) (BookmarkContent, error) {
	row := db.QueryRow(ctx, createBookmarkContent,
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
	)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
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
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createBookmarkContentTag = `-- name: CreateBookmarkContentTag :one
INSERT INTO bookmark_content_tags (name, user_id)
VALUES ($1, $2)
ON CONFLICT (name, user_id) DO UPDATE
    SET usage_count = bookmark_content_tags.usage_count + 1
RETURNING id, name, user_id, created_at, updated_at
`

type CreateBookmarkContentTagParams struct {
	Name   string
	UserID uuid.UUID
}

func (q *Queries) CreateBookmarkContentTag(ctx context.Context, db DBTX, arg CreateBookmarkContentTagParams) (BookmarkContentTag, error) {
	row := db.QueryRow(ctx, createBookmarkContentTag, arg.Name, arg.UserID)
	var i BookmarkContentTag
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteBookmark = `-- name: DeleteBookmark :exec
DELETE FROM bookmarks
WHERE id = $1 AND user_id = $2
`

type DeleteBookmarkParams struct {
	ID     uuid.UUID
	UserID pgtype.UUID
}

func (q *Queries) DeleteBookmark(ctx context.Context, db DBTX, arg DeleteBookmarkParams) error {
	_, err := db.Exec(ctx, deleteBookmark, arg.ID, arg.UserID)
	return err
}

const deleteBookmarkContentTag = `-- name: DeleteBookmarkContentTag :exec
DELETE FROM bookmark_content_tags
WHERE id = $1
  AND user_id = $2
`

type DeleteBookmarkContentTagParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteBookmarkContentTag(ctx context.Context, db DBTX, arg DeleteBookmarkContentTagParams) error {
	_, err := db.Exec(ctx, deleteBookmarkContentTag, arg.ID, arg.UserID)
	return err
}

const deleteBookmarksByUser = `-- name: DeleteBookmarksByUser :exec
DELETE FROM bookmarks
WHERE user_id = $1
`

func (q *Queries) DeleteBookmarksByUser(ctx context.Context, db DBTX, userID pgtype.UUID) error {
	_, err := db.Exec(ctx, deleteBookmarksByUser, userID)
	return err
}

const getBookmark = `-- name: GetBookmark :one
SELECT b.id, b.user_id, b.content_id, b.is_favorite, b.is_archive, b.is_public, b.reading_progress, b.metadata, b.created_at, b.updated_at,
       bc.id, bc.type, bc.title, bc.description, bc.url, bc.domain, bc.s3_key, bc.summary, bc.content, bc.html, bc.metadata, bc.created_at, bc.updated_at,
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
LIMIT 1
`

type GetBookmarkParams struct {
	ID     uuid.UUID
	UserID pgtype.UUID
}

type GetBookmarkRow struct {
	ID              uuid.UUID
	UserID          pgtype.UUID
	ContentID       pgtype.UUID
	IsFavorite      pgtype.Bool
	IsArchive       pgtype.Bool
	IsPublic        pgtype.Bool
	ReadingProgress pgtype.Int4
	Metadata        []byte
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	ID_2            uuid.UUID
	Type            string
	Title           string
	Description     pgtype.Text
	Url             pgtype.Text
	Domain          pgtype.Text
	S3Key           pgtype.Text
	Summary         pgtype.Text
	Content         pgtype.Text
	Html            pgtype.Text
	Metadata_2      []byte
	CreatedAt_2     pgtype.Timestamptz
	UpdatedAt_2     pgtype.Timestamptz
	Tags            interface{}
}

func (q *Queries) GetBookmark(ctx context.Context, db DBTX, arg GetBookmarkParams) (GetBookmarkRow, error) {
	row := db.QueryRow(ctx, getBookmark, arg.ID, arg.UserID)
	var i GetBookmarkRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.IsFavorite,
		&i.IsArchive,
		&i.IsPublic,
		&i.ReadingProgress,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ID_2,
		&i.Type,
		&i.Title,
		&i.Description,
		&i.Url,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Metadata_2,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
		&i.Tags,
	)
	return i, err
}

const isBookmarkExistWithURL = `-- name: IsBookmarkExistWithURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmarks b
           JOIN bookmark_content bc ON b.content_id = bc.id
  WHERE bc.url = $1
    AND b.user_id = $2
)
`

type IsBookmarkExistWithURLParams struct {
	Url    pgtype.Text
	UserID pgtype.UUID
}

func (q *Queries) IsBookmarkExistWithURL(ctx context.Context, db DBTX, arg IsBookmarkExistWithURLParams) (bool, error) {
	row := db.QueryRow(ctx, isBookmarkExistWithURL, arg.Url, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const linkBookmarkContentWithTags = `-- name: LinkBookmarkContentWithTags :exec
INSERT INTO bookmark_content_tags_mapping (content_id, tag_id)
SELECT $1, bct.id
FROM bookmark_content_tags bct
WHERE bct.name = ANY ($2::text[])
  AND bct.user_id = $3
`

type LinkBookmarkContentWithTagsParams struct {
	ContentID uuid.UUID
	Column2   []string
	UserID    uuid.UUID
}

func (q *Queries) LinkBookmarkContentWithTags(ctx context.Context, db DBTX, arg LinkBookmarkContentWithTagsParams) error {
	_, err := db.Exec(ctx, linkBookmarkContentWithTags, arg.ContentID, arg.Column2, arg.UserID)
	return err
}

const listBookmarkContentTags = `-- name: ListBookmarkContentTags :many
SELECT bct.name
FROM bookmark_content_tags bct
         JOIN bookmark_content_tags_mapping bctm ON bct.id = bctm.tag_id
WHERE bctm.content_id = $1
  AND bct.user_id = $2
`

type ListBookmarkContentTagsParams struct {
	ContentID uuid.UUID
	UserID    uuid.UUID
}

func (q *Queries) ListBookmarkContentTags(ctx context.Context, db DBTX, arg ListBookmarkContentTagsParams) ([]string, error) {
	rows, err := db.Query(ctx, listBookmarkContentTags, arg.ContentID, arg.UserID)
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

const listBookmarkDomains = `-- name: ListBookmarkDomains :many
SELECT bc.domain, count(*) as count
FROM bookmarks b
         JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.user_id = $1 
AND bc.domain IS NOT NULL
GROUP BY bc.domain
ORDER BY count DESC, domain ASC
`

type ListBookmarkDomainsRow struct {
	Domain pgtype.Text
	Count  int64
}

func (q *Queries) ListBookmarkDomains(ctx context.Context, db DBTX, userID pgtype.UUID) ([]ListBookmarkDomainsRow, error) {
	rows, err := db.Query(ctx, listBookmarkDomains, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBookmarkDomainsRow
	for rows.Next() {
		var i ListBookmarkDomainsRow
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

const listBookmarkTagsByUser = `-- name: ListBookmarkTagsByUser :many
SELECT bct.name, count(bctm.*) as count
FROM bookmark_content_tags bct
         JOIN bookmark_content_tags_mapping bctm ON bct.id = bctm.tag_id
WHERE bct.user_id = $1
GROUP BY bct.name
ORDER BY count DESC
`

type ListBookmarkTagsByUserRow struct {
	Name  string
	Count int64
}

// Tags related queries similar to content.sql but adapted for new schema
func (q *Queries) ListBookmarkTagsByUser(ctx context.Context, db DBTX, userID uuid.UUID) ([]ListBookmarkTagsByUserRow, error) {
	rows, err := db.Query(ctx, listBookmarkTagsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBookmarkTagsByUserRow
	for rows.Next() {
		var i ListBookmarkTagsByUserRow
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

const listBookmarks = `-- name: ListBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND ($4::text[] IS NULL OR bc.domain = ANY($4::text[]))
    AND ($5::text[] IS NULL OR bc.type = ANY($5::text[]))
    AND ($6::text[] IS NULL OR bct.name = ANY($6::text[]))
)
SELECT b.id, b.user_id, b.content_id, b.is_favorite, b.is_archive, b.is_public, b.reading_progress, b.metadata, b.created_at, b.updated_at,
       bc.id, bc.type, bc.title, bc.description, bc.url, bc.domain, bc.s3_key, bc.summary, bc.content, bc.html, bc.metadata, bc.created_at, bc.updated_at,
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
  AND ($4::text[] IS NULL OR bc.domain = ANY($4::text[]))
  AND ($5::text[] IS NULL OR bc.type = ANY($5::text[]))
  AND ($6::text[] IS NULL OR bct.name = ANY($6::text[]))
GROUP BY b.id, bc.id, t.total_count
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3
`

type ListBookmarksParams struct {
	UserID  pgtype.UUID
	Limit   int32
	Offset  int32
	Domains []string
	Types   []string
	Tags    []string
}

type ListBookmarksRow struct {
	ID              uuid.UUID
	UserID          pgtype.UUID
	ContentID       pgtype.UUID
	IsFavorite      pgtype.Bool
	IsArchive       pgtype.Bool
	IsPublic        pgtype.Bool
	ReadingProgress pgtype.Int4
	Metadata        []byte
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	ID_2            uuid.UUID
	Type            string
	Title           string
	Description     pgtype.Text
	Url             pgtype.Text
	Domain          pgtype.Text
	S3Key           pgtype.Text
	Summary         pgtype.Text
	Content         pgtype.Text
	Html            pgtype.Text
	Metadata_2      []byte
	CreatedAt_2     pgtype.Timestamptz
	UpdatedAt_2     pgtype.Timestamptz
	TotalCount      int64
	Tags            interface{}
}

func (q *Queries) ListBookmarks(ctx context.Context, db DBTX, arg ListBookmarksParams) ([]ListBookmarksRow, error) {
	rows, err := db.Query(ctx, listBookmarks,
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
	var items []ListBookmarksRow
	for rows.Next() {
		var i ListBookmarksRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ContentID,
			&i.IsFavorite,
			&i.IsArchive,
			&i.IsPublic,
			&i.ReadingProgress,
			&i.Metadata,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Type,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.Domain,
			&i.S3Key,
			&i.Summary,
			&i.Content,
			&i.Html,
			&i.Metadata_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const listExistingBookmarkTagsByTags = `-- name: ListExistingBookmarkTagsByTags :many
SELECT name
FROM bookmark_content_tags
WHERE name = ANY ($1::text[])
  AND user_id = $2
`

type ListExistingBookmarkTagsByTagsParams struct {
	Column1 []string
	UserID  uuid.UUID
}

func (q *Queries) ListExistingBookmarkTagsByTags(ctx context.Context, db DBTX, arg ListExistingBookmarkTagsByTagsParams) ([]string, error) {
	rows, err := db.Query(ctx, listExistingBookmarkTagsByTags, arg.Column1, arg.UserID)
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

const searchBookmarks = `-- name: SearchBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_content_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_content_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND ($4::text[] IS NULL OR bc.domain = ANY($4::text[]))
    AND ($5::text[] IS NULL OR bc.type = ANY($5::text[]))
    AND ($6::text[] IS NULL OR bct.name = ANY($6::text[]))
    AND (
      $7::text IS NULL
      OR bc.title @@@ $7
      OR bc.description @@@ $7
      OR bc.summary @@@ $7
      OR bc.content @@@ $7
      OR bc.metadata @@@ $7
    )
)
SELECT b.id, b.user_id, b.content_id, b.is_favorite, b.is_archive, b.is_public, b.reading_progress, b.metadata, b.created_at, b.updated_at,
       bc.id, bc.type, bc.title, bc.description, bc.url, bc.domain, bc.s3_key, bc.summary, bc.content, bc.html, bc.metadata, bc.created_at, bc.updated_at,
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
  AND ($4::text[] IS NULL OR bc.domain = ANY($4::text[]))
  AND ($5::text[] IS NULL OR bc.type = ANY($5::text[]))
  AND ($6::text[] IS NULL OR bct.name = ANY($6::text[]))
  AND (
    $7::text IS NULL
    OR bc.title @@@ $7
    OR bc.description @@@ $7
    OR bc.summary @@@ $7
    OR bc.content @@@ $7
    OR bc.metadata @@@ $7
  )
GROUP BY b.id, bc.id, t.total_count
ORDER BY b.created_at DESC
LIMIT $2 OFFSET $3
`

type SearchBookmarksParams struct {
	UserID  pgtype.UUID
	Limit   int32
	Offset  int32
	Domains []string
	Types   []string
	Tags    []string
	Query   pgtype.Text
}

type SearchBookmarksRow struct {
	ID              uuid.UUID
	UserID          pgtype.UUID
	ContentID       pgtype.UUID
	IsFavorite      pgtype.Bool
	IsArchive       pgtype.Bool
	IsPublic        pgtype.Bool
	ReadingProgress pgtype.Int4
	Metadata        []byte
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
	ID_2            uuid.UUID
	Type            string
	Title           string
	Description     pgtype.Text
	Url             pgtype.Text
	Domain          pgtype.Text
	S3Key           pgtype.Text
	Summary         pgtype.Text
	Content         pgtype.Text
	Html            pgtype.Text
	Metadata_2      []byte
	CreatedAt_2     pgtype.Timestamptz
	UpdatedAt_2     pgtype.Timestamptz
	TotalCount      int64
	Tags            interface{}
}

func (q *Queries) SearchBookmarks(ctx context.Context, db DBTX, arg SearchBookmarksParams) ([]SearchBookmarksRow, error) {
	rows, err := db.Query(ctx, searchBookmarks,
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
	var items []SearchBookmarksRow
	for rows.Next() {
		var i SearchBookmarksRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ContentID,
			&i.IsFavorite,
			&i.IsArchive,
			&i.IsPublic,
			&i.ReadingProgress,
			&i.Metadata,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ID_2,
			&i.Type,
			&i.Title,
			&i.Description,
			&i.Url,
			&i.Domain,
			&i.S3Key,
			&i.Summary,
			&i.Content,
			&i.Html,
			&i.Metadata_2,
			&i.CreatedAt_2,
			&i.UpdatedAt_2,
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

const unLinkBookmarkContentWithTags = `-- name: UnLinkBookmarkContentWithTags :exec
DELETE FROM bookmark_content_tags_mapping
WHERE content_id = $1
  AND tag_id IN (SELECT id
                 FROM bookmark_content_tags
                 WHERE name = ANY ($2::text[])
                   AND user_id = $3)
`

type UnLinkBookmarkContentWithTagsParams struct {
	ContentID uuid.UUID
	Column2   []string
	UserID    uuid.UUID
}

func (q *Queries) UnLinkBookmarkContentWithTags(ctx context.Context, db DBTX, arg UnLinkBookmarkContentWithTagsParams) error {
	_, err := db.Exec(ctx, unLinkBookmarkContentWithTags, arg.ContentID, arg.Column2, arg.UserID)
	return err
}

const updateBookmark = `-- name: UpdateBookmark :one
UPDATE bookmarks
SET is_favorite = COALESCE($3, is_favorite),
    is_archive = COALESCE($4, is_archive),
    is_public = COALESCE($5, is_public),
    reading_progress = COALESCE($6, reading_progress),
    metadata = COALESCE($7, metadata)
WHERE id = $1
  AND user_id = $2
RETURNING id, user_id, content_id, is_favorite, is_archive, is_public, reading_progress, metadata, created_at, updated_at
`

type UpdateBookmarkParams struct {
	ID              uuid.UUID
	UserID          pgtype.UUID
	IsFavorite      pgtype.Bool
	IsArchive       pgtype.Bool
	IsPublic        pgtype.Bool
	ReadingProgress pgtype.Int4
	Metadata        []byte
}

func (q *Queries) UpdateBookmark(ctx context.Context, db DBTX, arg UpdateBookmarkParams) (Bookmark, error) {
	row := db.QueryRow(ctx, updateBookmark,
		arg.ID,
		arg.UserID,
		arg.IsFavorite,
		arg.IsArchive,
		arg.IsPublic,
		arg.ReadingProgress,
		arg.Metadata,
	)
	var i Bookmark
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ContentID,
		&i.IsFavorite,
		&i.IsArchive,
		&i.IsPublic,
		&i.ReadingProgress,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateBookmarkContent = `-- name: UpdateBookmarkContent :one
UPDATE bookmark_content
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    url = COALESCE($4, url),
    domain = COALESCE($5, domain),
    s3_key = COALESCE($6, s3_key),
    summary = COALESCE($7, summary),
    content = COALESCE($8, content),
    html = COALESCE($9, html),
    metadata = COALESCE($10, metadata)
WHERE id = $1
RETURNING id, type, title, description, url, domain, s3_key, summary, content, html, metadata, created_at, updated_at
`

type UpdateBookmarkContentParams struct {
	ID          uuid.UUID
	Title       pgtype.Text
	Description pgtype.Text
	Url         pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Metadata    []byte
}

func (q *Queries) UpdateBookmarkContent(ctx context.Context, db DBTX, arg UpdateBookmarkContentParams) (BookmarkContent, error) {
	row := db.QueryRow(ctx, updateBookmarkContent,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.Url,
		arg.Domain,
		arg.S3Key,
		arg.Summary,
		arg.Content,
		arg.Html,
		arg.Metadata,
	)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
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
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
