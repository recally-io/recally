package auth

import (
	"time"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserDTO struct {
	ID                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	Email               string    `json:"email"`
	Telegram            string    `json:"telegram"`
	Github              string    `json:"github"`
	Google              string    `json:"google"`
	ActivateAssistantID uuid.UUID
	ActivateThreadID    uuid.UUID
	Status              string
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (t *UserDTO) Load(d *db.User) {
	t.ID = d.Uuid
	t.Username = d.Username.String
	t.Email = d.Email.String
	t.Telegram = d.Telegram.String
	t.Github = d.Github.String
	t.Google = d.Google.String
	t.ActivateAssistantID = d.ActivateAssistantID.Bytes
	t.ActivateThreadID = d.ActivateThreadID.Bytes
	t.Status = d.Status
	t.CreatedAt = d.CreatedAt.Time
	t.UpdatedAt = d.UpdatedAt.Time
}

func (t *UserDTO) Dump() *db.User {
	return &db.User{
		Uuid:                t.ID,
		Username:            pgtype.Text{String: t.Username, Valid: t.Username != ""},
		Email:               pgtype.Text{String: t.Email, Valid: t.Email != ""},
		Telegram:            pgtype.Text{String: t.Telegram, Valid: t.Telegram != ""},
		Github:              pgtype.Text{String: t.Github, Valid: t.Github != ""},
		Google:              pgtype.Text{String: t.Google, Valid: t.Google != ""},
		ActivateAssistantID: pgtype.UUID{Bytes: t.ActivateAssistantID, Valid: t.ActivateAssistantID != uuid.Nil},
		ActivateThreadID:    pgtype.UUID{Bytes: t.ActivateThreadID, Valid: t.ActivateThreadID != uuid.Nil},
		Status:              t.Status,
	}
}
