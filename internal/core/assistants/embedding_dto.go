package assistants

import (
	"time"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type EmbeddingMetadata struct{}

type EmbeddingDTO struct {
	ID           int64     `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	AttachmentID uuid.UUID `json:"attachment_id"`
	Text         string    `json:"text"`
	Embeddings   []float32 `json:"embeddings"`
	Metadata     []byte    `json:"metadata"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (e *EmbeddingDTO) Load(dbo *db.AssistantEmbeddding) {
	e.UserID = dbo.UserID.Bytes
	e.AttachmentID = dbo.AttachmentID.Bytes
	e.Text = dbo.Text
	e.Embeddings = dbo.Embeddings.Slice()
	e.Metadata = dbo.Metadata
	e.CreatedAt = dbo.CreatedAt.Time
	e.UpdatedAt = dbo.UpdatedAt.Time
}

func (e *EmbeddingDTO) Dump() *db.AssistantEmbeddding {
	vec := pgvector.NewVector(e.Embeddings)
	return &db.AssistantEmbeddding{
		UserID:       pgtype.UUID{Bytes: e.UserID, Valid: e.UserID != uuid.Nil},
		AttachmentID: pgtype.UUID{Bytes: e.AttachmentID, Valid: e.AttachmentID != uuid.Nil},
		Text:         e.Text,
		Embeddings:   &vec,
		Metadata:     e.Metadata,
	}
}
