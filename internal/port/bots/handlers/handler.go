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

func (h *Handler) initHandlerRequest(c telebot.Context) (context.Context, *auth.UserDTO, db.DBTX, error) {
	ctx := c.Get(contexts.ContextKeyContext).(context.Context)

	tx, ok := contexts.Get[pgx.Tx](ctx, contexts.ContextKeyTx)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to get dbtx from context")
	}

	userID, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to get userID from context")
	}
	user, err := h.authService.GetTelegramUser(ctx, tx, userID)
	if err != nil {
		if strings.Contains(err.Error(), db.ErrNotFound) {
			userName := ctx.Value(contexts.ContextKey(contexts.ContextKeyUserName)).(string)
			user, err = h.authService.CreateTelegramUser(ctx, tx, userName, userID)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return nil, nil, nil, fmt.Errorf("failed to get user: %w", err)
		}
	}

	return ctx, user, tx, nil
}
