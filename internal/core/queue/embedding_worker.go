package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/rag/document"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"github.com/riverqueue/river"
	"golang.org/x/sync/errgroup"
)

type attachmentEmbeddingWorkerDao interface {
	CreateAssistantEmbedding(ctx context.Context, db db.DBTX, arg db.CreateAssistantEmbeddingParams) error
	IsAssistantEmbeddingExists(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (bool, error)
}

type AttachmentEmbeddingWorkerArgs struct {
	UserID       uuid.UUID           `json:"user_id"`
	AttachmentID uuid.UUID           `json:"attachment_id"`
	Docs         []document.Document `json:"docs"`
}

func (AttachmentEmbeddingWorkerArgs) Kind() string {
	return "create_assistant_attachment_embeddings"
}

type AttachmentEmbeddingWorker struct {
	river.WorkerDefaults[AttachmentEmbeddingWorkerArgs]
	llm    *llms.LLM
	dao    attachmentEmbeddingWorkerDao
	dbPool *pgxpool.Pool
}

func NewAttachmentEmbeddingWorker(llm *llms.LLM, dao attachmentEmbeddingWorkerDao, dbPool *pgxpool.Pool) *AttachmentEmbeddingWorker {
	return &AttachmentEmbeddingWorker{
		llm:    llm,
		dao:    dao,
		dbPool: dbPool,
	}
}

func (w *AttachmentEmbeddingWorker) Work(ctx context.Context, args *river.Job[AttachmentEmbeddingWorkerArgs]) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(10) // Set max concurrent goroutines to 10

	for _, doc := range args.Args.Docs {
		doc := doc

		eg.Go(func() error {
			return w.createAssistantEmbedding(ctx, doc, args.Args.UserID, args.Args.AttachmentID)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("AttachmentEmbeddingWorker: failed to create assistant embedding: %w", err)
	}

	return nil
}

func (w *AttachmentEmbeddingWorker) createAssistantEmbedding(ctx context.Context, doc document.Document, userID, attachmentID uuid.UUID) error {
	tx, err := w.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AttachmentEmbeddingWorker: failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Commit(ctx); err != nil {
			fmt.Printf("AttachmentEmbeddingWorker: failed to commit transaction: %v\n", err)
		}
	}()

	if doc.ID != uuid.Nil {
		exists, err := w.dao.IsAssistantEmbeddingExists(ctx, tx, doc.ID)
		if err != nil {
			return fmt.Errorf("AttachmentEmbeddingWorker: failed to check if assistant embedding exists: %w", err)
		}

		if exists {
			logger.Default.Debug("AttachmentEmbeddingWorker: assistant embedding already exists", "id", doc.ID)

			return nil
		}
	} else {
		doc.ID = uuid.New()
	}

	embeddings, err := w.llm.CreateEmbeddings(ctx, doc.Content)
	if err != nil {
		return fmt.Errorf("AttachmentEmbeddingWorker: failed to create embeddings: %w", err)
	}

	doc.Metadata["id"] = doc.ID
	doc.Metadata["attachment_id"] = attachmentID

	metadata, err := json.Marshal(doc.Metadata)
	if err != nil {
		return fmt.Errorf("AttachmentEmbeddingWorker: failed to marshal metadata: %w", err)
	}

	vec := pgvector.NewVector(embeddings)

	arg := db.CreateAssistantEmbeddingParams{
		Uuid:         doc.ID,
		UserID:       pgtype.UUID{Bytes: [16]byte(userID), Valid: userID != uuid.Nil},
		AttachmentID: pgtype.UUID{Bytes: [16]byte(attachmentID), Valid: attachmentID != uuid.Nil},
		Text:         doc.Content,
		Embeddings:   &vec,
		Metadata:     metadata,
	}
	if err := w.dao.CreateAssistantEmbedding(ctx, tx, arg); err != nil {
		return fmt.Errorf("AttachmentEmbeddingWorker: failed to create assistant embedding: %w", err)
	}

	logger.Default.Debug("AttachmentEmbeddingWorker: assistant embedding created", "id", doc.ID)

	return nil
}
