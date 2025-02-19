// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: files.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createFile = `-- name: CreateFile :one
INSERT INTO files (
    original_url,
    user_id,
    s3_key,
    s3_url,
    file_name,
    file_type,
    file_size,
    file_hash,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at
`

type CreateFileParams struct {
	OriginalUrl string
	UserID      uuid.UUID
	S3Key       string
	S3Url       pgtype.Text
	FileName    pgtype.Text
	FileType    string
	FileSize    pgtype.Int8
	FileHash    pgtype.Text
	Metadata    []byte
}

func (q *Queries) CreateFile(ctx context.Context, db DBTX, arg CreateFileParams) (File, error) {
	row := db.QueryRow(ctx, createFile,
		arg.OriginalUrl,
		arg.UserID,
		arg.S3Key,
		arg.S3Url,
		arg.FileName,
		arg.FileType,
		arg.FileSize,
		arg.FileHash,
		arg.Metadata,
	)
	var i File
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OriginalUrl,
		&i.S3Key,
		&i.S3Url,
		&i.FileName,
		&i.FileType,
		&i.FileSize,
		&i.FileHash,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFile = `-- name: DeleteFile :exec
DELETE FROM files
WHERE id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, db DBTX, id uuid.UUID) error {
	_, err := db.Exec(ctx, deleteFile, id)
	return err
}

const deleteFileByOriginalURL = `-- name: DeleteFileByOriginalURL :exec
DELETE FROM files
WHERE original_url = $1
AND (user_id = $2 OR user_id = $3)
`

type DeleteFileByOriginalURLParams struct {
	OriginalUrl string
	UserID      uuid.UUID
	DummyUserID pgtype.UUID
}

func (q *Queries) DeleteFileByOriginalURL(ctx context.Context, db DBTX, arg DeleteFileByOriginalURLParams) error {
	_, err := db.Exec(ctx, deleteFileByOriginalURL, arg.OriginalUrl, arg.UserID, arg.DummyUserID)
	return err
}

const deleteFileByS3Key = `-- name: DeleteFileByS3Key :exec
DELETE FROM files
WHERE s3_key = $1
AND (user_id = $2 OR user_id = $3)
`

type DeleteFileByS3KeyParams struct {
	S3Key       string
	UserID      uuid.UUID
	DummyUserID pgtype.UUID
}

func (q *Queries) DeleteFileByS3Key(ctx context.Context, db DBTX, arg DeleteFileByS3KeyParams) error {
	_, err := db.Exec(ctx, deleteFileByS3Key, arg.S3Key, arg.UserID, arg.DummyUserID)
	return err
}

const getFileByID = `-- name: GetFileByID :one
SELECT id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at FROM files
WHERE id = $1
`

func (q *Queries) GetFileByID(ctx context.Context, db DBTX, id uuid.UUID) (File, error) {
	row := db.QueryRow(ctx, getFileByID, id)
	var i File
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OriginalUrl,
		&i.S3Key,
		&i.S3Url,
		&i.FileName,
		&i.FileType,
		&i.FileSize,
		&i.FileHash,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFileByOriginalURL = `-- name: GetFileByOriginalURL :one
SELECT id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at FROM files
WHERE original_url = $1 
AND (user_id = $2 OR user_id = $3)
`

type GetFileByOriginalURLParams struct {
	OriginalUrl string
	UserID      uuid.UUID
	DummyUserID pgtype.UUID
}

func (q *Queries) GetFileByOriginalURL(ctx context.Context, db DBTX, arg GetFileByOriginalURLParams) (File, error) {
	row := db.QueryRow(ctx, getFileByOriginalURL, arg.OriginalUrl, arg.UserID, arg.DummyUserID)
	var i File
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OriginalUrl,
		&i.S3Key,
		&i.S3Url,
		&i.FileName,
		&i.FileType,
		&i.FileSize,
		&i.FileHash,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getFileByS3Key = `-- name: GetFileByS3Key :one
SELECT id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at FROM files
WHERE s3_key = $1
AND (user_id = $2 OR user_id = $3)
`

type GetFileByS3KeyParams struct {
	S3Key       string
	UserID      uuid.UUID
	DummyUserID pgtype.UUID
}

func (q *Queries) GetFileByS3Key(ctx context.Context, db DBTX, arg GetFileByS3KeyParams) (File, error) {
	row := db.QueryRow(ctx, getFileByS3Key, arg.S3Key, arg.UserID, arg.DummyUserID)
	var i File
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OriginalUrl,
		&i.S3Key,
		&i.S3Url,
		&i.FileName,
		&i.FileType,
		&i.FileSize,
		&i.FileHash,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listFiles = `-- name: ListFiles :many
SELECT id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at FROM files
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3
`

type ListFilesParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) ListFiles(ctx context.Context, db DBTX, arg ListFilesParams) ([]File, error) {
	rows, err := db.Query(ctx, listFiles, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.OriginalUrl,
			&i.S3Key,
			&i.S3Url,
			&i.FileName,
			&i.FileType,
			&i.FileSize,
			&i.FileHash,
			&i.Metadata,
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

const searchFilesByType = `-- name: SearchFilesByType :many
SELECT id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at FROM files
WHERE file_type = $1
AND user_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
`

type SearchFilesByTypeParams struct {
	FileType string
	UserID   uuid.UUID
	Limit    int32
	Offset   int32
}

func (q *Queries) SearchFilesByType(ctx context.Context, db DBTX, arg SearchFilesByTypeParams) ([]File, error) {
	rows, err := db.Query(ctx, searchFilesByType,
		arg.FileType,
		arg.UserID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.OriginalUrl,
			&i.S3Key,
			&i.S3Url,
			&i.FileName,
			&i.FileType,
			&i.FileSize,
			&i.FileHash,
			&i.Metadata,
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

const updateFile = `-- name: UpdateFile :one
UPDATE files
SET 
    s3_url = COALESCE($2, s3_url),
    file_name = COALESCE($3, file_name),
    file_type = COALESCE($4, file_type),
    file_size = COALESCE($5, file_size),
    metadata = COALESCE($6, metadata)
WHERE id = $1
RETURNING id, user_id, original_url, s3_key, s3_url, file_name, file_type, file_size, file_hash, metadata, created_at, updated_at
`

type UpdateFileParams struct {
	ID       uuid.UUID
	S3Url    pgtype.Text
	FileName pgtype.Text
	FileType string
	FileSize pgtype.Int8
	Metadata []byte
}

func (q *Queries) UpdateFile(ctx context.Context, db DBTX, arg UpdateFileParams) (File, error) {
	row := db.QueryRow(ctx, updateFile,
		arg.ID,
		arg.S3Url,
		arg.FileName,
		arg.FileType,
		arg.FileSize,
		arg.Metadata,
	)
	var i File
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OriginalUrl,
		&i.S3Key,
		&i.S3Url,
		&i.FileName,
		&i.FileType,
		&i.FileSize,
		&i.FileHash,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
