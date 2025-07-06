package bookmarks

import (
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type BookmarkShareDTO struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	BookmarkID uuid.UUID `json:"bookmark_id"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *BookmarkShareDTO) Load(dbo *db.BookmarkShare) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.BookmarkID = dbo.BookmarkID.Bytes
	c.ExpiresAt = dbo.ExpiresAt.Time
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time
}

func (c *BookmarkShareDTO) Dump() db.CreateShareContentParams {
	return db.CreateShareContentParams{
		UserID:    c.UserID,
		ContentID: pgtype.UUID{Bytes: c.BookmarkID, Valid: c.BookmarkID != uuid.Nil},
		ExpiresAt: pgtype.Timestamptz{
			Time:  c.ExpiresAt,
			Valid: !c.ExpiresAt.IsZero(),
		},
	}
}
