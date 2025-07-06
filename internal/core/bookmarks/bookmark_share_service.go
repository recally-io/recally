package bookmarks

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID, expiresAt time.Time) (*BookmarkShareDTO, error) {
	cs, err := s.dao.CreateBookmarkShare(ctx, tx, db.CreateBookmarkShareParams{
		UserID:     userID,
		BookmarkID: pgtype.UUID{Bytes: bookmarkID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: !expiresAt.IsZero(),
		},
	})

	var dto BookmarkShareDTO

	dto.Load(&cs)

	return &dto, err
}

func (s *Service) GetBookmarkShareContent(ctx context.Context, tx db.DBTX, sharedID uuid.UUID) (*BookmarkContentDTO, error) {
	sharedContent, err := s.dao.GetBookmarkShareContent(ctx, tx, sharedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared content: %w", err)
	}

	var dto BookmarkContentDTO

	dto.Load(&sharedContent)

	return &dto, nil
}

func (s *Service) GetBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID) (*BookmarkShareDTO, error) {
	sharedContent, err := s.dao.GetBookmarkShare(ctx, tx, db.GetBookmarkShareParams{
		BookmarkID: pgtype.UUID{Bytes: bookmarkID, Valid: true},
		UserID:     userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get shared content: %w", err)
	}

	var dto BookmarkShareDTO

	dto.Load(&sharedContent)

	return &dto, nil
}

func (s *Service) UpdateBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID, expiresAt time.Time) (*BookmarkShareDTO, error) {
	sc, err := s.dao.UpdateBookmarkShareByBookmarkId(ctx, tx, db.UpdateBookmarkShareByBookmarkIdParams{
		ID:     bookmarkID,
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: !expiresAt.IsZero(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update shared content: %w", err)
	}

	var dto BookmarkShareDTO

	dto.Load(&sc)

	return &dto, nil
}

func (s *Service) DeleteBookmarkShare(ctx context.Context, tx db.DBTX, userID, bookmarkID uuid.UUID) error {
	if err := s.dao.DeleteShareContent(ctx, tx, db.DeleteShareContentParams{
		ID:     bookmarkID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("failed to delete shared content: %w", err)
	}

	return nil
}
