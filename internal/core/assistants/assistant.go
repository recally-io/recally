package assistants

import (
	"vibrain/internal/pkg/llms"

	"github.com/google/uuid"
)

type AssistantMetaData struct{}

type Assistant struct {
	Id           uuid.UUID         `json:"uuid"`
	UserId       uuid.UUID         `json:"user_id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	SystemPrompt string            `json:"system_prompt"`
	Model        string            `json:"model"`
	MetaData     AssistantMetaData `json:"metadata"`
}

type AssistantOption func(*Assistant)

func NewAssistant(userId uuid.UUID, opts ...AssistantOption) *Assistant {
	a := &Assistant{
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
	return func(a *Assistant) {
		a.Name = name
	}
}

func WithAssistantDescription(description string) AssistantOption {
	return func(a *Assistant) {
		a.Description = description
	}
}

func WithAssistantSystemPrompt(systemPrompt string) AssistantOption {
	return func(a *Assistant) {
		a.SystemPrompt = systemPrompt
	}
}

func WithAssistantModel(model string) AssistantOption {
	return func(a *Assistant) {
		a.Model = model
	}
}

func WithAssistantMetaData(metaData AssistantMetaData) AssistantOption {
	return func(a *Assistant) {
		a.MetaData = metaData
	}
}