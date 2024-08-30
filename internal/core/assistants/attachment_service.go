package assistants

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/rag/document"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateAttachment(ctx context.Context, tx db.DBTX, attachment *AttachmentDTO, docs []document.Document) (*AttachmentDTO, error) {
	model := attachment.Dump()
	ast, err := s.dao.CreateAssistantAttachment(ctx, tx, db.CreateAssistantAttachmentParams{
		UserID:      model.UserID,
		AssistantID: model.AssistantID,
		ThreadID:    model.ThreadID,
		Name:        model.Name,
		Type:        model.Type,
		Url:         model.Url,
		Size:        model.Size,
		Metadata:    model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}
	attachment.Load(&ast)
	return attachment, nil
}

func (s *Service) GetAttachment(ctx context.Context, tx db.DBTX, id uuid.UUID) (*AttachmentDTO, error) {
	ast, err := s.dao.GetAssistantAttachmentById(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}
	var attachment AttachmentDTO
	attachment.Load(&ast)
	return &attachment, nil
}

func (s *Service) UpdateAttachment(ctx context.Context, tx db.DBTX, attachment *AttachmentDTO) (*AttachmentDTO, error) {
	model := attachment.Dump()
	err := s.dao.UpdateAssistantAttachment(ctx, tx, db.UpdateAssistantAttachmentParams{
		Uuid:     attachment.Id,
		Name:     model.Name,
		Type:     model.Type,
		Url:      model.Url,
		Size:     model.Size,
		Metadata: model.Metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update attachment: %w", err)
	}
	return attachment, nil
}

func (s *Service) DeleteAttachment(ctx context.Context, tx db.DBTX, attachmentId uuid.UUID) error {
	if err := s.dao.DeleteAssistantAttachment(ctx, tx, attachmentId); err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}
	return nil
}

func (s *Service) ListAttachmentsByAssistant(ctx context.Context, tx db.DBTX, assistantId uuid.UUID) ([]AttachmentDTO, error) {
	asts, err := s.dao.ListAssistantAttachmentsByAssistantId(ctx, tx, pgtype.UUID{Bytes: assistantId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments by assistant: %w", err)
	}
	attachments := make([]AttachmentDTO, 0, len(asts))
	for _, ast := range asts {
		var a AttachmentDTO
		a.Load(&ast)
		attachments = append(attachments, a)
	}
	return attachments, nil
}

func (s *Service) ListAttachmentsByThread(ctx context.Context, tx db.DBTX, threadId uuid.UUID) ([]AttachmentDTO, error) {
	asts, err := s.dao.ListAssistantAttachmentsByThreadId(ctx, tx, pgtype.UUID{Bytes: threadId, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list attachments by thread: %w", err)
	}
	attachments := make([]AttachmentDTO, 0, len(asts))
	for _, ast := range asts {
		var a AttachmentDTO
		a.Load(&ast)
		attachments = append(attachments, a)
	}
	return attachments, nil
}
