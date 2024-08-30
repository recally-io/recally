package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"

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

func (s *Service) SimilaritySearchByThread(ctx context.Context, tx db.DBTX, threadID uuid.UUID, query []float32, limit int32) ([]SimilaritySearchResult, error) {
	results, err := s.dao.SimilaritySearchByThreadId(ctx, tx, db.SimilaritySearchByThreadIdParams{
		Uuid:       threadID,
		Embeddings: pgvector.NewVector(query),
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}

	searchResults := make([]SimilaritySearchResult, len(results))
	for i, result := range results {
		searchResults[i] = SimilaritySearchResult{
			ID:       result.ID,
			Text:     result.Text,
			Metadata: result.Metadata,
			Score:    float64(result.Score),
		}
	}

	return searchResults, nil
}

type SimilaritySearchResult struct {
	ID       int32   `json:"id"`
	Text     string  `json:"text"`
	Metadata []byte  `json:"metadata"`
	Score    float64 `json:"score"`
}
