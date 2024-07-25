package handlers

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                  uuid.UUID `json:"id"`
	Username            string    `json:"username"`
	Telegram            string    `json:"telegram"`
	ActivateAssistantID uuid.UUID
	ActivateThreadID    uuid.UUID
}

func (h *Handler) getOrCreateUser(ctx context.Context) (*User, error) {
	userID := ctx.Value(constant.ContextKey(constant.ContextKeyUserID)).(string)
	userName := ctx.Value(constant.ContextKey(constant.ContextKeyUserName)).(string)
	user, err := h.tx.GetTelegramUser(ctx, pgtype.Text{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		params := db.InserUserParams{
			Username: pgtype.Text{String: userName, Valid: true},
			Telegram: pgtype.Text{String: userID, Valid: true},
		}
		var err error
		user, err = h.tx.InserUser(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to insert user: %w", err)
		}
	}

	return &User{
		ID:                  user.Uuid,
		Username:            user.Username.String,
		Telegram:            user.Telegram.String,
		ActivateAssistantID: user.ActivateAssistantID.Bytes,
		ActivateThreadID:    user.ActivateThreadID.Bytes,
	}, nil
}

func (h *Handler) updateUser(ctx context.Context, user User, assistantId uuid.UUID, threadId uuid.UUID) (*User, error) {
	dbUser, err := h.tx.UpdateTelegramUser(ctx, db.UpdateTelegramUserParams{
		Telegram:            pgtype.Text{String: user.Telegram, Valid: true},
		ActivateAssistantID: pgtype.UUID{Bytes: assistantId, Valid: assistantId != uuid.Nil},
		ActivateThreadID:    pgtype.UUID{Bytes: threadId, Valid: threadId != uuid.Nil},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &User{
		ID:                  dbUser.Uuid,
		Username:            dbUser.Username.String,
		Telegram:            dbUser.Telegram.String,
		ActivateAssistantID: dbUser.ActivateAssistantID.Bytes,
		ActivateThreadID:    dbUser.ActivateThreadID.Bytes,
	}, nil
}
