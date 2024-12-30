package auth

import (
	"context"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) OwnerTransfer(ctx context.Context, tx db.DBTX, ownerID, newOwnerID uuid.UUID) error {
	return s.dao.OwnerTransferBookmark(ctx, tx, db.OwnerTransferBookmarkParams{
		UserID:   pgtype.UUID{Bytes: ownerID, Valid: ownerID != uuid.Nil},
		UserID_2: pgtype.UUID{Bytes: newOwnerID, Valid: newOwnerID != uuid.Nil},
	})
}
