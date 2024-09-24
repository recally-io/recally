package assistants

import (
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RagSettings struct {
	Enable       bool `json:"enable"`
	MultiQuery   bool `json:"multi_query"`
	QueryRewrite bool `json:"query_rewrite"`
	Rerank       bool `json:"rerank"`
}

type AssistantMetadata struct {
	// Tools is a list of tools that the assistant can use
	Tools       []string    `json:"tools,omitempty"`
	RagSettings RagSettings `json:"rag_settings,omitempty"`
}

type AssistantDTO struct {
	Id           uuid.UUID         `json:"id"`
	UserId       uuid.UUID         `json:"user_id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	SystemPrompt string            `json:"system_prompt"`
	Model        string            `json:"model"`
	Metadata     AssistantMetadata `json:"metadata"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// Load converts a database object to a domain object
func (a *AssistantDTO) Load(dbo *db.Assistant) {
	a.Id = dbo.Uuid
	a.UserId = dbo.UserID.Bytes
	a.Name = dbo.Name
	a.Description = dbo.Description.String
	a.SystemPrompt = dbo.SystemPrompt.String
	a.Model = dbo.Model
	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &a.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Assistant metadata", "err", err, "metadata", string(dbo.Metadata))
		}
	}
	a.CreatedAt = dbo.CreatedAt.Time
	a.UpdatedAt = dbo.UpdatedAt.Time
}

// Dump converts a domain object to a database object
func (a *AssistantDTO) Dump() *db.Assistant {
	metadata, _ := json.Marshal(a.Metadata)
	return &db.Assistant{
		UserID:       pgtype.UUID{Bytes: a.UserId, Valid: a.UserId != uuid.Nil},
		Name:         a.Name,
		Description:  pgtype.Text{String: a.Description, Valid: a.Description != ""},
		SystemPrompt: pgtype.Text{String: a.SystemPrompt, Valid: a.SystemPrompt != ""},
		Model:        a.Model,
		Metadata:     metadata,
	}
}

type AssistantOption func(*AssistantDTO)

func NewAssistant(userId uuid.UUID, opts ...AssistantOption) *AssistantDTO {
	a := &AssistantDTO{
		UserId:       userId,
		Model:        llms.OpenAIGPT4oMini,
		Name:         "Assistant" + uuid.New().String(),
		Description:  "I am an AI assistant. I can help you with a variety of tasks.",
		SystemPrompt: "You are a helpful AI assistant.",
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithAssistantName(name string) AssistantOption {
	return func(a *AssistantDTO) {
		a.Name = name
	}
}

func WithAssistantDescription(description string) AssistantOption {
	return func(a *AssistantDTO) {
		a.Description = description
	}
}

func WithAssistantSystemPrompt(systemPrompt string) AssistantOption {
	return func(a *AssistantDTO) {
		a.SystemPrompt = systemPrompt
	}
}

func WithAssistantModel(model string) AssistantOption {
	return func(a *AssistantDTO) {
		a.Model = model
	}
}

func WithAssistantMetaData(metaData AssistantMetadata) AssistantOption {
	return func(a *AssistantDTO) {
		a.Metadata = metaData
	}
}
