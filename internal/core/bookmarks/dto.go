package bookmarks

import (
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Metadata struct {
	Author      string    `json:"author,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
}

// BookmarkDTO represents the domain model for a bookmark
type BookmarkDTO struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	URL              string    `json:"url"`
	Title            string    `json:"title"`
	Summary          string    `json:"summary,omitempty"`
	SummaryEmbedding []float32 `json:"summary_embedding,omitempty"`
	Content          string    `json:"content,omitempty"`
	ContentEmbedding []float32 `json:"content_embedding,omitempty"`
	HTML             string    `json:"html,omitempty"`
	Metadata         Metadata  `json:"metadata,omitempty"`
	Screenshot       string    `json:"screenshot,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Load converts a database object to a domain object
func (b *BookmarkDTO) Load(dbo *db.Bookmark) {
	b.ID = dbo.Uuid
	b.UserID = dbo.UserID.Bytes
	b.URL = dbo.Url
	b.Title = dbo.Title.String
	b.Summary = dbo.Summary.String
	b.SummaryEmbedding = dbo.ContentEmbeddings.Slice()
	b.Content = dbo.Content.String
	b.ContentEmbedding = dbo.ContentEmbeddings.Slice()
	b.HTML = dbo.Html.String
	b.Screenshot = dbo.Screenshot.String
	b.CreatedAt = dbo.CreatedAt.Time
	b.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &b.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Bookmark metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

// Dump converts a domain object to a database object for creation
func (b *BookmarkDTO) Dump() db.CreateBookmarkParams {
	metadata, _ := json.Marshal(b.Metadata)
	id := b.ID
	if id == uuid.Nil {
		id = uuid.New()
	}
	bookmark := db.CreateBookmarkParams{
		Uuid:       id,
		UserID:     pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		Url:        b.URL,
		Title:      pgtype.Text{String: b.Title, Valid: b.Title != ""},
		Summary:    pgtype.Text{String: b.Summary, Valid: b.Summary != ""},
		Content:    pgtype.Text{String: b.Content, Valid: b.Content != ""},
		Html:       pgtype.Text{String: b.HTML, Valid: b.HTML != ""},
		Metadata:   metadata,
		Screenshot: pgtype.Text{String: b.Screenshot, Valid: b.Screenshot != ""},
	}
	return bookmark
}

// BookmarkOption defines a function type for configuring BookmarkDTO
type BookmarkOption func(*BookmarkDTO)

// NewBookmark creates a new BookmarkDTO with the given options
func NewBookmark(userID uuid.UUID, url string, opts ...BookmarkOption) *BookmarkDTO {
	b := &BookmarkDTO{
		ID:     uuid.New(),
		UserID: userID,
		URL:    url,
	}

	for _, opt := range opts {
		opt(b)
	}
	return b
}

// Option functions for configuring a new bookmark
func WithTitle(title string) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.Title = title
	}
}

func WithSummary(summary string) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.Summary = summary
	}
}

func WithContent(content string) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.Content = content
	}
}

func WithHTML(html string) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.HTML = html
	}
}

func WithMetadata(metadata Metadata) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.Metadata = metadata
	}
}

func WithScreenshot(screenshot string) BookmarkOption {
	return func(b *BookmarkDTO) {
		b.Screenshot = screenshot
	}
}
