package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	CreateAssistant(ctx context.Context, assistant *Assistant) error
	GetAssistant(ctx context.Context, id uuid.UUID) (*Assistant, error)

	CreateThread(ctx context.Context, thread *Thread) error
	GetThread(ctx context.Context, id uuid.UUID) (*Thread, error)
}

type repository struct {
	db *db.Queries
}

func NewRepository(pool *db.Pool) Repository {
	return &repository{db: db.New(pool)}
}

func (r *repository) CreateAssistant(ctx context.Context, assistant *Assistant) error {
	ast, err := r.db.CreateAssistant(ctx, db.CreateAssistantParams{
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

func (r *repository) GetAssistant(ctx context.Context, id uuid.UUID) (*Assistant, error) {
	ast, err := r.db.GetAssistant(ctx, id)
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

func (r *repository) CreateThread(ctx context.Context, thread *Thread) error {
	th, err := r.db.CreateAssistantThread(ctx, db.CreateAssistantThreadParams{
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

func (r *repository) GetThread(ctx context.Context, id uuid.UUID) (*Thread, error) {
	th, err := r.db.GetAssistantThread(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread: %w", err)
	}
	thread := &Thread{
		Id:          th.Uuid,
		UserId:      th.UserID.Bytes,
		AssistantId: th.AssistantID.Bytes,
		Name:        th.Name,
		Description: th.Description.String,
		Model:       th.Model,
	}
	return thread, nil
}
