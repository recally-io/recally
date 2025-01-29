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

const deleteBookmarksByUser = `-- name: DeleteBookmarksByUser :exec
DELETE FROM bookmarks
WHERE user_id = $1
`

func (q *Queries) DeleteBookmarksByUser(ctx context.Context, db DBTX, userID pgtype.UUID) error {
	_, err := db.Exec(ctx, deleteBookmarksByUser, userID)
	return err
}

const getBookmarkWithContent = `-- name: GetBookmarkWithContent :one
SELECT b.id, b.user_id, b.content_id, b.is_favorite, b.is_archive, b.is_public, b.reading_progress, b.metadata, b.created_at, b.updated_at,
       bc.id, bc.type, bc.url, bc.user_id, bc.title, bc.description, bc.domain, bc.s3_key, bc.summary, bc.content, bc.html, bc.tags, bc.metadata, bc.created_at, bc.updated_at,
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
LIMIT 1
`

type GetBookmarkWithContentParams struct {
	ID     uuid.UUID
	UserID pgtype.UUID
}

type GetBookmarkWithContentRow struct {
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
	Url             string
	UserID_2        pgtype.UUID
	Title           pgtype.Text
	Description     pgtype.Text
	Domain          pgtype.Text
	S3Key           pgtype.Text
	Summary         pgtype.Text
	Content         pgtype.Text
	Html            pgtype.Text
	Tags            []string
	Metadata_2      []byte
	CreatedAt_2     pgtype.Timestamptz
	UpdatedAt_2     pgtype.Timestamptz
	Tags_2          interface{}
}

func (q *Queries) GetBookmarkWithContent(ctx context.Context, db DBTX, arg GetBookmarkWithContentParams) (GetBookmarkWithContentRow, error) {
	row := db.QueryRow(ctx, getBookmarkWithContent, arg.ID, arg.UserID)
	var i GetBookmarkWithContentRow
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
		&i.Url,
		&i.UserID_2,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata_2,
		&i.CreatedAt_2,
		&i.UpdatedAt_2,
		&i.Tags_2,
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
	Url    string
	UserID pgtype.UUID
}

func (q *Queries) IsBookmarkExistWithURL(ctx context.Context, db DBTX, arg IsBookmarkExistWithURLParams) (bool, error) {
	row := db.QueryRow(ctx, isBookmarkExistWithURL, arg.Url, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const listBookmarkDomainsByUser = `-- name: ListBookmarkDomainsByUser :many
SELECT bc.domain, count(*) as cnt
FROM bookmarks b
  JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.user_id = $1 
AND bc.domain IS NOT NULL
GROUP BY bc.domain
ORDER BY cnt DESC, domain ASC
`

type ListBookmarkDomainsByUserRow struct {
	Domain pgtype.Text
	Cnt    int64
}

func (q *Queries) ListBookmarkDomainsByUser(ctx context.Context, db DBTX, userID pgtype.UUID) ([]ListBookmarkDomainsByUserRow, error) {
	rows, err := db.Query(ctx, listBookmarkDomainsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListBookmarkDomainsByUserRow
	for rows.Next() {
		var i ListBookmarkDomainsByUserRow
		if err := rows.Scan(&i.Domain, &i.Cnt); err != nil {
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
           LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
  WHERE b.user_id = $1
    AND ($4::text[] IS NULL OR bc.domain = ANY($4::text[]))
    AND ($5::text[] IS NULL OR bc.type = ANY($5::text[]))
    AND ($6::text[] IS NULL OR bct.name = ANY($6::text[]))
)
SELECT b.id, b.user_id, b.content_id, b.is_favorite, b.is_archive, b.is_public, b.reading_progress, b.metadata, b.created_at, b.updated_at,
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

const ownerTransferBookmark = `-- name: OwnerTransferBookmark :exec
UPDATE bookmarks 
SET 
    user_id = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
`

type OwnerTransferBookmarkParams struct {
	UserID   pgtype.UUID
	UserID_2 pgtype.UUID
}

func (q *Queries) OwnerTransferBookmark(ctx context.Context, db DBTX, arg OwnerTransferBookmarkParams) error {
	_, err := db.Exec(ctx, ownerTransferBookmark, arg.UserID, arg.UserID_2)
	return err
}

const searchBookmarks = `-- name: SearchBookmarks :many
WITH total AS (
  SELECT COUNT(DISTINCT b.*) AS total_count
  FROM bookmarks AS b
           JOIN bookmark_content AS bc ON b.content_id = bc.id
           LEFT JOIN bookmark_tags_mapping AS bctm ON bc.id = bctm.content_id
           LEFT JOIN bookmark_tags AS bct ON bctm.tag_id = bct.id
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
