package assistants

import (
	"time"

	"github.com/google/uuid"
)

type ThreadMetaData struct{}

type ThreadMessage struct {
	UserID    uuid.UUID
	ThreadID  uuid.UUID
	Model     string
	Token     int
	Role      string
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Thread struct {
	Id           uuid.UUID       `json:"uuid"`
	UserId       uuid.UUID       `json:"user_id"`
	AssistantId  uuid.UUID       `json:"assistant_id"`
	SystemPrompt string          `json:"system_prompt"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Model        string          `json:"model"`
	MetaData     ThreadMetaData  `json:"metadata"`
	Messages     []ThreadMessage `json:"messages"`
}

func (t *Thread) AddMessage(role, text string) {
	t.Messages = append(t.Messages, ThreadMessage{
		Role: role,
		Text: text,
	})
}

type ThreadOption func(*Thread)

func NewThread(userId uuid.UUID, assistant Assistant, opts ...ThreadOption) *Thread {
	t := &Thread{
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
	return func(t *Thread) {
		t.Name = name
	}
}

func WithThreadDescription(description string) ThreadOption {
	return func(t *Thread) {
		t.Description = description
	}
}

func WithThreadModel(model string) ThreadOption {
	return func(t *Thread) {
		t.Model = model
	}
}

func WithThreadMetaData(metaData ThreadMetaData) ThreadOption {
	return func(t *Thread) {
		t.MetaData = metaData
	}
}

func AddThreadMessage(message ThreadMessage) ThreadOption {
	return func(t *Thread) {
		t.Messages = append(t.Messages, message)
	}
}

func WithThreadMessages(messages []ThreadMessage) ThreadOption {
	return func(t *Thread) {
		t.Messages = messages
	}
}
