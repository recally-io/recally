// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pgv "github.com/pgvector/pgvector-go"
)

type Assistant struct {
	ID           int32
	Uuid         uuid.UUID
	UserID       pgtype.UUID
	Name         string
	Description  pgtype.Text
	SystemPrompt pgtype.Text
	Model        string
	Metadata     []byte
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type AssistantAttachment struct {
	ID          int32
	Uuid        uuid.UUID
	UserID      pgtype.UUID
	AssistantID pgtype.UUID
	ThreadID    pgtype.UUID
	Name        pgtype.Text
	Type        pgtype.Text
	Url         pgtype.Text
	Size        pgtype.Int4
	Metadata    []byte
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

type AssistantEmbeddding struct {
	ID           int32
	UserID       pgtype.UUID
	AttachmentID pgtype.UUID
	Text         string
	Embeddings   *pgv.Vector
	Metadata     []byte
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
	Uuid         uuid.UUID
}

type AssistantMessage struct {
	ID              int32
	Uuid            uuid.UUID
	UserID          pgtype.UUID
	AssistantID     pgtype.UUID
	ThreadID        pgtype.UUID
	Model           pgtype.Text
	Role            string
	Text            pgtype.Text
	PromptToken     pgtype.Int4
	CompletionToken pgtype.Int4
	Embeddings      *pgv.Vector
	Metadata        []byte
	CreatedAt       pgtype.Timestamptz
	UpdatedAt       pgtype.Timestamptz
}

type AssistantThread struct {
	ID           int32
	Uuid         uuid.UUID
	UserID       pgtype.UUID
	AssistantID  pgtype.UUID
	Name         string
	Description  pgtype.Text
	SystemPrompt pgtype.Text
	Model        string
	Metadata     []byte
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

type Bookmark struct {
	ID                int32
	Uuid              uuid.UUID
	UserID            pgtype.UUID
	Url               string
	Title             pgtype.Text
	Summary           pgtype.Text
	SummaryEmbeddings *pgv.Vector
	Content           pgtype.Text
	ContentEmbeddings *pgv.Vector
	Html              pgtype.Text
	Metadata          []byte
	Screenshot        pgtype.Text
	CreatedAt         pgtype.Timestamptz
	UpdatedAt         pgtype.Timestamptz
}

type Cache struct {
	ID        int32
	Domain    string
	Key       string
	Value     []byte
	ExpiresAt pgtype.Timestamp
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type User struct {
	ID                  int32
	Uuid                uuid.UUID
	Username            pgtype.Text
	PasswordHash        pgtype.Text
	Email               pgtype.Text
	Github              pgtype.Text
	Google              pgtype.Text
	Telegram            pgtype.Text
	ActivateAssistantID pgtype.UUID
	ActivateThreadID    pgtype.UUID
	Status              string
	CreatedAt           pgtype.Timestamptz
	UpdatedAt           pgtype.Timestamptz
}
