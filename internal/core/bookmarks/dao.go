package bookmarks

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// DAO provides data access operations for bookmarks
type DAO interface {
	CreateBookmark(ctx context.Context, tx db.DBTX, arg db.CreateBookmarkParams) (db.Bookmark, error)
	DeleteBookmark(ctx context.Context, db db.DBTX, arg db.DeleteBookmarkParams) error
	DeleteBookmarksByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) error
	GetBookmarkByUUID(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.Bookmark, error)
	GetBookmarkByURL(ctx context.Context, db db.DBTX, arg db.GetBookmarkByURLParams) (db.Bookmark, error)
	ListBookmarks(ctx context.Context, db db.DBTX, arg db.ListBookmarksParams) ([]db.Bookmark, error)
	UpdateBookmark(ctx context.Context, db db.DBTX, arg db.UpdateBookmarkParams) (db.Bookmark, error)
}
