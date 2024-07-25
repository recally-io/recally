package assistants

import (
	"github.com/google/uuid"
)

type ThreadMetaData struct{}

type Thread struct {
	Id          uuid.UUID      `json:"uuid"`
	UserId      uuid.UUID      `json:"user_id"`
	AssistantId uuid.UUID      `json:"assistant_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Model       string         `json:"model"`
	MetaData    ThreadMetaData `json:"metadata"`
}

type ThreadOption func(*Thread)

func NewThread(userId uuid.UUID, assistant Assistant, opts ...ThreadOption) *Thread {
	t := &Thread{
		UserId:      userId,
		AssistantId: assistant.Id,
		Model:       assistant.Model,
		Name:        "Thread" + uuid.New().String(),
		Description: "I am a conversation thread.",
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

