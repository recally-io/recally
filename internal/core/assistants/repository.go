package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	CreateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) error
	ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]AssistantDTO, error)
	GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AssistantDTO, error)
	// DeleteAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) error

	ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]ThreadDTO, error)
	CreateThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) error
	GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadDTO, error)
	// DeleteThread(ctx context.Context, tx db.DBTX, id uuid.UUID) error

	ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessageDTO, error)
	CreateThreadMessage(ctx context.Context, tx db.DBTX, threadID uuid.UUID, message ThreadMessageDTO) error

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

func (r *repository) ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]AssistantDTO, error) {
	asts, err := r.db.ListAssistantsByUser(ctx, tx, pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get assistants: %w", err)
	}
	asstants := make([]AssistantDTO, 0, len(asts))
	for _, ast := range asts {
		var a AssistantDTO
		a.Load(&ast)
		asstants = append(asstants, a)
	}
	return asstants, nil
}

func (r *repository) CreateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) error {
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

func (r *repository) GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AssistantDTO, error) {
	ast, err := r.db.GetAssistant(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant: %w", err)
	}
	var assistant AssistantDTO
	assistant.Load(&ast)
	return &assistant, nil
}

func (r *repository) CreateThread(ctx context.Context, tx db.DBTX, thread *ThreadDTO) error {
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

func (r *repository) ListThreads(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) ([]ThreadDTO, error) {
	threads, err := r.db.ListAssistantThreads(ctx, tx, pgtype.UUID{Bytes: assistantID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get threads: %w", err)
	}

	var result []ThreadDTO
	for _, th := range threads {
		var t ThreadDTO
		t.Load(&th)
		result = append(result, t)
	}

	return result, nil
}

func (r *repository) GetThread(ctx context.Context, tx db.DBTX, id uuid.UUID) (*ThreadDTO, error) {
	th, err := r.db.GetAssistantThread(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	var t ThreadDTO
	t.Load(&th)

	messages, err := r.ListThreadMessages(ctx, tx, th.Uuid)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}
	t.Messages = messages
	return &t, nil
}

func (r *repository) ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]ThreadMessageDTO, error) {
	messages, err := r.db.ListThreadMessages(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	var result []ThreadMessageDTO
	for _, msg := range messages {
		var m ThreadMessageDTO
		m.Load(&msg)

		result = append(result, m)
	}

	return result, nil
}

func (r *repository) CreateThreadMessage(ctx context.Context, tx db.DBTX, threadID uuid.UUID, message ThreadMessageDTO) error {
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
