package bookmarks

import (
	"encoding/json"
	"net/url"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type Highlight struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
	Note        string `json:"note,omitempty"`
}

type Metadata struct {
	Author      string    `json:"author,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	Description string    `json:"description,omitempty"`
	SiteName    string    `json:"site_name,omitempty"`
	Domain      string    `json:"domain,omitempty"`

	Image   string `json:"image,omitempty"`
	Favicon string `json:"favicon"`
	Cover   string `json:"cover,omitempty"`

	Tags       []string    `json:"tags,omitempty"`
	Highlights []Highlight `json:"highlights,omitempty"`

	Share *ContentShareDTO `json:"share,omitempty"`
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

type ContentType string

const (
	ContentTypeBookmark   ContentType = "bookmark"
	ContentTypePDF        ContentType = "pdf"
	ContentTypeEPUB       ContentType = "epub"
	ContentTypeRSS        ContentType = "rss"
	ContentTypeNewsletter ContentType = "newsletter"
	ContentTypeImage      ContentType = "image"
	ContentTypePodcast    ContentType = "podcast"
	ContentTypeVideo      ContentType = "video"
)

type ContentDTO struct {
	ID          uuid.UUID   `json:"id"`
	UserID      uuid.UUID   `json:"user_id"`
	Type        ContentType `json:"type"`
	URL         string      `json:"url,omitempty"`
	Domain      string      `json:"domain,omitempty"`
	S3Key       string      `json:"s3_key,omitempty"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Content     string      `json:"content,omitempty"`
	HTML        string      `json:"html,omitempty"`
	Metadata    Metadata    `json:"metadata,omitempty"`
	IsFavorite  bool        `json:"is_favorite,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func loadTag(dbTags interface{}) []string {
	if dbTags != nil {
		tags := make([]string, 0, len(dbTags.([]interface{})))
		for _, tag := range dbTags.([]interface{}) {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			}
		}
		return tags
	}
	return nil
}

func (c *ContentDTO) Load(dbo *db.Content) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.Type = ContentType(dbo.Type)
	c.URL = dbo.Url.String
	c.Domain = dbo.Domain.String
	c.S3Key = dbo.S3Key.String
	c.Title = dbo.Title
	c.Description = dbo.Description.String
	c.Summary = dbo.Summary.String
	c.Content = dbo.Content.String
	c.HTML = dbo.Html.String
	c.IsFavorite = dbo.IsFavorite.Bool
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &c.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Content metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

func (c *ContentDTO) LoadWithTags(dbo *db.GetContentRow) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.Type = ContentType(dbo.Type)
	c.URL = dbo.Url.String
	c.Domain = dbo.Domain.String
	c.S3Key = dbo.S3Key.String
	c.Title = dbo.Title
	c.Description = dbo.Description.String
	c.Summary = dbo.Summary.String
	c.Content = dbo.Content.String
	c.HTML = dbo.Html.String
	c.IsFavorite = dbo.IsFavorite.Bool
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time
	c.Tags = loadTag(dbo.Tags)

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &c.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Content metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

func (c *ContentDTO) LoadWithTagsAndTotalCount(dbo *db.ListContentsRow) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.Type = ContentType(dbo.Type)
	c.URL = dbo.Url.String
	c.Domain = dbo.Domain.String
	c.S3Key = dbo.S3Key.String
	c.Title = dbo.Title
	c.Description = dbo.Description.String
	c.Summary = dbo.Summary.String
	c.Content = dbo.Content.String
	c.HTML = dbo.Html.String
	c.IsFavorite = dbo.IsFavorite.Bool
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time
	c.Tags = loadTag(dbo.Tags)

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &c.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Content metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

func (c *ContentDTO) LoadWithTagsAndTotalCountFromSearch(dbo *db.SearchContentsWithFilterRow) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.Type = ContentType(dbo.Type)
	c.URL = dbo.Url.String
	c.Domain = dbo.Domain.String
	c.S3Key = dbo.S3Key.String
	c.Title = dbo.Title
	c.Description = dbo.Description.String
	c.Summary = dbo.Summary.String
	c.Content = dbo.Content.String
	c.HTML = dbo.Html.String
	c.IsFavorite = dbo.IsFavorite.Bool
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time
	c.Tags = loadTag(dbo.Tags)

	if dbo.Metadata != nil {
		if err := json.Unmarshal(dbo.Metadata, &c.Metadata); err != nil {
			logger.Default.Warn("failed to unmarshal Content metadata",
				"err", err, "metadata", string(dbo.Metadata))
		}
	}
}

func (c *ContentDTO) Dump() db.CreateContentParams {
	metadata, _ := json.Marshal(c.Metadata)

	if (c.Domain == "") && (c.URL != "") {
		u, _ := url.Parse(c.URL)
		c.Domain = u.Host
	}

	return db.CreateContentParams{
		UserID:      c.UserID,
		Type:        string(c.Type),
		Title:       c.Title,
		Description: pgtype.Text{String: c.Description, Valid: c.Description != ""},
		Url:         pgtype.Text{String: c.URL, Valid: c.URL != ""},
		Domain:      pgtype.Text{String: c.Domain, Valid: c.Domain != ""},
		S3Key:       pgtype.Text{String: c.S3Key, Valid: c.S3Key != ""},
		Summary:     pgtype.Text{String: c.Summary, Valid: c.Summary != ""},
		Content:     pgtype.Text{String: c.Content, Valid: c.Content != ""},
		Html:        pgtype.Text{String: c.HTML, Valid: c.HTML != ""},
		IsFavorite:  pgtype.Bool{Bool: c.IsFavorite, Valid: true},
		Metadata:    metadata,
	}
}

func (c *ContentDTO) DumpToUpdateParams() db.UpdateContentParams {
	metadata, _ := json.Marshal(c.Metadata)
	return db.UpdateContentParams{
		ID:          c.ID,
		UserID:      c.UserID,
		Title:       pgtype.Text{String: c.Title, Valid: c.Title != ""},
		Description: pgtype.Text{String: c.Description, Valid: c.Description != ""},
		Url:         pgtype.Text{String: c.URL, Valid: c.URL != ""},
		Domain:      pgtype.Text{String: c.Domain, Valid: c.Domain != ""},
		S3Key:       pgtype.Text{String: c.S3Key, Valid: c.S3Key != ""},
		Summary:     pgtype.Text{String: c.Summary, Valid: c.Summary != ""},
		Content:     pgtype.Text{String: c.Content, Valid: c.Content != ""},
		Html:        pgtype.Text{String: c.HTML, Valid: c.HTML != ""},
		IsFavorite:  pgtype.Bool{Bool: c.IsFavorite, Valid: true},
		Metadata:    metadata,
	}
}

// Load converts a database object to a domain object
func (b *BookmarkDTO) Load(dbo *db.Bookmark) {
	b.ID = dbo.Uuid
	b.UserID = dbo.UserID.Bytes
	b.URL = dbo.Url
	b.Title = dbo.Title.String
	b.Summary = dbo.Summary.String
	if dbo.SummaryEmbeddings != nil {
		b.SummaryEmbedding = dbo.SummaryEmbeddings.Slice()
	}
	b.Content = dbo.Content.String
	if dbo.ContentEmbeddings != nil {
		b.ContentEmbedding = dbo.ContentEmbeddings.Slice()
	}
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

// Load converts a database object to a domain object
func (b *BookmarkDTO) LoadWithCount(dbo *db.ListBookmarksRow) {
	b.ID = dbo.Uuid
	b.UserID = dbo.UserID.Bytes
	b.URL = dbo.Url
	b.Title = dbo.Title.String
	b.Summary = dbo.Summary.String
	if dbo.SummaryEmbeddings != nil {
		b.SummaryEmbedding = dbo.SummaryEmbeddings.Slice()
	}
	b.Content = dbo.Content.String
	if dbo.ContentEmbeddings != nil {
		b.ContentEmbedding = dbo.ContentEmbeddings.Slice()
	}
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
	if len(b.ContentEmbedding) > 0 {
		v := pgvector.NewVector(b.ContentEmbedding)
		bookmark.ContentEmbeddings = &v
	}
	if len(b.SummaryEmbedding) > 0 {
		v := pgvector.NewVector(b.SummaryEmbedding)
		bookmark.SummaryEmbeddings = &v
	}
	return bookmark
}

// Dump to UpdateBookmarkParams
func (b *BookmarkDTO) DumpToUpdateParams() db.UpdateBookmarkParams {
	metadata, _ := json.Marshal(b.Metadata)
	p := db.UpdateBookmarkParams{
		Uuid:       b.ID,
		UserID:     pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		Title:      pgtype.Text{String: b.Title, Valid: b.Title != ""},
		Summary:    pgtype.Text{String: b.Summary, Valid: b.Summary != ""},
		Content:    pgtype.Text{String: b.Content, Valid: b.Content != ""},
		Html:       pgtype.Text{String: b.HTML, Valid: b.HTML != ""},
		Screenshot: pgtype.Text{String: b.Screenshot, Valid: b.Screenshot != ""},
		Metadata:   metadata,
	}

	if len(b.ContentEmbedding) > 0 {
		v := pgvector.NewVector(b.ContentEmbedding)
		p.ContentEmbeddings = &v
	}
	if len(b.SummaryEmbedding) > 0 {
		v := pgvector.NewVector(b.SummaryEmbedding)
		p.SummaryEmbeddings = &v
	}

	return p
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

type ContentShareDTO struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	ContentID uuid.UUID `json:"content_id"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *ContentShareDTO) Load(dbo *db.ContentShare) {
	c.ID = dbo.ID
	c.UserID = dbo.UserID
	c.ContentID = dbo.ContentID.Bytes
	c.ExpiresAt = dbo.ExpiresAt.Time
	c.CreatedAt = dbo.CreatedAt.Time
	c.UpdatedAt = dbo.UpdatedAt.Time
}

func (c *ContentShareDTO) Dump() db.CreateShareContentParams {
	return db.CreateShareContentParams{
		UserID:    c.UserID,
		ContentID: pgtype.UUID{Bytes: c.ContentID, Valid: c.ContentID != uuid.Nil},
		ExpiresAt: pgtype.Timestamptz{
			Time:  c.ExpiresAt,
			Valid: !c.ExpiresAt.IsZero(),
		},
	}
}
