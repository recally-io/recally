package handlers

import (
	"context"
	"fmt"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"
	"strings"

	"github.com/jackc/pgx/v5"
	"gopkg.in/telebot.v3"
)

func (h *Handler) initHandlerRequest(c telebot.Context) (context.Context, *auth.UserDTO, pgx.Tx, error) {
	ctx := c.Get(contexts.ContextKeyContext).(context.Context)

	tx, ok := contexts.Get[pgx.Tx](ctx, contexts.ContextKeyTx)
	if !ok {
		return ctx, nil, nil, fmt.Errorf("failed to get dbtx from context")
	}

	telegramID, ok := contexts.Get[string](ctx, contexts.ContextKeyTelegramID)
	if !ok {
		return ctx, nil, tx, fmt.Errorf("failed to get telegram ID from context")
	}
	user, err := h.authService.GetTelegramUser(ctx, tx, telegramID)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrNotFound) {
			telegramName, _ := contexts.Get[string](ctx, contexts.ContextKeyTelegramName)
			user, err = h.authService.CreateTelegramUser(ctx, tx, telegramName, telegramID)
			if err != nil {
				return ctx, nil, tx, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, nil, tx, fmt.Errorf("failed to get user: %w", err)
		}
	}
	ctx = auth.SetUserToContext(ctx, user)
	return ctx, user, tx, nil
}
