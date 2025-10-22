package bookmarks

import (
	"encoding/json"
	"slices"
	"time"

	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"

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
	LastReadAt      time.Time `json:"last_read_at"`

	Highlights []Highlight `json:"highlights,omitempty"`
}

type BookmarkDTO struct {
	ID         uuid.UUID           `json:"id"`
	UserID     uuid.UUID           `json:"user_id"`
	ContentID  uuid.UUID           `json:"content_id"`
	IsFavorite bool                `json:"is_favorite"`
	IsArchive  bool                `json:"is_archive"`
	IsPublic   bool                `json:"is_public"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Tags       []string            `json:"tags"`
	Metadata   *BookmarkMetadata   `json:"metadata"`
	Content    *BookmarkContentDTO `json:"content"`
	Share      *BookmarkShareDTO   `json:"share,omitempty"`
}

func (b *BookmarkDTO) Load(dbo *db.Bookmark) {
	b.ID = dbo.ID
	b.UserID = dbo.UserID.Bytes
	b.ContentID = dbo.ContentID.Bytes
	b.IsFavorite = dbo.IsFavorite
	b.IsArchive = dbo.IsArchive
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
	b.Content = &content

	if dbo.BookmarkShare.ID != uuid.Nil {
		var share BookmarkShareDTO

		share.Load(&dbo.BookmarkShare)
		b.Share = &share
		b.IsPublic = true
	}

	// Load tags from the aggregated tags field
	b.Tags = loadBookmarkTags(dbo.Tags)
}

func (b *BookmarkDTO) Dump() db.CreateBookmarkParams {
	return db.CreateBookmarkParams{
		UserID:     pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		ContentID:  pgtype.UUID{Bytes: b.ContentID, Valid: b.ContentID != uuid.Nil},
		IsFavorite: b.IsFavorite,
		IsArchive:  b.IsArchive,
		Metadata:   dumpBookmarkMetadata(b.Metadata),
	}
}

func (b *BookmarkDTO) DumpToUpdateParams() db.UpdateBookmarkParams {
	return db.UpdateBookmarkParams{
		ID:     b.ID,
		UserID: pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		IsFavorite: pgtype.Bool{
			Bool:  b.IsFavorite,
			Valid: b.IsFavorite,
		},
		IsArchive: pgtype.Bool{
			Bool:  b.IsArchive,
			Valid: b.IsArchive,
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
		content.Content = ""
		content.Html = ""
		b.Content = &content
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
		content.Content = ""
		content.Html = ""
		b.Content = &content
		// Load tags
		b.Tags = loadBookmarkTags(dbo.Tags)
	}

	return bookmarks
}

func loadBookmarkTags(input any) []string {
	if input == nil {
		return nil
	}

	// Handle case where input is []interface{}
	if interfaceSlice, ok := input.([]any); ok {
		tags := make([]string, len(interfaceSlice))

		for i, v := range interfaceSlice {
			if str, ok := v.(string); ok {
				tags[i] = str
			}
		}

		slices.Sort(tags)

		return tags
	}

	// Handle case where input is already []string
	if tags, ok := input.([]string); ok {
		return tags
	}

	return nil
}

func loadBookmarkMetadata(input any) *BookmarkMetadata {
	if input == nil {
		return nil
	}

	metadata := BookmarkMetadata{}
	if err := json.Unmarshal(input.([]byte), &metadata); err != nil {
		logger.Default.Warn("failed to unmarshal Bookmark metadata", "err", err, "metadata", string(input.(string)))
	}

	return &metadata
}

func dumpBookmarkMetadata(input *BookmarkMetadata) []byte {
	if input == nil {
		return nil
	}

	metadata, _ := json.Marshal(input)

	return metadata
}

func loadBookmarkContentMetadata(input any) *BookmarkContentMetadata {
	if input == nil {
		return nil
	}

	metadata := BookmarkContentMetadata{}
	if err := json.Unmarshal(input.([]byte), &metadata); err != nil {
		logger.Default.Warn("failed to unmarshal BookmarkContent metadata", "err", err, "metadata", string(input.([]byte)))
	}

	return &metadata
}

func dumpBookmarkContentMetadata(input *BookmarkContentMetadata) []byte {
	if input == nil {
		return nil
	}

	metadata, _ := json.Marshal(input)

	return metadata
}
