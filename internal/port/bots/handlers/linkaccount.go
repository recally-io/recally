package handlers

import (
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/logger"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) LinkAccountHandler(c tele.Context) error {
	ctx, user, tx, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	token := strings.TrimSpace(strings.TrimPrefix(c.Text(), "/linkaccount"))
	if token == "" {
		_ = c.Reply("Invalid token")
		return nil
	}

	tgUserID, _ := contexts.Get[string](ctx, contexts.ContextKeyTelegramID)
	oAuthUser := auth.OAuth2User{
		Provider: "telegram",
		ID:       tgUserID,
		Name:     user.Username,
	}

	if err := h.authService.LinkAccount(ctx, tx, oAuthUser, token); err != nil {
		_ = c.Reply("Failed to link user: " + err.Error())
		return nil
	}

	return c.Reply("User linked successfully")
}
