package assistants

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type dao interface {
	ListAssistantsByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) ([]db.Assistant, error)
	CreateAssistant(ctx context.Context, db db.DBTX, arg db.CreateAssistantParams) (db.Assistant, error)
	GetAssistant(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.Assistant, error)
	UpdateAssistant(ctx context.Context, db db.DBTX, arg db.UpdateAssistantParams) (db.Assistant, error)

	CreateAssistantThread(ctx context.Context, db db.DBTX, arg db.CreateAssistantThreadParams) (db.AssistantThread, error)
	ListAssistantThreads(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) ([]db.AssistantThread, error)
	GetAssistantThread(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.AssistantThread, error)
	UpdateAssistantThread(ctx context.Context, db db.DBTX, arg db.UpdateAssistantThreadParams) (db.AssistantThread, error)
	DeleteAssistantThread(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	DeleteAssistantThreadsByAssistant(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error

	ListThreadMessages(ctx context.Context, db db.DBTX, threadID pgtype.UUID) ([]db.AssistantMessage, error)
	CreateThreadMessage(ctx context.Context, db db.DBTX, arg db.CreateThreadMessageParams) (db.AssistantMessage, error)
	GetThreadMessage(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.AssistantMessage, error)
	DeleteThreadMessage(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	DeleteThreadMessagesByThread(ctx context.Context, db db.DBTX, threadID pgtype.UUID) error
	DeleteThreadMessageByThreadAndCreatedAt(ctx context.Context, db db.DBTX, arg db.DeleteThreadMessageByThreadAndCreatedAtParams) error
	DeleteThreadMessagesByAssistant(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error
}
