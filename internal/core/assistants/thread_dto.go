package assistants

import (
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ThreadMetadata struct {
	IsGeneratedTitle bool     `json:"is_generated_title"`
	Tools            []string `json:"tools"`
}

type ThreadDTO struct {
	Id           uuid.UUID          `json:"id"`
	UserId       uuid.UUID          `json:"user_id"`
	AssistantId  uuid.UUID          `json:"assistant_id"`
	SystemPrompt string             `json:"system_prompt"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Model        string             `json:"model"`
	Metadata     ThreadMetadata     `json:"metadata"`
	Messages     []ThreadMessageDTO `json:"messages"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

func (t *ThreadDTO) Load(dbo *db.AssistantThread) {
	t.Id = dbo.Uuid
	t.UserId = dbo.UserID.Bytes
	t.AssistantId = dbo.AssistantID.Bytes
	t.SystemPrompt = dbo.SystemPrompt.String
	t.Name = dbo.Name
	t.Description = dbo.Description.String
	t.Model = dbo.Model
	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &t.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Thread metadata", "err", err, "metadata", string(dbo.Metadata))
		}
	}

	t.CreatedAt = dbo.CreatedAt.Time
	t.UpdatedAt = dbo.UpdatedAt.Time
}

func (d *ThreadDTO) Dump() *db.AssistantThread {
	metadata, _ := json.Marshal(d.Metadata)
	if d.Id == uuid.Nil {
		d.Id = uuid.New()
	}
	return &db.AssistantThread{
		Uuid:         d.Id,
		UserID:       pgtype.UUID{Bytes: d.UserId, Valid: true},
		AssistantID:  pgtype.UUID{Bytes: d.AssistantId, Valid: true},
		SystemPrompt: pgtype.Text{String: d.SystemPrompt, Valid: d.SystemPrompt != ""},
		Name:         d.Name,
		Description:  pgtype.Text{String: d.Description, Valid: d.Description != ""},
		Model:        d.Model,
		Metadata:     metadata,
	}
}

func (t *ThreadDTO) AddMessage(role, text string) {
	t.Messages = append(t.Messages, ThreadMessageDTO{
		Role: role,
		Text: text,
	})
}

type ThreadOption func(*ThreadDTO)

func NewThread(userId uuid.UUID, assistant AssistantDTO, opts ...ThreadOption) *ThreadDTO {
	t := &ThreadDTO{
		UserId:       userId,
		AssistantId:  assistant.Id,
		SystemPrompt: assistant.SystemPrompt,
		Model:        assistant.Model,
		Name:         "Thread" + uuid.New().String(),
		Description:  "I am a conversation thread.",
	}

	for _, opt := range opts {
		opt(t)
	}
	return t
}

func WithThreadName(name string) ThreadOption {
	return func(t *ThreadDTO) {
		t.Name = name
	}
}

func WithThreadDescription(description string) ThreadOption {
	return func(t *ThreadDTO) {
		t.Description = description
	}
}

func WithThreadModel(model string) ThreadOption {
	return func(t *ThreadDTO) {
		t.Model = model
	}
}

func WithThreadMetaData(metaData ThreadMetadata) ThreadOption {
	return func(t *ThreadDTO) {
		t.Metadata = metaData
	}
}

func AddThreadMessage(message ThreadMessageDTO) ThreadOption {
	return func(t *ThreadDTO) {
		t.Messages = append(t.Messages, message)
	}
}

func WithThreadMessages(messages []ThreadMessageDTO) ThreadOption {
	return func(t *ThreadDTO) {
		t.Messages = messages
	}
}
