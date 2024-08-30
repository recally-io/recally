package assistants

import (
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type MessageMetadata struct {
	Tools  []string `json:"tools"`
	Images []string `json:"images"`
}

type MessageDTO struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	AssistantID     uuid.UUID       `json:"assistant_id"`
	ThreadID        uuid.UUID       `json:"thread_id"`
	Model           string          `json:"model"`
	Role            string          `json:"role"`
	Text            string          `json:"text"`
	PromptToken     int32           `json:"prompt_token"`
	CompletionToken int32           `json:"completion_token"`
	Embeddings      []float32       `json:"embeddings"`
	Metadata        MessageMetadata `json:"metadata"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

func (d *MessageDTO) Load(dbo *db.AssistantMessage) {
	d.ID = dbo.Uuid
	d.UserID = dbo.UserID.Bytes
	d.AssistantID = dbo.AssistantID.Bytes
	d.ThreadID = dbo.ThreadID.Bytes
	d.Model = dbo.Model.String
	d.Role = dbo.Role
	d.Text = dbo.Text.String
	d.PromptToken = dbo.PromptToken.Int32
	d.CompletionToken = dbo.CompletionToken.Int32
	d.Embeddings = dbo.Embeddings.Slice()
	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &d.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal ThreadMessage metadata", "err", err, "metadata", string(dbo.Metadata))
		}
	}
	d.CreatedAt = dbo.CreatedAt.Time
	d.UpdatedAt = dbo.UpdatedAt.Time
}

func (d *MessageDTO) Dump() *db.AssistantMessage {
	metadata, _ := json.Marshal(d.Metadata)
	if d.Embeddings == nil {
		d.Embeddings = []float32{}
	}
	return &db.AssistantMessage{
		UserID:          pgtype.UUID{Bytes: d.UserID, Valid: true},
		AssistantID:     pgtype.UUID{Bytes: d.AssistantID, Valid: true},
		ThreadID:        pgtype.UUID{Bytes: d.ThreadID, Valid: true},
		Model:           pgtype.Text{String: d.Model, Valid: d.Model != ""},
		Role:            d.Role,
		Text:            pgtype.Text{String: d.Text, Valid: d.Text != ""},
		PromptToken:     pgtype.Int4{Int32: d.PromptToken, Valid: d.PromptToken != 0},
		CompletionToken: pgtype.Int4{Int32: d.CompletionToken, Valid: d.CompletionToken != 0},
		Embeddings:      pgvector.NewVector(d.Embeddings),
		Metadata:        metadata,
	}
}