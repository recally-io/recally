package handlers

import "github.com/google/uuid"

type User struct {
	ID                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	Telegram            string    `json:"telegram"`
	ActivateAssistantID uuid.UUID
	ActivateThreadID    uuid.UUID
}
