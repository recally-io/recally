package assistants

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

func (s *Service) CreateEmbedding(ctx context.Context, tx db.DBTX, embedding *EmbeddingDTO) error {
	model := embedding.Dump()

	err := s.dao.CreateAssistantEmbedding(ctx, tx, db.CreateAssistantEmbeddingParams{
		UserID:       model.UserID,
		AttachmentID: model.AttachmentID,
		Text:         model.Text,
		Embeddings:   model.Embeddings,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to create embedding: %w", err)
	}

	return nil
}

func (s *Service) DeleteEmbedding(ctx context.Context, tx db.DBTX, id int32) error {
	if err := s.dao.DeleteAssistantEmbeddings(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to delete embedding: %w", err)
	}

	return nil
}

func (s *Service) DeleteEmbeddingsByAssistant(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) error {
	if err := s.dao.DeleteAssistantEmbeddingsByAssistantId(ctx, tx, pgtype.UUID{Bytes: assistantID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete embeddings by assistant: %w", err)
	}

	return nil
}

func (s *Service) DeleteEmbeddingsByAttachment(ctx context.Context, tx db.DBTX, attachmentID uuid.UUID) error {
	if err := s.dao.DeleteAssistantEmbeddingsByAttachmentId(ctx, tx, pgtype.UUID{Bytes: attachmentID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete embeddings by attachment: %w", err)
	}

	return nil
}

func (s *Service) DeleteEmbeddingsByThread(ctx context.Context, tx db.DBTX, threadID uuid.UUID) error {
	if err := s.dao.DeleteAssistantEmbeddingsByThreadId(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete embeddings by thread: %w", err)
	}

	return nil
}

func (s *Service) SimilaritySearchByThread(ctx context.Context, tx db.DBTX, threadID uuid.UUID, query []float32, limit int32) ([]EmbeddingDTO, error) {
	vec := pgvector.NewVector(query)

	results, err := s.dao.SimilaritySearchByThreadId(ctx, tx, db.SimilaritySearchByThreadIdParams{
		Uuid:       threadID,
		Embeddings: &vec,
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}

	searchResults := make([]EmbeddingDTO, len(results))

	for i, result := range results {
		var v EmbeddingDTO

		v.Load(&result)
		searchResults[i] = v
	}

	return searchResults, nil
}
