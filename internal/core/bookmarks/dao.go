package bookmarks

import (
	"context"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// DAO provides data access operations for bookmarks
type DAO interface {
	CreateBookmark(ctx context.Context, tx db.DBTX, arg db.CreateBookmarkParams) (db.Bookmark, error)
	DeleteBookmark(ctx context.Context, db db.DBTX, arg db.DeleteBookmarkParams) error
	DeleteBookmarksByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) error
	ListBookmarks(ctx context.Context, db db.DBTX, arg db.ListBookmarksParams) ([]db.ListBookmarksRow, error)
	UpdateBookmark(ctx context.Context, db db.DBTX, arg db.UpdateBookmarkParams) (db.Bookmark, error)

	CreateShareContent(ctx context.Context, db db.DBTX, arg db.CreateShareContentParams) (db.ContentShare, error)
	DeleteShareContent(ctx context.Context, db db.DBTX, arg db.DeleteShareContentParams) error
	GetSharedContent(ctx context.Context, db db.DBTX, id uuid.UUID) (db.Content, error)
	GetShareContent(ctx context.Context, db db.DBTX, arg db.GetShareContentParams) (db.ContentShare, error)
	UpdateShareContent(ctx context.Context, db db.DBTX, arg db.UpdateShareContentParams) (db.ContentShare, error)
}
