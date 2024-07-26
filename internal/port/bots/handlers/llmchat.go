package handlers

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	"gopkg.in/telebot.v3"
)

func (h *Handler) LLMChatHandler(c telebot.Context) error {
	ctx := c.Get(constant.ContextKeyContext).(context.Context)

	user, err := h.getOrCreateUser(ctx)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to get or create user", "err", err)
		return c.Reply("Failed to get or create user " + err.Error())
	}

	thread, err := h.getActivateThread(ctx, user)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to get thread", "err", err)
		return c.Reply("Failed to get thread " + err.Error())
	}

	if err := h.assistant.AddThreadMessage(ctx, thread, "user", strings.TrimSpace(c.Text())); err != nil {
		logger.FromContext(ctx).Error("Failed to add message to thread", "err", err)
		return c.Reply("Failed to add message to thread " + err.Error())
	}

	message, err := h.assistant.RunThread(ctx, thread)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to run thread", "err", err)
		return c.Reply("Failed to run thread " + err.Error())
	}
	md := convertToTGMarkdown(message.Text)
	return c.Reply(md, telebot.ModeMarkdownV2)
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

	user.ActivateAssistantID = assistant.Id
	_, err = h.repository.UpdateUser(ctx, *user)
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

	assistant, err := h.getActivateAssistant(ctx, user)
	if err != nil {
		logger.FromContext(ctx).Error("TextHandler", "error", err)
		return c.Reply("Failed to get assistant " + err.Error())
	}

	thread := assistants.NewThread(user.ID, *assistant)
	if err := h.assistant.CreateThread(ctx, thread); err != nil {
		return c.Reply("Failed to create thread " + err.Error())
	}

	user.ActivateThreadID = thread.Id
	user.ActivateAssistantID = assistant.Id
	_, err = h.repository.UpdateUser(ctx, *user)
	if err != nil {
		return c.Reply("Failed to update user " + err.Error())
	}

	return c.Reply("Assistant Thread created successfully. Thread ID: " + assistant.Id.String())
}

func (h *Handler) LLMChatListThreadHandler(c telebot.Context) error {
	return c.Reply("Not Implemented")
}

func (h *Handler) getActivateAssistant(ctx context.Context, user *User) (*assistants.Assistant, error) {
	assistant, err := h.assistant.GetAssistant(ctx, user.ActivateAssistantID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			assistant = assistants.NewAssistant(user.ID)
			if err := h.assistant.CreateAssistant(ctx, assistant); err != nil {
				return nil, fmt.Errorf("failed to create assistant: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get assistant: %w", err)
		}
	}
	return assistant, nil
}

func (h *Handler) getActivateThread(ctx context.Context, user *User) (*assistants.Thread, error) {
	thread, err := h.assistant.GetThread(ctx, user.ActivateThreadID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			assistant, err := h.getActivateAssistant(ctx, user)
			if err != nil {
				return nil, fmt.Errorf("failed to get activate assistant when : %w", err)
			}
			thread = assistants.NewThread(user.ID, *assistant)
			if err := h.assistant.CreateThread(ctx, thread); err != nil {
				return nil, fmt.Errorf("failed to create thread: %w", err)
			}
			user.ActivateThreadID = thread.Id
			user.ActivateAssistantID = assistant.Id
			_, err = h.repository.UpdateUser(ctx, *user)
			if err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get thread: %w", err)
		}
	}
	return thread, nil
}
