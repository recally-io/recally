package bookmarks

import (
	"context"
	"fmt"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) ShareContent(ctx context.Context, tx db.DBTX, contentID uuid.UUID, expiresAt time.Time) (*ContentShareDTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cs, err := s.dao.CreateShareContent(ctx, tx, db.CreateShareContentParams{
		UserID:    user.ID,
		ContentID: pgtype.UUID{Bytes: contentID, Valid: true},
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: !expiresAt.IsZero(),
		},
	})

	var dto ContentShareDTO
	dto.Load(&cs)
	return &dto, err
}

func (s *Service) GetSharedContent(ctx context.Context, tx db.DBTX, sharedID uuid.UUID) (*ContentDTO, error) {
	sharedContent, err := s.dao.GetSharedContent(ctx, tx, sharedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared content: %w", err)
	}

	var dto ContentDTO
	dto.Load(&sharedContent)
	return &dto, nil
}

func (s *Service) GetShareContent(ctx context.Context, tx db.DBTX, contentID uuid.UUID) (*ContentShareDTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	sharedContent, err := s.dao.GetShareContent(ctx, tx, db.GetShareContentParams{
		ContentID: pgtype.UUID{Bytes: contentID, Valid: true},
		UserID:    user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get shared content: %w", err)
	}

	var dto ContentShareDTO
	dto.Load(&sharedContent)
	return &dto, nil
}

func (s *Service) UpdateSharedContent(ctx context.Context, tx db.DBTX, contentID uuid.UUID, expiresAt time.Time) (*ContentShareDTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	sc, err := s.dao.UpdateShareContent(ctx, tx, db.UpdateShareContentParams{
		ID:     contentID,
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: !expiresAt.IsZero(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update shared content: %w", err)
	}

	var dto ContentShareDTO
	dto.Load(&sc)
	return &dto, nil
}

func (s *Service) DeleteSharedContent(ctx context.Context, tx db.DBTX, contentID uuid.UUID) error {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return err
	}

	if err := s.dao.DeleteShareContent(ctx, tx, db.DeleteShareContentParams{
		ID:     contentID,
		UserID: user.ID,
	}); err != nil {
		return fmt.Errorf("failed to delete shared content: %w", err)
	}
	return err
}
