package handlers

import (
	"context"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	"gopkg.in/telebot.v3"
)

func (h *Handler) LLMChatHandler(c telebot.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	return nil
}

func (h *Handler) LLMMemChatHandler(c telebot.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	return nil
}

func (h *Handler) LLMChatClearContextHandler(c telebot.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	return nil
}

func (h *Handler) LLMChatNewContextHandler(c telebot.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	return nil
}
