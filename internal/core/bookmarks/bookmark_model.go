package bookmarks

import (
	"encoding/json"
	"recally/internal/pkg/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Highlight struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
	Note        string `json:"note,omitempty"`
}

type BookmarkMetadata struct {
	ReadingProgress int       `json:"reading_progress,omitempty"`
	LastReadAt      time.Time `json:"last_read_at,omitempty"`

	Highlights []Highlight       `json:"highlights,omitempty"`
	Share      *BookmarkShareDTO `json:"share,omitempty"`
}

type BookmarkDTO struct {
	ID              uuid.UUID          `json:"id"`
	UserID          uuid.UUID          `json:"user_id"`
	ContentID       uuid.UUID          `json:"content_id"`
	IsFavorite      bool               `json:"is_favorite"`
	IsArchive       bool               `json:"is_archive"`
	IsPublic        bool               `json:"is_public"`
	ReadingProgress int                `json:"reading_progress"`
	Metadata        BookmarkMetadata   `json:"metadata"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Tags            []string           `json:"tags"`
	Content         BookmarkContentDTO `json:"content"`
}

func (b *BookmarkDTO) Load(dbo *db.Bookmark) {
	b.ID = dbo.ID
	b.UserID = dbo.UserID.Bytes
	b.ContentID = dbo.ContentID.Bytes
	b.IsFavorite = dbo.IsFavorite.Bool
	b.IsArchive = dbo.IsArchive.Bool
	b.IsPublic = dbo.IsPublic.Bool
	b.ReadingProgress = int(dbo.ReadingProgress.Int32)
	b.CreatedAt = dbo.CreatedAt.Time
	b.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		b.Metadata = loadBookmarkMetadata(dbo.Metadata)
	}
}

func (b *BookmarkDTO) LoadWithContent(dbo *db.GetBookmarkWithContentRow) {
	// Load bookmark data
	b.ID = dbo.ID
	b.UserID = dbo.UserID.Bytes
	b.ContentID = dbo.ID_2
	b.IsFavorite = dbo.IsFavorite.Bool
	b.IsArchive = dbo.IsArchive.Bool
	b.IsPublic = dbo.IsPublic.Bool
	b.ReadingProgress = int(dbo.ReadingProgress.Int32)
	b.CreatedAt = dbo.CreatedAt.Time
	b.UpdatedAt = dbo.UpdatedAt.Time

	// Load bookmark metadata
	if dbo.Metadata != nil {
		b.Metadata = loadBookmarkMetadata(dbo.Metadata)
	}

	// Load tags from the aggregated tags field
	b.Tags = loadBookmarkTags(dbo.Tags)

	// Load content data
	b.Content.ID = dbo.ID_2
	b.Content.Type = ContentType(dbo.Type)
	b.Content.URL = dbo.Url
	b.Content.UserID = dbo.UserID_2.Bytes
	b.Content.Title = dbo.Title.String
	b.Content.Description = dbo.Description.String
	b.Content.Domain = dbo.Domain.String
	b.Content.S3Key = dbo.S3Key.String
	b.Content.Summary = dbo.Summary.String
	b.Content.Content = dbo.Content.String
	b.Content.Html = dbo.Html.String
	b.Content.Tags = dbo.Tags
	b.Content.CreatedAt = dbo.CreatedAt_2.Time
	b.Content.UpdatedAt = dbo.UpdatedAt_2.Time

	// Load content metadata
	if dbo.Metadata_2 != nil {
		b.Content.Metadata = loadBookmarkContentMetadata(dbo.Metadata_2)
	}
}

func (b *BookmarkDTO) Dump() db.CreateBookmarkParams {
	return db.CreateBookmarkParams{
		UserID:    pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		ContentID: pgtype.UUID{Bytes: b.ContentID, Valid: b.ContentID != uuid.Nil},
		IsFavorite: pgtype.Bool{
			Bool:  b.IsFavorite,
			Valid: true,
		},
		IsArchive: pgtype.Bool{
			Bool:  b.IsArchive,
			Valid: true,
		},
		IsPublic: pgtype.Bool{
			Bool:  b.IsPublic,
			Valid: true,
		},
		ReadingProgress: pgtype.Int4{
			Int32: int32(b.ReadingProgress),
			Valid: true,
		},
		Metadata: dumpBookmarkMetadata(b.Metadata),
	}
}

func (b *BookmarkDTO) DumpToUpdateParams() db.UpdateBookmarkParams {
	return db.UpdateBookmarkParams{
		ID:     b.ID,
		UserID: pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		IsFavorite: pgtype.Bool{
			Bool:  b.IsFavorite,
			Valid: true,
		},
		IsArchive: pgtype.Bool{
			Bool:  b.IsArchive,
			Valid: true,
		},
		IsPublic: pgtype.Bool{
			Bool:  b.IsPublic,
			Valid: true,
		},
		ReadingProgress: pgtype.Int4{
			Int32: int32(b.ReadingProgress),
			Valid: true,
		},
		Metadata: dumpBookmarkMetadata(b.Metadata),
	}
}

func loadListBookmarks(dbos []db.ListBookmarksRow) []BookmarkDTO {
	bookmarks := make([]BookmarkDTO, len(dbos))
	for i, dbo := range dbos {
		b := &bookmarks[i]
		b.ID = dbo.ID
		b.UserID = dbo.UserID.Bytes
		b.ContentID = dbo.ContentID.Bytes
		b.IsFavorite = dbo.IsFavorite.Bool
		b.IsArchive = dbo.IsArchive.Bool
		b.IsPublic = dbo.IsPublic.Bool
		b.ReadingProgress = int(dbo.ReadingProgress.Int32)
		b.CreatedAt = dbo.CreatedAt.Time
		b.UpdatedAt = dbo.UpdatedAt.Time

		if dbo.Metadata != nil {
			b.Metadata = loadBookmarkMetadata(dbo.Metadata)
		}

		b.Tags = loadBookmarkTags(dbo.Tags)
	}
	return bookmarks
}

func loadSearchBookmarks(dbos []db.SearchBookmarksRow) []BookmarkDTO {
	bookmarks := make([]BookmarkDTO, len(dbos))
	for i, dbo := range dbos {
		b := &bookmarks[i]
		b.ID = dbo.ID
		b.UserID = dbo.UserID.Bytes
		b.ContentID = dbo.ContentID.Bytes
		b.IsFavorite = dbo.IsFavorite.Bool
		b.IsArchive = dbo.IsArchive.Bool
		b.IsPublic = dbo.IsPublic.Bool
		b.ReadingProgress = int(dbo.ReadingProgress.Int32)
		b.CreatedAt = dbo.CreatedAt.Time
		b.UpdatedAt = dbo.UpdatedAt.Time

		if dbo.Metadata != nil {
			b.Metadata = loadBookmarkMetadata(dbo.Metadata)
		}

		b.Tags = loadBookmarkTags(dbo.Tags)
	}
	return bookmarks
}

func loadBookmarkTags(input interface{}) []string {
	if tags, ok := input.([]string); ok {
		return tags
	}
	return nil
}

func loadBookmarkMetadata(input interface{}) BookmarkMetadata {
	if metadata, ok := input.(BookmarkMetadata); ok {
		return metadata
	}
	return BookmarkMetadata{}
}

func dumpBookmarkMetadata(input BookmarkMetadata) []byte {
	metadata, _ := json.Marshal(input)
	return metadata
}

func loadBookmarkContentMetadata(input interface{}) BookmarkContentMetadata {
	if metadata, ok := input.(BookmarkContentMetadata); ok {
		return metadata
	}
	return BookmarkContentMetadata{}
}

func dumpBookmarkContentMetadata(input BookmarkContentMetadata) []byte {
	metadata, _ := json.Marshal(input)
	return metadata
}
