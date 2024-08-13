package assistants

import (
	"time"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
)

type ThreadMetaData struct{}

type ThreadMessage struct {
	UserID    uuid.UUID `json:"user_id"`
	ThreadID  uuid.UUID `json:"thread_id"`
	Model     string    `json:"model"`
	Token     int       `json:"token"`
	Role      string    `json:"role"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *ThreadMessage) FromDBO(dbo *db.AssistantMessage) {
	t.UserID = dbo.UserID.Bytes
	t.ThreadID = dbo.ThreadID.Bytes
	t.Model = dbo.Model.String
	t.Token = int(dbo.Token.Int32)
	t.Role = dbo.Role
	t.Text = dbo.Text.String
	t.CreatedAt = dbo.CreatedAt.Time
	t.UpdatedAt = dbo.UpdatedAt.Time
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
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (t *Thread) FromDBO(dbo *db.AssistantThread) {
	t.Id = dbo.Uuid
	t.UserId = dbo.UserID.Bytes
	t.AssistantId = dbo.AssistantID.Bytes
	t.SystemPrompt = dbo.SystemPrompt.String
	t.Name = dbo.Name
	t.Description = dbo.Description.String
	t.Model = dbo.Model
	t.CreatedAt = dbo.CreatedAt.Time
	t.UpdatedAt = dbo.UpdatedAt.Time
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
