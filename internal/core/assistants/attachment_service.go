package assistants

import (
	"context"
	"fmt"
	"recally/internal/core/queue"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/rag/document"

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

	result, err := s.queue.Insert(ctx, queue.AttachmentEmbeddingWorkerArgs{
		AttachmentID: ast.Uuid,
		UserID:       ast.UserID.Bytes,
		Docs:         docs,
	}, nil)
	if err != nil {
		logger.Default.Error("failed to enqueue attachment embedding worker", "err", err)

		return nil, fmt.Errorf("failed to enqueue attachment embedding worker: %w", err)
	} else {
		logger.Default.Info("successfully enqueued attachment embedding worker", "result", result)
	}

	attachment.Load(&ast)

	if attachment.ThreadId != uuid.Nil {
		// update thread to enbale rag
		t, err := s.GetThread(ctx, tx, attachment.ThreadId)
		if err != nil {
			logger.FromContext(ctx).Error("failed to get thread", "err", err)

			return attachment, nil
		}

		if !t.Metadata.RagSettings.Enable {
			t.Metadata.RagSettings.Enable = true
			if _, err := s.UpdateThread(ctx, tx, t); err != nil {
				logger.FromContext(ctx).Error("failed to update thread", "err", err)
			}
		}
	} else {
		// update assistant to enbale rag
		a, err := s.GetAssistant(ctx, tx, attachment.AssistantId)
		if err != nil {
			logger.FromContext(ctx).Error("failed to get assistant", "err", err)

			return attachment, nil
		}

		if !a.Metadata.RagSettings.Enable {
			a.Metadata.RagSettings.Enable = true
			if _, err := s.UpdateAssistant(ctx, tx, a); err != nil {
				logger.FromContext(ctx).Error("failed to update assistant", "err", err)
			}
		}
	}

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
