package handlers

import (
	"context"
	"fmt"
	"strings"
	"vibrain/internal/core/assistants"
	"vibrain/internal/pkg/auth"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
	"gopkg.in/telebot.v3"
)

func (h *Handler) LLMChatHandler(c telebot.Context) error {
	ctx, user, tx, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	thread, err := h.getActivateThread(ctx, tx, user)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to get thread", "err", err)
		_ = c.Reply("Failed to get thread " + err.Error())
		return err
	}

	if _, err := h.assistantService.AddThreadMessage(ctx, tx, thread, "user", strings.TrimSpace(c.Text())); err != nil {
		logger.FromContext(ctx).Error("Failed to add message to thread", "err", err)
		_ = c.Reply("Failed to add message to thread " + err.Error())
		return err
	}

	message, err := h.assistantService.RunThread(ctx, tx, thread)
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
	ctx, user, tx, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	assistant := assistants.NewAssistant(user.ID)
	if _, err := h.assistantService.CreateAssistant(ctx, tx, assistant); err != nil {
		_ = c.Reply("Failed to create assistant " + err.Error())
		return err
	}

	user.ActivateAssistantID = assistant.Id
	_, err = h.authService.UpdateTelegramUser(ctx, tx, user)
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
	ctx, user, tx, err := h.initHandlerRequest(c)
	if err != nil {
		logger.FromContext(ctx).Error("init request error", "err", err)
		_ = c.Reply("Failed to processing message, please retry.")
		return err
	}

	assistant, err := h.getActivateAssistant(ctx, tx, user)
	if err != nil {
		logger.FromContext(ctx).Error("TextHandler", "error", err)
		_ = c.Reply("Failed to get assistant " + err.Error())
		return err
	}

	thread := assistants.NewThread(user.ID, *assistant)
	if _, err := h.assistantService.CreateThread(ctx, tx, thread); err != nil {
		_ = c.Reply("Failed to create thread " + err.Error())
		return err
	}

	user.ActivateThreadID = thread.Id
	user.ActivateAssistantID = assistant.Id
	_, err = h.authService.UpdateTelegramUser(ctx, tx, user)
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

func (h *Handler) getActivateAssistant(ctx context.Context, tx db.DBTX, user *auth.UserDTO) (*assistants.AssistantDTO, error) {
	assistant, err := h.assistantService.GetAssistant(ctx, tx, user.ActivateAssistantID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			assistant = assistants.NewAssistant(user.ID)
			if _, err := h.assistantService.CreateAssistant(ctx, tx, assistant); err != nil {
				return nil, fmt.Errorf("failed to create assistant: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get assistant: %w", err)
		}
	}
	return assistant, nil
}

func (h *Handler) getActivateThread(ctx context.Context, tx db.DBTX, user *auth.UserDTO) (*assistants.ThreadDTO, error) {
	thread, err := h.assistantService.GetThread(ctx, tx, user.ActivateThreadID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			assistant, err := h.getActivateAssistant(ctx, tx, user)
			if err != nil {
				return nil, fmt.Errorf("failed to get activate assistant when : %w", err)
			}
			thread = assistants.NewThread(user.ID, *assistant)
			if _, err := h.assistantService.CreateThread(ctx, tx, thread); err != nil {
				return nil, fmt.Errorf("failed to create thread: %w", err)
			}
			user.ActivateThreadID = thread.Id
			user.ActivateAssistantID = assistant.Id
			_, err = h.authService.UpdateTelegramUser(ctx, tx, user)
			if err != nil {
				return nil, fmt.Errorf("failed to update user: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get thread: %w", err)
		}
	}
	return thread, nil
}

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
		if db.IsNotFound(err) {
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
