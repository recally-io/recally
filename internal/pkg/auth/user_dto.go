package auth

import (
	"encoding/json"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SummaryConfig struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	Language string `json:"language"`
}

type UserSettings struct {
	SummaryOptions      SummaryConfig `json:"summary_options"`
	IsLinkedTelegramBot bool          `json:"is_linked_telegram_bot"`
}

type UserDTO struct {
	ID                  uuid.UUID    `json:"id"`
	Username            string       `json:"username"`
	Email               string       `json:"email"`
	Phone               string       `json:"phone"`
	Password            string       `json:"password"`
	ActivateAssistantID uuid.UUID    `json:"activate_assistant_id"`
	ActivateThreadID    uuid.UUID    `json:"activate_thread_id"`
	Status              string       `json:"status"`
	Settings            UserSettings `json:"settings"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
}

func (t *UserDTO) Load(d *db.User) {
	t.ID = d.Uuid
	t.Username = d.Username.String
	t.Email = d.Email.String
	t.Phone = d.Phone.String
	t.Password = d.PasswordHash.String
	t.ActivateAssistantID = d.ActivateAssistantID.Bytes
	t.ActivateThreadID = d.ActivateThreadID.Bytes
	t.Status = d.Status
	if d.Settings != nil {
		if err := json.Unmarshal(d.Settings, &t.Settings); err != nil {
			logger.Default.Warn("failed to unmarshal user settings", "err", err, "settings", string(d.Settings))
		}
	}
	t.CreatedAt = d.CreatedAt.Time
	t.UpdatedAt = d.UpdatedAt.Time
}

func (t *UserDTO) Dump() *db.User {
	settings, _ := json.Marshal(t.Settings)
	return &db.User{
		Uuid:                t.ID,
		Username:            pgtype.Text{String: t.Username, Valid: t.Username != ""},
		Email:               pgtype.Text{String: t.Email, Valid: t.Email != ""},
		Phone:               pgtype.Text{String: t.Phone, Valid: t.Phone != ""},
		PasswordHash:        pgtype.Text{String: t.Password, Valid: t.Password != ""},
		ActivateAssistantID: pgtype.UUID{Bytes: t.ActivateAssistantID, Valid: t.ActivateAssistantID != uuid.Nil},
		ActivateThreadID:    pgtype.UUID{Bytes: t.ActivateThreadID, Valid: t.ActivateThreadID != uuid.Nil},
		Status:              t.Status,
		Settings:            settings,
	}
}
