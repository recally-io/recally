package handlers

import (
	"context"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"
	"strings"

	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) LinkAccountHandler(c tele.Context) error {
	ctx := c.Get(contexts.ContextKeyContext).(context.Context)

	tx, ok := contexts.Get[pgx.Tx](ctx, contexts.ContextKeyTx)
	if !ok {
		_ = c.Reply("Failed to get dbtx from context")
		return nil
	}

	userID, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	if !ok {
		_ = c.Reply("Failed to get userID from context")
		return nil
	}

	token := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/linkaccount"))
	if token == "" {
		_ = c.Reply("Invalid token")
		return nil
	}

	oAuthUser := auth.OAuth2User{
		Provider: "telegram",
		ID:       userID,
		Name:     c.Sender().Username,
	}

	if err := h.authService.LinkAccount(ctx, tx, oAuthUser, token); err != nil {
		_ = c.Reply("Failed to link user: " + err.Error())
		return nil
	}

	return c.Reply("User linked successfully")
}
