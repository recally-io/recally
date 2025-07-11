package auth

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) OwnerTransfer(ctx context.Context, tx db.DBTX, ownerID, newOwnerID uuid.UUID) error {
	// Transfer ownership of bookmarks
	if err := s.dao.OwnerTransferBookmark(ctx, tx, db.OwnerTransferBookmarkParams{
		UserID:    pgtype.UUID{Bytes: ownerID, Valid: ownerID != uuid.Nil},
		NewUserID: pgtype.UUID{Bytes: newOwnerID, Valid: newOwnerID != uuid.Nil},
	}); err != nil {
		return fmt.Errorf("failed to transfer ownership of bookmarks: %w", err)
	}

	// Transfer ownership of bookmark content
	if err := s.dao.OwnerTransferBookmarkContent(ctx, tx, db.OwnerTransferBookmarkContentParams{
		UserID:    pgtype.UUID{Bytes: ownerID, Valid: ownerID != uuid.Nil},
		NewUserID: pgtype.UUID{Bytes: newOwnerID, Valid: newOwnerID != uuid.Nil},
	}); err != nil {
		return fmt.Errorf("failed to transfer ownership of bookmark content: %w", err)
	}

	// Transfer ownership of bookmark shares
	if err := s.dao.OwnerTransferBookmarkShare(ctx, tx, db.OwnerTransferBookmarkShareParams{
		UserID:    pgtype.UUID{Bytes: ownerID, Valid: ownerID != uuid.Nil},
		NewUserID: pgtype.UUID{Bytes: newOwnerID, Valid: newOwnerID != uuid.Nil},
	}); err != nil {
		return fmt.Errorf("failed to transfer ownership of bookmark shares: %w", err)
	}

	// Transfer ownership of bookmark tags
	if err := s.dao.OwnerTransferBookmarkTag(ctx, tx, db.OwnerTransferBookmarkTagParams{
		UserID:    pgtype.UUID{Bytes: ownerID, Valid: ownerID != uuid.Nil},
		NewUserID: pgtype.UUID{Bytes: newOwnerID, Valid: newOwnerID != uuid.Nil},
	}); err != nil {
		return fmt.Errorf("failed to transfer ownership of bookmark tags: %w", err)
	}

	logger.FromContext(ctx).Info("all owner transfer tasks completed successfully")

	return nil
}
