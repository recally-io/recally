package bookmarks

import (
	"context"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// DAO provides data access operations for bookmarks
type DAO interface {
	IsBookmarkContentExistByURL(ctx context.Context, db db.DBTX, url string) (bool, error)
	CreateBookmarkContent(ctx context.Context, db db.DBTX, arg db.CreateBookmarkContentParams) (db.BookmarkContent, error)
	GetBookmarkContentByID(ctx context.Context, db db.DBTX, id uuid.UUID) (db.BookmarkContent, error)
	GetBookmarkContentByURL(ctx context.Context, db db.DBTX, arg db.GetBookmarkContentByURLParams) (db.BookmarkContent, error)
	UpdateBookmarkContent(ctx context.Context, db db.DBTX, arg db.UpdateBookmarkContentParams) (db.BookmarkContent, error)

	CreateBookmark(ctx context.Context, db db.DBTX, arg db.CreateBookmarkParams) (db.Bookmark, error)
	GetBookmarkWithContent(ctx context.Context, db db.DBTX, arg db.GetBookmarkWithContentParams) (db.GetBookmarkWithContentRow, error)
	ListBookmarks(ctx context.Context, db db.DBTX, arg db.ListBookmarksParams) ([]db.ListBookmarksRow, error)
	SearchBookmarks(ctx context.Context, db db.DBTX, arg db.SearchBookmarksParams) ([]db.SearchBookmarksRow, error)
	DeleteBookmark(ctx context.Context, db db.DBTX, arg db.DeleteBookmarkParams) error
	DeleteBookmarksByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) error

	CreateBookmarkShare(ctx context.Context, db db.DBTX, arg db.CreateBookmarkShareParams) (db.BookmarkShare, error)
	GetBookmarkShareContent(ctx context.Context, db db.DBTX, id uuid.UUID) (db.BookmarkContent, error)
	GetBookmarkShare(ctx context.Context, db db.DBTX, arg db.GetBookmarkShareParams) (db.BookmarkShare, error)
	UpdateBookmarkShareByBookmarkId(ctx context.Context, db db.DBTX, arg db.UpdateBookmarkShareByBookmarkIdParams) (db.BookmarkShare, error)
	DeleteShareContent(ctx context.Context, db db.DBTX, arg db.DeleteShareContentParams) error

	ListExistingBookmarkTagsByTags(ctx context.Context, db db.DBTX, arg db.ListExistingBookmarkTagsByTagsParams) ([]string, error)
	CreateBookmarkTag(ctx context.Context, db db.DBTX, arg db.CreateBookmarkTagParams) (db.BookmarkTag, error)
	ListBookmarkTagsByBookmarkId(ctx context.Context, db db.DBTX, bookmarkID uuid.UUID) ([]string, error)
	LinkBookmarkWithTags(ctx context.Context, db db.DBTX, arg db.LinkBookmarkWithTagsParams) error
	UnLinkBookmarkWithTags(ctx context.Context, db db.DBTX, arg db.UnLinkBookmarkWithTagsParams) error

	ListBookmarkTagsByUser(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]db.ListBookmarkTagsByUserRow, error)
	ListBookmarkDomainsByUser(ctx context.Context, db db.DBTX, userID pgtype.UUID) ([]db.ListBookmarkDomainsByUserRow, error)
}
