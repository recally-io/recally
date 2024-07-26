package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	CreateAssistant(ctx context.Context, tx db.DBTX, assistant *Assistant) error
	GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*Assistant, error)

	CreateThread(ctx context.Context, tx db.DBTX, thread *Thread) error
	GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*Thread, error)

	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessage, error)
	CreateThreadMessage(ctx context.Context, tx db.DBTX, threadID uuid.UUID, message ThreadMessage) error

	GetTelegramUser(ctx context.Context, tx db.DBTX, userID string) (*User, error)
	CreateTelegramUser(ctx context.Context, tx db.DBTX, userName string, userID string) (*User, error)
	UpdateTelegramUser(ctx context.Context, tx db.DBTX, user User) (*User, error)
}

type repository struct {
	db *db.Queries
}

func NewRepository() Repository {
	return &repository{db: db.New()}
}

func (r *repository) CreateAssistant(ctx context.Context, tx db.DBTX, assistant *Assistant) error {
	ast, err := r.db.CreateAssistant(ctx, tx, db.CreateAssistantParams{
		UserID:       pgtype.UUID{Bytes: assistant.UserId, Valid: true},
		Name:         assistant.Name,
		Description:  pgtype.Text{String: assistant.Description, Valid: assistant.Description != ""},
		SystemPrompt: pgtype.Text{String: assistant.SystemPrompt, Valid: assistant.SystemPrompt != ""},
		Model:        assistant.Model,
	})
	if err != nil {
		return fmt.Errorf("failed to create assistant: %w", err)
	}

	assistant.Id = ast.Uuid
	return nil
}

func (r *repository) GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*Assistant, error) {
	ast, err := r.db.GetAssistant(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant: %w", err)
	}

	assistant := &Assistant{
		Id:           ast.Uuid,
		UserId:       ast.UserID.Bytes,
		Name:         ast.Name,
		Description:  ast.Description.String,
		SystemPrompt: ast.SystemPrompt.String,
		Model:        ast.Model,
	}
	return assistant, nil
}

func (r *repository) CreateThread(ctx context.Context, tx db.DBTX, thread *Thread) error {
	th, err := r.db.CreateAssistantThread(ctx, tx, db.CreateAssistantThreadParams{
		UserID:      pgtype.UUID{Bytes: thread.UserId, Valid: true},
		AssistantID: pgtype.UUID{Bytes: thread.AssistantId, Valid: true},
		Name:        thread.Name,
		Description: pgtype.Text{String: thread.Description, Valid: thread.Description != ""},
		Model:       thread.Model,
	})
	if err != nil {
		return fmt.Errorf("failed to create thread: %w", err)
	}

	thread.Id = th.Uuid

	return nil
}

func (r *repository) GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*Thread, error) {
	th, err := r.db.GetAssistantThread(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	thread := &Thread{
		Id:           th.Uuid,
		UserId:       th.UserID.Bytes,
		AssistantId:  th.AssistantID.Bytes,
		Name:         th.Name,
		Description:  th.Description.String,
		Model:        th.Model,
		SystemPrompt: th.SystemPrompt.String,
	}

	messages, err := r.ListThreadMessages(ctx, tx, th.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	thread.Messages = messages
	return thread, nil
}

func (r *repository) ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessage, error) {
	messages, err := r.db.ListThreadMessages(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	var result []ThreadMessage
	for _, msg := range messages {
		result = append(result, ThreadMessage{
			Role:      msg.Role,
			Text:      msg.Text.String,
			CreatedAt: msg.CreatedAt.Time,
			UpdatedAt: msg.UpdatedAt.Time,
		})
	}

	return result, nil
}

func (r *repository) CreateThreadMessage(ctx context.Context, tx db.DBTX, threadID uuid.UUID, message ThreadMessage) error {
	_, err := r.db.CreateThreadMessage(ctx, tx, db.CreateThreadMessageParams{
		UserID:   pgtype.UUID{Bytes: message.UserID, Valid: true},
		ThreadID: pgtype.UUID{Bytes: threadID, Valid: true},
		Model:    pgtype.Text{String: message.Model, Valid: message.Model != ""},
		Role:     message.Role,
		Text:     pgtype.Text{String: message.Text, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to save thread message: %w", err)
	}
	return nil
}

func (r *repository) GetTelegramUser(ctx context.Context, tx db.DBTX, userID string) (*User, error) {
	user, err := r.db.GetTelegramUser(ctx, tx, pgtype.Text{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  user.Uuid,
		Username:            user.Username.String,
		Telegram:            user.Telegram.String,
		ActivateAssistantID: user.ActivateAssistantID.Bytes,
		ActivateThreadID:    user.ActivateThreadID.Bytes,
	}, nil
}

func (r *repository) CreateTelegramUser(ctx context.Context, tx db.DBTX, userName string, userID string) (*User, error) {
	params := db.InserUserParams{
		Username: pgtype.Text{String: userName, Valid: true},
		Telegram: pgtype.Text{String: userID, Valid: true},
	}
	user, err := r.db.InserUser(ctx, tx, params)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  user.Uuid,
		Username:            user.Username.String,
		Telegram:            user.Telegram.String,
		ActivateAssistantID: user.ActivateAssistantID.Bytes,
	}, nil
}

func (r *repository) UpdateTelegramUser(ctx context.Context, tx db.DBTX, user User) (*User, error) {
	dbUser, err := r.db.UpdateTelegramUser(ctx, tx, db.UpdateTelegramUserParams{
		Telegram:            pgtype.Text{String: user.Telegram, Valid: true},
		ActivateAssistantID: pgtype.UUID{Bytes: user.ActivateAssistantID, Valid: user.ActivateAssistantID != uuid.Nil},
		ActivateThreadID:    pgtype.UUID{Bytes: user.ActivateThreadID, Valid: user.ActivateThreadID != uuid.Nil},
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:                  dbUser.Uuid,
		Username:            dbUser.Username.String,
		Telegram:            dbUser.Telegram.String,
		ActivateAssistantID: dbUser.ActivateAssistantID.Bytes,
		ActivateThreadID:    dbUser.ActivateThreadID.Bytes,
	}, nil
}
