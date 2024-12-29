package assistants

import (
	"encoding/json"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ThreadMetadata struct {
	AssistantMetadata
	IsGeneratedTitle bool `json:"is_generated_title"`
}

// Merge merges the thread metadata with the assistant metadata
func (m *ThreadMetadata) Merge(am AssistantMetadata) {
	mergedRagSettings := am.RagSettings // Start with assistant's RagSettings

	// Override assistant's RagSettings if thread has non-default values
	mergedRagSettings.Enable = m.RagSettings.Enable
	mergedRagSettings.MultiQuery = m.RagSettings.MultiQuery
	mergedRagSettings.QueryRewrite = m.RagSettings.QueryRewrite
	mergedRagSettings.Rerank = m.RagSettings.Rerank

	mergedTools := am.Tools // Take assistant tools first
	if len(m.Tools) > 0 {
		mergedTools = m.Tools // Override if thread has its own tools
	}

	m.Tools = mergedTools
	m.RagSettings = mergedRagSettings
}

type ThreadDTO struct {
	Id           uuid.UUID      `json:"id"`
	UserId       uuid.UUID      `json:"user_id"`
	AssistantId  uuid.UUID      `json:"assistant_id"`
	SystemPrompt string         `json:"system_prompt"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Model        string         `json:"model"`
	Metadata     ThreadMetadata `json:"metadata"`
	Messages     []MessageDTO   `json:"messages"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
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
	t.Messages = append(t.Messages, MessageDTO{
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

func AddThreadMessage(message MessageDTO) ThreadOption {
	return func(t *ThreadDTO) {
		t.Messages = append(t.Messages, message)
	}
}

func WithThreadMessages(messages []MessageDTO) ThreadOption {
	return func(t *ThreadDTO) {
		t.Messages = messages
	}
}
