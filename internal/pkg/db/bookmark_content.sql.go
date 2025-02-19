// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: bookmark_content.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createBookmarkContent = `-- name: CreateBookmarkContent :one
INSERT INTO bookmark_content (
  type,
  url,
  user_id,
  title,
  description,
  domain,
  s3_key,
  summary,
  content,
  html,
  tags,
  metadata
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING id, type, url, user_id, title, description, domain, s3_key, summary, content, html, tags, metadata, created_at, updated_at
`

type CreateBookmarkContentParams struct {
	Type        string
	Url         string
	UserID      pgtype.UUID
	Title       pgtype.Text
	Description pgtype.Text
	Domain      pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Tags        []string
	Metadata    []byte
}

func (q *Queries) CreateBookmarkContent(ctx context.Context, db DBTX, arg CreateBookmarkContentParams) (BookmarkContent, error) {
	row := db.QueryRow(ctx, createBookmarkContent,
		arg.Type,
		arg.Url,
		arg.UserID,
		arg.Title,
		arg.Description,
		arg.Domain,
		arg.S3Key,
		arg.Summary,
		arg.Content,
		arg.Html,
		arg.Tags,
		arg.Metadata,
	)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getBookmarkContentByBookmarkID = `-- name: GetBookmarkContentByBookmarkID :one
SELECT bc.id, bc.type, bc.url, bc.user_id, bc.title, bc.description, bc.domain, bc.s3_key, bc.summary, bc.content, bc.html, bc.tags, bc.metadata, bc.created_at, bc.updated_at
FROM bookmarks b 
  JOIN bookmark_content bc ON b.content_id = bc.id
WHERE b.id = $1
`

func (q *Queries) GetBookmarkContentByBookmarkID(ctx context.Context, db DBTX, id uuid.UUID) (BookmarkContent, error) {
	row := db.QueryRow(ctx, getBookmarkContentByBookmarkID, id)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getBookmarkContentByID = `-- name: GetBookmarkContentByID :one
SELECT id, type, url, user_id, title, description, domain, s3_key, summary, content, html, tags, metadata, created_at, updated_at
FROM bookmark_content
WHERE id = $1
`

func (q *Queries) GetBookmarkContentByID(ctx context.Context, db DBTX, id uuid.UUID) (BookmarkContent, error) {
	row := db.QueryRow(ctx, getBookmarkContentByID, id)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getBookmarkContentByURL = `-- name: GetBookmarkContentByURL :one
SELECT id, type, url, user_id, title, description, domain, s3_key, summary, content, html, tags, metadata, created_at, updated_at
FROM bookmark_content
WHERE url = $1 AND (user_id = $2 OR user_id IS NULL)
LIMIT 1
`

type GetBookmarkContentByURLParams struct {
	Url    string
	UserID pgtype.UUID
}

// First try to get user specific content, then the shared content
func (q *Queries) GetBookmarkContentByURL(ctx context.Context, db DBTX, arg GetBookmarkContentByURLParams) (BookmarkContent, error) {
	row := db.QueryRow(ctx, getBookmarkContentByURL, arg.Url, arg.UserID)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const isBookmarkContentExistByURL = `-- name: IsBookmarkContentExistByURL :one
SELECT EXISTS (
  SELECT 1
  FROM bookmark_content
  WHERE url = $1
)
`

func (q *Queries) IsBookmarkContentExistByURL(ctx context.Context, db DBTX, url string) (bool, error) {
	row := db.QueryRow(ctx, isBookmarkContentExistByURL, url)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const ownerTransferBookmarkContent = `-- name: OwnerTransferBookmarkContent :exec
UPDATE bookmark_content
SET 
    user_id = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $2
`

type OwnerTransferBookmarkContentParams struct {
	NewUserID pgtype.UUID
	UserID    pgtype.UUID
}

func (q *Queries) OwnerTransferBookmarkContent(ctx context.Context, db DBTX, arg OwnerTransferBookmarkContentParams) error {
	_, err := db.Exec(ctx, ownerTransferBookmarkContent, arg.NewUserID, arg.UserID)
	return err
}

const updateBookmarkContent = `-- name: UpdateBookmarkContent :one
UPDATE bookmark_content
SET title = COALESCE($2, title),
    description = COALESCE($3, description),
    s3_key = COALESCE($4, s3_key),
    summary = COALESCE($5, summary),
    content = COALESCE($6, content),
    html = COALESCE($7, html),
    tags = COALESCE($8, tags),
    metadata = COALESCE($9, metadata)
WHERE id = $1
RETURNING id, type, url, user_id, title, description, domain, s3_key, summary, content, html, tags, metadata, created_at, updated_at
`

type UpdateBookmarkContentParams struct {
	ID          uuid.UUID
	Title       pgtype.Text
	Description pgtype.Text
	S3Key       pgtype.Text
	Summary     pgtype.Text
	Content     pgtype.Text
	Html        pgtype.Text
	Tags        []string
	Metadata    []byte
}

func (q *Queries) UpdateBookmarkContent(ctx context.Context, db DBTX, arg UpdateBookmarkContentParams) (BookmarkContent, error) {
	row := db.QueryRow(ctx, updateBookmarkContent,
		arg.ID,
		arg.Title,
		arg.Description,
		arg.S3Key,
		arg.Summary,
		arg.Content,
		arg.Html,
		arg.Tags,
		arg.Metadata,
	)
	var i BookmarkContent
	err := row.Scan(
		&i.ID,
		&i.Type,
		&i.Url,
		&i.UserID,
		&i.Title,
		&i.Description,
		&i.Domain,
		&i.S3Key,
		&i.Summary,
		&i.Content,
		&i.Html,
		&i.Tags,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
