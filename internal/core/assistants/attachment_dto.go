package assistants

import (
	"encoding/json"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AttachmentMetadata struct{}

type AttachmentDTO struct {
	Id          uuid.UUID          `json:"id"`
	UserId      uuid.UUID          `json:"user_id"`
	AssistantId uuid.UUID          `json:"assistant_id"`
	ThreadId    uuid.UUID          `json:"thread_id"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	URL         string             `json:"url"`
	Size        int                `json:"size"`
	Metadata    AttachmentMetadata `json:"metadata"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (a *AttachmentDTO) Load(dbo *db.AssistantAttachment) {
	a.Id = dbo.Uuid
	a.UserId = dbo.UserID.Bytes
	a.AssistantId = dbo.AssistantID.Bytes
	a.ThreadId = dbo.ThreadID.Bytes
	a.Name = dbo.Name.String
	a.Type = dbo.Type.String
	a.URL = dbo.Url.String
	a.Size = int(dbo.Size.Int32)

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &a.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Attachment metadata", "err", err, "metadata", string(dbo.Metadata))
		}
	}

	a.CreatedAt = dbo.CreatedAt.Time
	a.UpdatedAt = dbo.UpdatedAt.Time
}

func (a *AttachmentDTO) Dump() *db.AssistantAttachment {
	metadata, _ := json.Marshal(a.Metadata)

	if a.Id == uuid.Nil {
		a.Id = uuid.New()
	}

	return &db.AssistantAttachment{
		Uuid:        a.Id,
		UserID:      pgtype.UUID{Bytes: a.UserId, Valid: a.UserId != uuid.Nil},
		AssistantID: pgtype.UUID{Bytes: a.AssistantId, Valid: a.AssistantId != uuid.Nil},
		ThreadID:    pgtype.UUID{Bytes: a.ThreadId, Valid: a.ThreadId != uuid.Nil},
		Name:        pgtype.Text{String: a.Name, Valid: a.Name != ""},
		Type:        pgtype.Text{String: a.Type, Valid: a.Type != ""},
		Url:         pgtype.Text{String: a.URL, Valid: a.URL != ""},
		Size:        pgtype.Int4{Int32: int32(a.Size), Valid: a.Size > 0},
		Metadata:    metadata,
	}
}
