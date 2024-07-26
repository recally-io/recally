package handlers

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/logger"

	"gopkg.in/telebot.v3"
)

func (h *Handler) LLMChatHandler(c telebot.Context) error {
	ctx, user, repo, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	thread, err := h.getActivateThread(ctx, repo, user)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to get thread", "err", err)
		_ = c.Reply("Failed to get thread " + err.Error())
		return err
	}

	if err := h.assistant.AddThreadMessage(ctx, thread, "user", strings.TrimSpace(c.Text())); err != nil {
		logger.FromContext(ctx).Error("Failed to add message to thread", "err", err)
		_ = c.Reply("Failed to add message to thread " + err.Error())
		return err
	}

	message, err := h.assistant.RunThread(ctx, thread)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to run thread", "err", err)
		_ = c.Reply("Failed to run thread " + err.Error())
		return err
	}
	md := convertToTGMarkdown(message.Text)
	_ = c.Reply(md, telebot.ModeMarkdownV2)
	return nil
}

func (h *Handler) LLMChatNewAssistanthandler(c telebot.Context) error {
	ctx, user, repo, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	assistant := assistants.NewAssistant(user.ID)
	if err := h.assistant.CreateAssistant(ctx, assistant); err != nil {
		_ = c.Reply("Failed to create assistant " + err.Error())
		return err
	}

	user.ActivateAssistantID = assistant.Id
	_, err = repo.UpdateUser(ctx, *user)
	if err != nil {
		_ = c.Reply("Failed to update user " + err.Error())
		return err
	}

	_ = c.Reply("Assistant created successfully. Assistant ID: " + assistant.Id.String())
	return nil
}

func (h *Handler) LLMChatListAssistantshandler(c telebot.Context) error {
	return c.Reply("Not Implemented")
}

func (h *Handler) LLMChatNewThreadHandler(c telebot.Context) error {
	ctx, user, repo, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	assistant, err := h.getActivateAssistant(ctx, user)
	if err != nil {
		logger.FromContext(ctx).Error("TextHandler", "error", err)
		_ = c.Reply("Failed to get assistant " + err.Error())
		return err
	}

	thread := assistants.NewThread(user.ID, *assistant)
	if err := h.assistant.CreateThread(ctx, thread); err != nil {
		_ = c.Reply("Failed to create thread " + err.Error())
		return err
	}

	user.ActivateThreadID = thread.Id
	user.ActivateAssistantID = assistant.Id
	_, err = repo.UpdateUser(ctx, *user)
	if err != nil {
		_ = c.Reply("Failed to update user " + err.Error())
		return err
	}

	_ = c.Reply("Assistant Thread created successfully. Thread ID: " + assistant.Id.String())
	return nil
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

func (h *Handler) getActivateThread(ctx context.Context, repo Repository, user *User) (*assistants.Thread, error) {
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
			_, err = repo.UpdateUser(ctx, *user)
			if err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get thread: %w", err)
		}
	}
	return thread, nil
}

func (h *Handler) initHandlerRequest(c telebot.Context) (context.Context, *User, Repository, error) {
	ctx := c.Get(contexts.ContextKeyContext).(context.Context)
	repo, err := NewRepositoryFromContext(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get repository: %w", err)
	}
	user, err := repo.GetOrCreateUser(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get or create user: %w", err)
	}

	return ctx, user, repo, nil
}
