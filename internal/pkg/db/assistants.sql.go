// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: assistants.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

const createAssistant = `-- name: CreateAssistant :exec

INSERT INTO assistants (user_id, name, description, system_prompt, model, metadata)
VALUES ($1, $2, $3, $4, $5, $6)
`

type CreateAssistantParams struct {
	UserID       pgtype.UUID
	Name         string
	Description  pgtype.Text
	SystemPrompt pgtype.Text
	Model        string
	Metadata     []byte
}

// CRUD for assistants
func (q *Queries) CreateAssistant(ctx context.Context, arg CreateAssistantParams) error {
	_, err := q.db.Exec(ctx, createAssistant,
		arg.UserID,
		arg.Name,
		arg.Description,
		arg.SystemPrompt,
		arg.Model,
		arg.Metadata,
	)
	return err
}

const createAssistantEmbedding = `-- name: CreateAssistantEmbedding :exec
INSERT INTO assistant_embedddings (user_id, message_id, attachment_id, embeddings)
VALUES ($1, $2, $3, $4)
`

type CreateAssistantEmbeddingParams struct {
	UserID       pgtype.UUID
	MessageID    pgtype.UUID
	AttachmentID pgtype.UUID
	Embeddings   pgvector.Vector
}

// CRUD for assistant_message_embedddings
func (q *Queries) CreateAssistantEmbedding(ctx context.Context, arg CreateAssistantEmbeddingParams) error {
	_, err := q.db.Exec(ctx, createAssistantEmbedding,
		arg.UserID,
		arg.MessageID,
		arg.AttachmentID,
		arg.Embeddings,
	)
	return err
}

const createAssistantThread = `-- name: CreateAssistantThread :exec
INSERT INTO assistant_threads (user_id, assistant_id, name, description, model, is_long_term_memory, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateAssistantThreadParams struct {
	UserID           pgtype.UUID
	AssistantID      pgtype.UUID
	Name             string
	Description      pgtype.Text
	Model            string
	IsLongTermMemory pgtype.Bool
	Metadata         []byte
}

// CRUD for assistant_threads
func (q *Queries) CreateAssistantThread(ctx context.Context, arg CreateAssistantThreadParams) error {
	_, err := q.db.Exec(ctx, createAssistantThread,
		arg.UserID,
		arg.AssistantID,
		arg.Name,
		arg.Description,
		arg.Model,
		arg.IsLongTermMemory,
		arg.Metadata,
	)
	return err
}

const createAttachment = `-- name: CreateAttachment :exec
INSERT INTO assistant_attachments (user_id, entity, entity_id, file_type, file_url, size, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

type CreateAttachmentParams struct {
	UserID   pgtype.UUID
	Entity   string
	EntityID pgtype.UUID
	FileType pgtype.Text
	FileUrl  pgtype.Text
	Size     pgtype.Int4
	Metadata []byte
}

// CRUD for assistant_attachments
func (q *Queries) CreateAttachment(ctx context.Context, arg CreateAttachmentParams) error {
	_, err := q.db.Exec(ctx, createAttachment,
		arg.UserID,
		arg.Entity,
		arg.EntityID,
		arg.FileType,
		arg.FileUrl,
		arg.Size,
		arg.Metadata,
	)
	return err
}

const createThreadMessage = `-- name: CreateThreadMessage :exec
INSERT INTO assistant_messages (user_id, thread_id, model, token, role, text, attachments, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateThreadMessageParams struct {
	UserID      pgtype.UUID
	ThreadID    pgtype.UUID
	Model       pgtype.Text
	Token       pgtype.Int4
	Role        string
	Text        pgtype.Text
	Attachments []pgtype.UUID
	Metadata    []byte
}

// CRUD for assistant_thread_messages
func (q *Queries) CreateThreadMessage(ctx context.Context, arg CreateThreadMessageParams) error {
	_, err := q.db.Exec(ctx, createThreadMessage,
		arg.UserID,
		arg.ThreadID,
		arg.Model,
		arg.Token,
		arg.Role,
		arg.Text,
		arg.Attachments,
		arg.Metadata,
	)
	return err
}

const deleteAssistant = `-- name: DeleteAssistant :exec
DELETE FROM assistants WHERE uuid = $1
`

func (q *Queries) DeleteAssistant(ctx context.Context, uuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAssistant, uuid)
	return err
}

const deleteAssistantEmbeddings = `-- name: DeleteAssistantEmbeddings :exec
DELETE FROM assistant_embedddings WHERE id = $1
`

func (q *Queries) DeleteAssistantEmbeddings(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteAssistantEmbeddings, id)
	return err
}

const deleteAssistantThread = `-- name: DeleteAssistantThread :exec
DELETE FROM assistant_threads WHERE uuid = $1
`

func (q *Queries) DeleteAssistantThread(ctx context.Context, uuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAssistantThread, uuid)
	return err
}

const deleteAttachment = `-- name: DeleteAttachment :exec
DELETE FROM assistant_attachments WHERE uuid = $1
`

func (q *Queries) DeleteAttachment(ctx context.Context, uuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteAttachment, uuid)
	return err
}

const deleteThreadMessage = `-- name: DeleteThreadMessage :exec
DELETE FROM assistant_messages WHERE uuid = $1
`

func (q *Queries) DeleteThreadMessage(ctx context.Context, uuid pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteThreadMessage, uuid)
	return err
}

const getAssistant = `-- name: GetAssistant :one
SELECT id, uuid, user_id, name, description, system_prompt, model, metadata, created_at, updated_at FROM assistants WHERE uuid = $1
`

func (q *Queries) GetAssistant(ctx context.Context, uuid pgtype.UUID) (Assistant, error) {
	row := q.db.QueryRow(ctx, getAssistant, uuid)
	var i Assistant
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
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

const getAssistantThread = `-- name: GetAssistantThread :one
SELECT id, uuid, user_id, assistant_id, name, description, model, is_long_term_memory, metadata, created_at, updated_at FROM assistant_threads WHERE uuid = $1
`

func (q *Queries) GetAssistantThread(ctx context.Context, uuid pgtype.UUID) (AssistantThread, error) {
	row := q.db.QueryRow(ctx, getAssistantThread, uuid)
	var i AssistantThread
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.AssistantID,
		&i.Name,
		&i.Description,
		&i.Model,
		&i.IsLongTermMemory,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAttachment = `-- name: GetAttachment :one
SELECT id, uuid, user_id, entity, entity_id, file_type, file_url, size, metadata, created_at, updated_at FROM assistant_attachments WHERE uuid = $1
`

func (q *Queries) GetAttachment(ctx context.Context, uuid pgtype.UUID) (AssistantAttachment, error) {
	row := q.db.QueryRow(ctx, getAttachment, uuid)
	var i AssistantAttachment
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.Entity,
		&i.EntityID,
		&i.FileType,
		&i.FileUrl,
		&i.Size,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getThreadMessage = `-- name: GetThreadMessage :one
SELECT id, uuid, user_id, thread_id, model, token, role, text, attachments, metadata, created_at, updated_at FROM assistant_messages WHERE uuid = $1
`

func (q *Queries) GetThreadMessage(ctx context.Context, uuid pgtype.UUID) (AssistantMessage, error) {
	row := q.db.QueryRow(ctx, getThreadMessage, uuid)
	var i AssistantMessage
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.UserID,
		&i.ThreadID,
		&i.Model,
		&i.Token,
		&i.Role,
		&i.Text,
		&i.Attachments,
		&i.Metadata,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAssistantThreads = `-- name: ListAssistantThreads :many
SELECT id, uuid, user_id, assistant_id, name, description, model, is_long_term_memory, metadata, created_at, updated_at FROM assistant_threads WHERE assistant_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAssistantThreads(ctx context.Context, assistantID pgtype.UUID) ([]AssistantThread, error) {
	rows, err := q.db.Query(ctx, listAssistantThreads, assistantID)
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
			&i.Model,
			&i.IsLongTermMemory,
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
SELECT id, uuid, user_id, assistant_id, name, description, model, is_long_term_memory, metadata, created_at, updated_at FROM assistant_threads WHERE user_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAssistantThreadsByUser(ctx context.Context, userID pgtype.UUID) ([]AssistantThread, error) {
	rows, err := q.db.Query(ctx, listAssistantThreadsByUser, userID)
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
			&i.Model,
			&i.IsLongTermMemory,
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

const listAssistantsByUser = `-- name: ListAssistantsByUser :many
SELECT id, uuid, user_id, name, description, system_prompt, model, metadata, created_at, updated_at FROM assistants WHERE user_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAssistantsByUser(ctx context.Context, userID pgtype.UUID) ([]Assistant, error) {
	rows, err := q.db.Query(ctx, listAssistantsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Assistant
	for rows.Next() {
		var i Assistant
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
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

const listAttachments = `-- name: ListAttachments :many
SELECT id, uuid, user_id, entity, entity_id, file_type, file_url, size, metadata, created_at, updated_at FROM assistant_attachments WHERE entity = $1 AND entity_id = $2 ORDER BY created_at DESC
`

type ListAttachmentsParams struct {
	Entity   string
	EntityID pgtype.UUID
}

func (q *Queries) ListAttachments(ctx context.Context, arg ListAttachmentsParams) ([]AssistantAttachment, error) {
	rows, err := q.db.Query(ctx, listAttachments, arg.Entity, arg.EntityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssistantAttachment
	for rows.Next() {
		var i AssistantAttachment
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
			&i.Entity,
			&i.EntityID,
			&i.FileType,
			&i.FileUrl,
			&i.Size,
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

const listAttachmentsByUser = `-- name: ListAttachmentsByUser :many
SELECT id, uuid, user_id, entity, entity_id, file_type, file_url, size, metadata, created_at, updated_at FROM assistant_attachments WHERE user_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListAttachmentsByUser(ctx context.Context, userID pgtype.UUID) ([]AssistantAttachment, error) {
	rows, err := q.db.Query(ctx, listAttachmentsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssistantAttachment
	for rows.Next() {
		var i AssistantAttachment
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
			&i.Entity,
			&i.EntityID,
			&i.FileType,
			&i.FileUrl,
			&i.Size,
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

const listThreadMessages = `-- name: ListThreadMessages :many
SELECT id, uuid, user_id, thread_id, model, token, role, text, attachments, metadata, created_at, updated_at FROM assistant_messages WHERE thread_id = $1 ORDER BY created_at DESC
`

func (q *Queries) ListThreadMessages(ctx context.Context, threadID pgtype.UUID) ([]AssistantMessage, error) {
	rows, err := q.db.Query(ctx, listThreadMessages, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AssistantMessage
	for rows.Next() {
		var i AssistantMessage
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.UserID,
			&i.ThreadID,
			&i.Model,
			&i.Token,
			&i.Role,
			&i.Text,
			&i.Attachments,
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

const similaritySearchForThreadByCosineDistance = `-- name: SimilaritySearchForThreadByCosineDistance :many
SELECT ae.id, ae.text, 1 - (embeddings <=> $2) AS score  
FROM assistant_embedddings ae
JOIN assistant_attachments aa ON ae.attachment_id = aa.uuid
WHERE ae.attachment_id IN (
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistant_messages am ON aa.entity_id = am.uuid
        JOIN assistant_threads at ON am.thread_id = at.uuid
        WHERE at.uuid = $1
    UNION
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistant_threads at ON aa.entity_id = at.uuid
        WHERE at.uuid = $1
    UNION
    SELECT aa.uuid FROM assistant_attachments aa
        JOIN assistants a ON aa.entity_id = a.uuid
        JOIN assistant_threads at ON a.uuid = at.assistant_id
        WHERE at.uuid = $1
)
AND embeddings <=> $2
ORDER BY 1 - (embedding <=> $2) LIMIT $3
`

type SimilaritySearchForThreadByCosineDistanceParams struct {
	Uuid       pgtype.UUID
	Embeddings pgvector.Vector
	Limit      int32
}

type SimilaritySearchForThreadByCosineDistanceRow struct {
	ID    int32
	Text  string
	Score int32
}

// It need combine all these results to get the final result:
// 1. assistants -> assistant_attachments -> assistant_message_embedddings
// 2. assistant_threads -> assistant_attachments -> assistant_message_embedddings
// 3. assistant_threads -> assistant_messages -> assistant_attachments -> assistant_message_embedddings
func (q *Queries) SimilaritySearchForThreadByCosineDistance(ctx context.Context, arg SimilaritySearchForThreadByCosineDistanceParams) ([]SimilaritySearchForThreadByCosineDistanceRow, error) {
	rows, err := q.db.Query(ctx, similaritySearchForThreadByCosineDistance, arg.Uuid, arg.Embeddings, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SimilaritySearchForThreadByCosineDistanceRow
	for rows.Next() {
		var i SimilaritySearchForThreadByCosineDistanceRow
		if err := rows.Scan(&i.ID, &i.Text, &i.Score); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAssistant = `-- name: UpdateAssistant :exec
UPDATE assistants SET name = $2, description = $3, system_prompt = $4, model = $5, metadata = $6
WHERE uuid = $1
`

type UpdateAssistantParams struct {
	Uuid         pgtype.UUID
	Name         string
	Description  pgtype.Text
	SystemPrompt pgtype.Text
	Model        string
	Metadata     []byte
}

func (q *Queries) UpdateAssistant(ctx context.Context, arg UpdateAssistantParams) error {
	_, err := q.db.Exec(ctx, updateAssistant,
		arg.Uuid,
		arg.Name,
		arg.Description,
		arg.SystemPrompt,
		arg.Model,
		arg.Metadata,
	)
	return err
}

const updateAssistantThread = `-- name: UpdateAssistantThread :exec
UPDATE assistant_threads SET name = $2, description = $3, model = $4, is_long_term_memory = $5, metadata = $6 
WHERE uuid = $1
`

type UpdateAssistantThreadParams struct {
	Uuid             pgtype.UUID
	Name             string
	Description      pgtype.Text
	Model            string
	IsLongTermMemory pgtype.Bool
	Metadata         []byte
}

func (q *Queries) UpdateAssistantThread(ctx context.Context, arg UpdateAssistantThreadParams) error {
	_, err := q.db.Exec(ctx, updateAssistantThread,
		arg.Uuid,
		arg.Name,
		arg.Description,
		arg.Model,
		arg.IsLongTermMemory,
		arg.Metadata,
	)
	return err
}

const updateAttachment = `-- name: UpdateAttachment :exec
UPDATE assistant_attachments SET file_type = $2, file_url = $3, size = $4, metadata = $5 WHERE uuid = $1
`

type UpdateAttachmentParams struct {
	Uuid     pgtype.UUID
	FileType pgtype.Text
	FileUrl  pgtype.Text
	Size     pgtype.Int4
	Metadata []byte
}

func (q *Queries) UpdateAttachment(ctx context.Context, arg UpdateAttachmentParams) error {
	_, err := q.db.Exec(ctx, updateAttachment,
		arg.Uuid,
		arg.FileType,
		arg.FileUrl,
		arg.Size,
		arg.Metadata,
	)
	return err
}

const updateThreadMessage = `-- name: UpdateThreadMessage :exec
UPDATE assistant_messages SET text = $2, attachments = $3, metadata = $4 WHERE uuid = $1
`

type UpdateThreadMessageParams struct {
	Uuid        pgtype.UUID
	Text        pgtype.Text
	Attachments []pgtype.UUID
	Metadata    []byte
}

func (q *Queries) UpdateThreadMessage(ctx context.Context, arg UpdateThreadMessageParams) error {
	_, err := q.db.Exec(ctx, updateThreadMessage,
		arg.Uuid,
		arg.Text,
		arg.Attachments,
		arg.Metadata,
	)
	return err
}
