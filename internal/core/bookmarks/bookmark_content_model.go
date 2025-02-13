package bookmarks

import (
	"context"
	"recally/internal/core/files"
	"recally/internal/pkg/db"
	"recally/internal/pkg/webreader"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ContentType string

const (
	ContentTypeBookmark   ContentType = "bookmark"
	ContentTypePDF        ContentType = "pdf"
	ContentTypeEPUB       ContentType = "epub"
	ContentTypeRSS        ContentType = "rss"
	ContentTypeNewsletter ContentType = "newsletter"
	ContentTypeImage      ContentType = "image"
	ContentTypeAudio      ContentType = "audio"
	ContentTypeVideo      ContentType = "video"
)

type BookmarkContentFileMetadata struct {
	Name      string `json:"name,omitempty"`
	Extension string `json:"extension,omitempty"`
	MIMEType  string `json:"mime_type,omitempty"`
	Size      int64  `json:"size,omitempty"`
	PageCount int    `json:"page_count,omitempty"`
}

type BookmarkContentMetadata struct {
	Author      string    `json:"author,omitempty"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	Description string    `json:"description,omitempty"`
	SiteName    string    `json:"site_name,omitempty"`
	Domain      string    `json:"domain,omitempty"`

	Favicon string `json:"favicon"`
	Cover   string `json:"cover,omitempty"`

	File BookmarkContentFileMetadata `json:"file,omitempty"`
}

type BookmarkContentDTO struct {
	ID          uuid.UUID                `json:"id"`
	Type        ContentType              `json:"type"`
	URL         string                   `json:"url"`
	UserID      uuid.UUID                `json:"user_id"`
	Title       string                   `json:"title"`
	Description string                   `json:"description"`
	Domain      string                   `json:"domain"`
	S3Key       string                   `json:"s3_key"`
	Summary     string                   `json:"summary"`
	Content     string                   `json:"content"`
	Html        string                   `json:"html"`
	Tags        []string                 `json:"tags"`
	Metadata    *BookmarkContentMetadata `json:"metadata"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

func (b *BookmarkContentDTO) Load(dbo *db.BookmarkContent) {
	b.ID = dbo.ID
	b.Type = ContentType(dbo.Type)
	b.URL = dbo.Url
	b.UserID = dbo.UserID.Bytes
	b.Title = dbo.Title.String
	b.Description = dbo.Description.String
	b.Domain = dbo.Domain.String
	b.S3Key = dbo.S3Key.String
	b.Summary = dbo.Summary.String
	b.Content = dbo.Content.String
	b.Html = dbo.Html.String
	b.Tags = dbo.Tags
	b.CreatedAt = dbo.CreatedAt.Time
	b.UpdatedAt = dbo.UpdatedAt.Time

	if dbo.Metadata != nil {
		b.Metadata = loadBookmarkContentMetadata(dbo.Metadata)
	}
}

func (b *BookmarkContentDTO) Dump() db.CreateBookmarkContentParams {
	metadata := dumpBookmarkContentMetadata(b.Metadata)

	return db.CreateBookmarkContentParams{
		Type:   string(b.Type),
		Url:    b.URL,
		UserID: pgtype.UUID{Bytes: b.UserID, Valid: b.UserID != uuid.Nil},
		Title: pgtype.Text{
			String: b.Title,
			Valid:  b.Title != "",
		},
		Description: pgtype.Text{
			String: b.Description,
			Valid:  b.Description != "",
		},
		Domain: pgtype.Text{
			String: b.Domain,
			Valid:  b.Domain != "",
		},
		S3Key: pgtype.Text{
			String: b.S3Key,
			Valid:  b.S3Key != "",
		},
		Summary: pgtype.Text{
			String: b.Summary,
			Valid:  b.Summary != "",
		},
		Content: pgtype.Text{
			String: b.Content,
			Valid:  b.Content != "",
		},
		Html: pgtype.Text{
			String: b.Html,
			Valid:  b.Html != "",
		},
		Tags:     b.Tags,
		Metadata: metadata,
	}
}

func (b *BookmarkContentDTO) DumpToUpdateParams() db.UpdateBookmarkContentParams {
	return db.UpdateBookmarkContentParams{
		ID: b.ID,
		Title: pgtype.Text{
			String: b.Title,
			Valid:  b.Title != "",
		},
		Description: pgtype.Text{
			String: b.Description,
			Valid:  b.Description != "",
		},
		S3Key: pgtype.Text{
			String: b.S3Key,
			Valid:  b.S3Key != "",
		},
		Summary: pgtype.Text{
			String: b.Summary,
			Valid:  b.Summary != "",
		},
		Content: pgtype.Text{
			String: b.Content,
			Valid:  b.Content != "",
		},
		Html: pgtype.Text{
			String: b.Html,
			Valid:  b.Html != "",
		},
		Tags:     b.Tags,
		Metadata: dumpBookmarkContentMetadata(b.Metadata),
	}
}

func (c *BookmarkContentDTO) FromReaderContent(article *webreader.Content) {
	c.URL = article.URL
	c.Content = article.Markwdown
	c.Title = article.Title
	c.Html = article.Html

	// Update metadata
	c.Metadata = &BookmarkContentMetadata{
		Author:      article.Author,
		SiteName:    article.SiteName,
		Description: article.Description,
		Cover:       article.Cover,
		Favicon:     article.Favicon,
	}

	if article.PublishedTime != nil {
		c.Metadata.PublishedAt = *article.PublishedTime
	}
}

func (c *BookmarkContentDTO) GetFilePublicURL(ctx context.Context) (string, error) {
	if c.S3Key == "" {
		return c.URL, nil
	}

	return files.DefaultService.GetPublicURL(ctx, c.S3Key)
}
