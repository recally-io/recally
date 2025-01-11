// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: assistants_thread.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createAssistantThread = `-- name: CreateAssistantThread :one
INSERT INTO assistant_threads (uuid, user_id, assistant_id, name, description, system_prompt, model, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, uuid, user_id, assistant_id, name, description, system_prompt, model, metadata, created_at, updated_at
`

type CreateAssistantThreadParams struct {
	Uuid         uuid.UUID
	UserID       pgtype.UUID
	AssistantID  pgtype.UUID
	Name         string
	Description  pgtype.Text
	SystemPrompt pgtype.Text
	Model        string
	Metadata     []byte
}

// CRUD for assistant_threads
func (q *Queries) CreateAssistantThread(ctx context.Context, db DBTX, arg CreateAssistantThreadParams) (AssistantThread, error) {
	row := db.QueryRow(ctx, createAssistantThread,
		arg.Uuid,
		arg.UserID,
		arg.AssistantID,
		arg.Name,
		arg.Description,
		arg.SystemPrompt,
		arg.Model,
		arg.Metadata,
	)
	var i AssistantThread
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.AssistantID,
		&i.Name,
		&i.Description,
		&i.SystemPrompt,
		&i.Model,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAssistantThread = `-- name: DeleteAssistantThread :exec
DELETE FROM assistant_threads WHERE uuid = $1
`

func (q *Queries) DeleteAssistantThread(ctx context.Context, db DBTX, argUuid uuid.UUID) error {
	_, err := db.Exec(ctx, deleteAssistantThread, argUuid)
	return err
}

const deleteAssistantThreadsByAssistant = `-- name: DeleteAssistantThreadsByAssistant :exec
DELETE FROM assistant_threads WHERE assistant_id = $1
`

func (q *Queries) DeleteAssistantThreadsByAssistant(ctx context.Context, db DBTX, assistantID pgtype.UUID) error {
	_, err := db.Exec(ctx, deleteAssistantThreadsByAssistant, assistantID)
	return err
}

const getAssistantThread = `-- name: GetAssistantThread :one
SELECT id, uuid, user_id, assistant_id, name, description, system_prompt, model, metadata, created_at, updated_at FROM assistant_threads WHERE uuid = $1
`

func (q *Queries) GetAssistantThread(ctx context.Context, db DBTX, argUuid uuid.UUID) (AssistantThread, error) {
	row := db.QueryRow(ctx, getAssistantThread, argUuid)
	var i AssistantThread
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.AssistantID,
		&i.Name,
		&i.Description,
		&i.SystemPrompt,
		&i.Model,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAssistantThreads = `-- name: ListAssistantThreads :many
SELECT id, uuid, user_id, assistant_id, name, description, system_prompt, model, metadata, created_at, updated_at FROM assistant_threads WHERE assistant_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAssistantThreads(ctx context.Context, db DBTX, assistantID pgtype.UUID) ([]AssistantThread, error) {
	rows, err := db.Query(ctx, listAssistantThreads, assistantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssistantThread
	for rows.Next() {
		var i AssistantThread
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
			&i.AssistantID,
			&i.Name,
			&i.Description,
			&i.SystemPrompt,
			&i.Model,
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

const listAssistantThreadsByUser = `-- name: ListAssistantThreadsByUser :many
SELECT id, uuid, user_id, assistant_id, name, description, system_prompt, model, metadata, created_at, updated_at FROM assistant_threads WHERE user_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAssistantThreadsByUser(ctx context.Context, db DBTX, userID pgtype.UUID) ([]AssistantThread, error) {
	rows, err := db.Query(ctx, listAssistantThreadsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssistantThread
	for rows.Next() {
		var i AssistantThread
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
			&i.AssistantID,
			&i.Name,
			&i.Description,
			&i.SystemPrompt,
			&i.Model,
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

const updateAssistantThread = `-- name: UpdateAssistantThread :one
UPDATE assistant_threads SET name = $2, description = $3, model = $4, metadata = $5, system_prompt = $6
WHERE uuid = $1
RETURNING id, uuid, user_id, assistant_id, name, description, system_prompt, model, metadata, created_at, updated_at
`

type UpdateAssistantThreadParams struct {
	Uuid         uuid.UUID
	Name         string
	Description  pgtype.Text
	Model        string
	Metadata     []byte
	SystemPrompt pgtype.Text
}

func (q *Queries) UpdateAssistantThread(ctx context.Context, db DBTX, arg UpdateAssistantThreadParams) (AssistantThread, error) {
	row := db.QueryRow(ctx, updateAssistantThread,
		arg.Uuid,
		arg.Name,
		arg.Description,
		arg.Model,
		arg.Metadata,
		arg.SystemPrompt,
	)
	var i AssistantThread
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.AssistantID,
		&i.Name,
		&i.Description,
		&i.SystemPrompt,
		&i.Model,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
