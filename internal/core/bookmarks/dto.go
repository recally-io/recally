package bookmarks

import (
	"encoding/json"
	"net/url"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
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

func (c *ContentDTO) FromReaderContent(article *webreader.Content) {
	c.Content = article.Markwdown
	c.Title = article.Title
	c.HTML = article.Html

	// Update metadata
	c.Metadata.Author = article.Author
	c.Metadata.SiteName = article.SiteName
	c.Metadata.Description = article.Description

	c.Metadata.Cover = article.Cover
	c.Metadata.Favicon = article.Favicon
	if article.Cover != "" {
		c.Metadata.Image = article.Cover
	} else {
		c.Metadata.Image = article.Favicon
	}

	if article.PublishedTime != nil {
		c.Metadata.PublishedAt = *article.PublishedTime
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
