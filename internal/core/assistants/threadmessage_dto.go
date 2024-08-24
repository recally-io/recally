package assistants

import (
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ThreadMessageMetadata struct {
	Tools []string `json:"tools"`
}

type ThreadMessageDTO struct {
	ID          uuid.UUID             `json:"id"`
	UserID      uuid.UUID             `json:"user_id"`
	ThreadID    uuid.UUID             `json:"thread_id"`
	Model       string                `json:"model"`
	Token       int                   `json:"token"`
	Role        string                `json:"role"`
	Text        string                `json:"text"`
	Attachments []uuid.UUID           `json:"attachments"`
	Metadata    ThreadMessageMetadata `json:"metadata"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

func (d *ThreadMessageDTO) Load(dbo *db.AssistantMessage) {
	d.ID = dbo.Uuid
	d.UserID = dbo.UserID.Bytes
	d.ThreadID = dbo.ThreadID.Bytes
	d.Model = dbo.Model.String
	d.Token = int(dbo.Token.Int32)
	d.Role = dbo.Role
	d.Text = dbo.Text.String
	d.Attachments = dbo.Attachments
	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &d.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal ThreadMessage metadata", "err", err, "metadata", string(dbo.Metadata))
		}
	}
	d.CreatedAt = dbo.CreatedAt.Time
	d.UpdatedAt = dbo.UpdatedAt.Time
}

func (d *ThreadMessageDTO) Dump() *db.AssistantMessage {
	metadata, _ := json.Marshal(d.Metadata)
	return &db.AssistantMessage{
		UserID:      pgtype.UUID{Bytes: d.UserID, Valid: true},
		ThreadID:    pgtype.UUID{Bytes: d.ThreadID, Valid: true},
		Model:       pgtype.Text{String: d.Model, Valid: d.Model != ""},
		Token:       pgtype.Int4{Int32: int32(d.Token), Valid: d.Token != 0},
		Role:        d.Role,
		Text:        pgtype.Text{String: d.Text, Valid: d.Text != ""},
		Attachments: d.Attachments,
		Metadata:    metadata,
	}
}
