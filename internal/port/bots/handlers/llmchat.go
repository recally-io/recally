package handlers

import (
	"context"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"gopkg.in/telebot.v3"
)

func (h *Handler) LLMChatHandler(c telebot.Context) error {
	user := c.Sender()
	text := c.Text()
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	logger.FromContext(ctx).Info("TextHandler", "user", user.Username, "text", text)
	return c.Reply("Hello, " + user.Username)
}

func (h *Handler) LLMChatNewAssistanthandler(c telebot.Context) error {
	var err error
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	user, err := h.getOrCreateUser(ctx)
	if err != nil {
		logger.FromContext(ctx).Error("TextHandler", "error", err)
		return c.Reply("Failed to get or create user " + err.Error())
	}

	assistant := assistants.NewAssistant(user.ID)
	if err := h.assistant.CreateAssistant(ctx, assistant); err != nil {
		return c.Reply("Failed to create assistant " + err.Error())
	}

	_, err = h.updateUser(ctx, *user, assistant.Id, user.ActivateThreadID)
	if err != nil {
		return c.Reply("Failed to update user " + err.Error())
	}

	return c.Reply("Assistant created successfully. Assistant ID: " + assistant.Id.String())
}

func (h *Handler) LLMChatListAssistantshandler(c telebot.Context) error {
	return c.Reply("Not Implemented")
}

func (h *Handler) LLMChatNewThreadHandler(c telebot.Context) error {
	var err error
	ctx := c.Get(constant.ContextKeyContext).(context.Context)
	user, err := h.getOrCreateUser(ctx)
	if err != nil {
		logger.FromContext(ctx).Error("TextHandler", "error", err)
		return c.Reply("Failed to get or create user " + err.Error())
	}

	var assistant *assistants.Assistant

	if user.ActivateAssistantID == uuid.Nil {
		assistant = assistants.NewAssistant(user.ID)
		if err := h.assistant.CreateAssistant(ctx, assistant); err != nil {
			return c.Reply("Failed to create assistant " + err.Error())
		}
		user, err = h.updateUser(ctx, *user, assistant.Id, user.ActivateThreadID)
		if err != nil {
			return c.Reply("Failed to update user " + err.Error())
		}
	} else {
		assistant, err = h.assistant.GetAssistant(ctx, user.ActivateAssistantID.String())
		if err != nil {
			return c.Reply("Failed to get assistant " + err.Error())
		}
	}

	thread := assistants.NewThread(user.ID, *assistant)
	if err := h.assistant.CreateThread(ctx, thread); err != nil {
		return c.Reply("Failed to create thread " + err.Error())
	}

	_, err = h.updateUser(ctx, *user, assistant.Id, thread.Id)
	if err != nil {
		return c.Reply("Failed to update user " + err.Error())
	}

	return c.Reply("Assistant Thread created successfully. Thread ID: " + assistant.Id.String())
}

func (h *Handler) LLMChatListThreadHandler(c telebot.Context) error {
	return c.Reply("Not Implemented")
}
