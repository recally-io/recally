package assistants

import (
	"context"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type dao interface {
	ListAssistantsByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) ([]db.Assistant, error)
	CreateAssistant(ctx context.Context, db db.DBTX, arg db.CreateAssistantParams) (db.Assistant, error)
	GetAssistant(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.Assistant, error)
	UpdateAssistant(ctx context.Context, db db.DBTX, arg db.UpdateAssistantParams) (db.Assistant, error)
	DeleteAssistant(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error

	CreateAssistantThread(ctx context.Context, db db.DBTX, arg db.CreateAssistantThreadParams) (db.AssistantThread, error)
	ListAssistantThreads(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) ([]db.AssistantThread, error)
	GetAssistantThread(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.AssistantThread, error)
	UpdateAssistantThread(ctx context.Context, db db.DBTX, arg db.UpdateAssistantThreadParams) (db.AssistantThread, error)
	DeleteAssistantThread(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	DeleteAssistantThreadsByAssistant(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error

	ListThreadMessages(ctx context.Context, db db.DBTX, threadID pgtype.UUID) ([]db.AssistantMessage, error)
	CreateThreadMessage(ctx context.Context, db db.DBTX, arg db.CreateThreadMessageParams) (db.AssistantMessage, error)
	GetThreadMessage(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.AssistantMessage, error)
	UpdateThreadMessage(ctx context.Context, db db.DBTX, arg db.UpdateThreadMessageParams) error
	DeleteThreadMessage(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	DeleteThreadMessagesByThread(ctx context.Context, db db.DBTX, threadID pgtype.UUID) error
	DeleteThreadMessageByThreadAndCreatedAt(ctx context.Context, db db.DBTX, arg db.DeleteThreadMessageByThreadAndCreatedAtParams) error
	DeleteThreadMessagesByAssistant(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error
	SimilaritySearchMessages(ctx context.Context, db db.DBTX, arg db.SimilaritySearchMessagesParams) ([]db.AssistantMessage, error)

	CreateAssistantAttachment(ctx context.Context, db db.DBTX, arg db.CreateAssistantAttachmentParams) (db.AssistantAttachment, error)
	DeleteAssistantAttachment(ctx context.Context, db db.DBTX, argUuid uuid.UUID) error
	DeleteAssistantAttachmentsByAssistantId(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error
	DeleteAssistantAttachmentsByThreadId(ctx context.Context, db db.DBTX, threadID pgtype.UUID) error
	GetAssistantAttachmentById(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.AssistantAttachment, error)
	ListAssistantAttachmentsByAssistantId(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) ([]db.AssistantAttachment, error)
	ListAssistantAttachmentsByThreadId(ctx context.Context, db db.DBTX, threadID pgtype.UUID) ([]db.AssistantAttachment, error)
	UpdateAssistantAttachment(ctx context.Context, db db.DBTX, arg db.UpdateAssistantAttachmentParams) error

	CreateAssistantEmbedding(ctx context.Context, db db.DBTX, arg db.CreateAssistantEmbeddingParams) error
	DeleteAssistantEmbeddings(ctx context.Context, db db.DBTX, id int32) error
	DeleteAssistantEmbeddingsByAssistantId(ctx context.Context, db db.DBTX, assistantID pgtype.UUID) error
	DeleteAssistantEmbeddingsByAttachmentId(ctx context.Context, db db.DBTX, attachmentID pgtype.UUID) error
	DeleteAssistantEmbeddingsByThreadId(ctx context.Context, db db.DBTX, threadID pgtype.UUID) error
	SimilaritySearchByThreadId(ctx context.Context, db db.DBTX, arg db.SimilaritySearchByThreadIdParams) ([]db.AssistantEmbeddding, error)
}
