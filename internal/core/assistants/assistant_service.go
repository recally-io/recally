package assistants

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) ListAssistants(ctx context.Context, tx db.DBTX, userId uuid.UUID) ([]AssistantDTO, error) {
	asts, err := s.dao.ListAssistantsByUser(ctx, tx, pgtype.UUID{Bytes: userId, Valid: true})
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

func (s *Service) CreateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) (*AssistantDTO, error) {
	model := assistant.Dump()

	ast, err := s.dao.CreateAssistant(ctx, tx, db.CreateAssistantParams{
		UserID:       model.UserID,
		Name:         model.Name,
		Description:  model.Description,
		SystemPrompt: model.SystemPrompt,
		Model:        model.Model,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create assistant: %w", err)
	}

	assistant.Load(&ast)

	return assistant, nil
}

func (s *Service) UpdateAssistant(ctx context.Context, tx db.DBTX, assistant *AssistantDTO) (*AssistantDTO, error) {
	model := assistant.Dump()

	ast, err := s.dao.UpdateAssistant(ctx, tx, db.UpdateAssistantParams{
		Uuid:         assistant.Id,
		Name:         model.Name,
		Description:  model.Description,
		SystemPrompt: model.SystemPrompt,
		Model:        model.Model,
		Metadata:     model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update assistant: %w", err)
	}

	assistant.Load(&ast)

	return assistant, nil
}

func (s *Service) GetAssistant(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AssistantDTO, error) {
	ast, err := s.dao.GetAssistant(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get assistant: %w", err)
	}

	var assistant AssistantDTO

	assistant.Load(&ast)

	return &assistant, nil
}

func (s *Service) DeleteAssistant(ctx context.Context, tx db.DBTX, assistantId uuid.UUID) error {
	// Delete associated attachments
	if err := s.dao.DeleteAssistantAttachmentsByAssistantId(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete assistant attachments: %w", err)
	}
	// Delete associated threads and messages
	if err := s.dao.DeleteThreadMessagesByAssistant(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete thread messages by assistant: %w", err)
	}

	if err := s.dao.DeleteAssistantThreadsByAssistant(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true}); err != nil {
		return fmt.Errorf("failed to delete assistant threads: %w", err)
	}

	// Delete the assistant
	if err := s.dao.DeleteAssistant(ctx, tx, assistantId); err != nil {
		return fmt.Errorf("failed to delete assistant: %w", err)
	}

	return nil
}
