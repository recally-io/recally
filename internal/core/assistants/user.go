package assistants

import "github.com/google/uuid"

type User struct {
	ID                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Telegram            string    `json:"telegram"`
	Github              string    `json:"github"`
	Google              string    `json:"google"`
	ActivateAssistantID uuid.UUID
	ActivateThreadID    uuid.UUID
}
