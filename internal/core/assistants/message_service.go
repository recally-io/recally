package assistants

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

func (s *Service) ListThreadMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID) ([]MessageDTO, error) {
	messages, err := s.dao.ListThreadMessages(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get thread messages: %w", err)
	}

	var result []MessageDTO
	for _, msg := range messages {
		var m MessageDTO
		m.Load(&msg)

		result = append(result, m)
	}

	return result, nil
}

func (s *Service) CreateThreadMessage(ctx context.Context, tx db.DBTX, threadId uuid.UUID, message *MessageDTO) (*MessageDTO, error) {
	model := message.Dump()
	tm, err := s.dao.CreateThreadMessage(ctx, tx, db.CreateThreadMessageParams{
		UserID:      model.UserID,
		AssistantID: model.AssistantID,
		ThreadID:    model.ThreadID,
		Model:       model.Model,
		Role:        model.Role,
		Text:        model.Text,
		Embeddings:  model.Embeddings,
		Metadata:    model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save thread message: %w", err)
	}
	message.Load(&tm)
	return message, nil
}

func (s *Service) GetThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) (*MessageDTO, error) {
	msg, err := s.dao.GetThreadMessage(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thread message: %w", err)
	}
	var m MessageDTO
	m.Load(&msg)
	return &m, nil
}

func (s *Service) AddThreadMessage(ctx context.Context, tx db.DBTX, thread *ThreadDTO, role, text string, metadata MessageMetadata) (*MessageDTO, error) {
	thread.AddMessage(role, text)
	message := &MessageDTO{
		UserID:      thread.UserId,
		AssistantID: thread.AssistantId,
		ThreadID:    thread.Id,
		Model:       thread.Model,
		Role:        role,
		Text:        text,
		Metadata:    metadata,
	}
	return s.CreateThreadMessage(ctx, tx, thread.Id, message)
}

func (s *Service) DeleteThreadMessage(ctx context.Context, tx db.DBTX, id uuid.UUID) error {
	msg, err := s.dao.GetThreadMessage(ctx, tx, id)
	if err != nil {
		return fmt.Errorf("failed to get thread message: %w", err)
	}
	if err := s.dao.DeleteThreadMessageByThreadAndCreatedAt(ctx, tx, db.DeleteThreadMessageByThreadAndCreatedAtParams{
		ThreadID:  msg.ThreadID,
		CreatedAt: msg.CreatedAt,
	}); err != nil {
		return fmt.Errorf("failed to delete thread message: %w", err)
	}
	return nil
}

func (s *Service) DeleteMessagesByAssistant(ctx context.Context, tx db.DBTX, assistantID uuid.UUID) error {
	if err := s.dao.DeleteThreadMessagesByAssistant(ctx, tx, pgtype.UUID{Bytes: assistantID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete thread messages by assistant: %w", err)
	}
	return nil
}

func (s *Service) DeleteMessagesByThread(ctx context.Context, tx db.DBTX, threadID uuid.UUID) error {
	if err := s.dao.DeleteThreadMessagesByThread(ctx, tx, pgtype.UUID{Bytes: threadID, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete thread messages by thread: %w", err)
	}
	return nil
}

func (s *Service) UpdateThreadMessage(ctx context.Context, tx db.DBTX, message *MessageDTO) error {
	dbo := message.Dump()
	err := s.dao.UpdateThreadMessage(ctx, tx, db.UpdateThreadMessageParams{
		Uuid:            dbo.Uuid,
		Role:            dbo.Role,
		Text:            dbo.Text,
		Model:           dbo.Model,
		PromptToken:     dbo.PromptToken,
		CompletionToken: dbo.CompletionToken,
		Embeddings:      dbo.Embeddings,
		Metadata:        dbo.Metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to update thread message: %w", err)
	}
	return nil
}

func (s *Service) SimilaritySearchMessages(ctx context.Context, tx db.DBTX, threadID uuid.UUID, text string, limit int32) ([]MessageDTO, error) {
	embeddings, err := s.llm.CreateEmbeddings(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("failed to create embeddings: %w", err)
	}
	vec := pgvector.NewVector(embeddings)
	messages, err := s.dao.SimilaritySearchMessages(ctx, tx, db.SimilaritySearchMessagesParams{
		ThreadID:   pgtype.UUID{Bytes: threadID, Valid: true},
		Embeddings: &vec,
		Limit:      limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to similarity search messages: %w", err)
	}
	var result []MessageDTO
	for _, msg := range messages {
		var m MessageDTO
		m.Load(&msg)
		result = append(result, m)
	}
	return result, nil
}
