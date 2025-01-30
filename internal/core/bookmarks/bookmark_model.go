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
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	ContentID  uuid.UUID          `json:"content_id"`
	IsFavorite bool               `json:"is_favorite"`
	IsArchive  bool               `json:"is_archive"`
	IsPublic   bool               `json:"is_public"`
	Metadata   BookmarkMetadata   `json:"metadata"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	Tags       []string           `json:"tags"`
	Content    BookmarkContentDTO `json:"content"`
}

func (b *BookmarkDTO) Load(dbo *db.Bookmark) {
	b.ID = dbo.ID
	b.UserID = dbo.UserID.Bytes
	b.ContentID = dbo.ContentID.Bytes
	b.IsFavorite = dbo.IsFavorite
	b.IsArchive = dbo.IsArchive
	b.IsPublic = dbo.IsPublic
	b.CreatedAt = dbo.CreatedAt.Time
	b.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		b.Metadata = loadBookmarkMetadata(dbo.Metadata)
	}
}

func (b *BookmarkDTO) LoadWithContent(dbo *db.GetBookmarkWithContentRow) {
	// Load bookmark data
	b.Load(&dbo.Bookmark)
	// load bookmark content
	var content BookmarkContentDTO
	content.Load(&dbo.BookmarkContent)
	b.Content = content
	// Load tags from the aggregated tags field
	b.Tags = loadBookmarkTags(dbo.Tags)
}

func (b *BookmarkDTO) Dump() db.CreateBookmarkParams {
	return db.CreateBookmarkParams{
		UserID:     pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		ContentID:  pgtype.UUID{Bytes: b.ContentID, Valid: b.ContentID != uuid.Nil},
		IsFavorite: b.IsFavorite,
		IsArchive:  b.IsArchive,
		IsPublic:   b.IsPublic,
		Metadata:   dumpBookmarkMetadata(b.Metadata),
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
		Metadata: dumpBookmarkMetadata(b.Metadata),
	}
}

func loadListBookmarks(dbos []db.ListBookmarksRow) []BookmarkDTO {
	bookmarks := make([]BookmarkDTO, len(dbos))
	for i, dbo := range dbos {
		b := &bookmarks[i]
		b.Load(&dbo.Bookmark)
		// load bookmaek content
		var content BookmarkContentDTO
		content.Load(&dbo.BookmarkContent)
		b.Content = content
		// Load tags
		b.Tags = loadBookmarkTags(dbo.Tags)
	}
	return bookmarks
}

func loadSearchBookmarks(dbos []db.SearchBookmarksRow) []BookmarkDTO {
	bookmarks := make([]BookmarkDTO, len(dbos))
	for i, dbo := range dbos {
		b := &bookmarks[i]
		b.Load(&dbo.Bookmark)
		// load bookmaek content
		var content BookmarkContentDTO
		content.Load(&dbo.BookmarkContent)
		b.Content = content
		// Load tags
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
