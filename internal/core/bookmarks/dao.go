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
	GetBookmarkByUUID(ctx context.Context, db db.DBTX, argUuid uuid.UUID) (db.Bookmark, error)
	GetBookmarkByURL(ctx context.Context, db db.DBTX, arg db.GetBookmarkByURLParams) (db.Bookmark, error)
	ListBookmarks(ctx context.Context, db db.DBTX, arg db.ListBookmarksParams) ([]db.ListBookmarksRow, error)
	UpdateBookmark(ctx context.Context, db db.DBTX, arg db.UpdateBookmarkParams) (db.Bookmark, error)

	CreateContent(ctx context.Context, db db.DBTX, arg db.CreateContentParams) (db.Content, error)
	CreateContentTag(ctx context.Context, db db.DBTX, arg db.CreateContentTagParams) (db.ContentTag, error)
	DeleteContent(ctx context.Context, db db.DBTX, arg db.DeleteContentParams) error
	DeleteContentTag(ctx context.Context, db db.DBTX, arg db.DeleteContentTagParams) error
	DeleteContentsByUser(ctx context.Context, db db.DBTX, userID uuid.UUID) error
	GetContent(ctx context.Context, db db.DBTX, arg db.GetContentParams) (db.GetContentRow, error)
	IsContentExistWithURL(ctx context.Context, db db.DBTX, arg db.IsContentExistWithURLParams) (bool, error)
	LinkContentWithTags(ctx context.Context, db db.DBTX, arg db.LinkContentWithTagsParams) error
	ListContentTags(ctx context.Context, db db.DBTX, arg db.ListContentTagsParams) (interface{}, error)
	ListContentDomains(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]db.ListContentDomainsRow, error)
	ListContents(ctx context.Context, db db.DBTX, arg db.ListContentsParams) ([]db.ListContentsRow, error)
	ListTagsByUser(ctx context.Context, db db.DBTX, userID uuid.UUID) ([]db.ContentTag, error)
	OwnerTransferContent(ctx context.Context, db db.DBTX, arg db.OwnerTransferContentParams) error
	UpdateContent(ctx context.Context, db db.DBTX, arg db.UpdateContentParams) (db.Content, error)
}
