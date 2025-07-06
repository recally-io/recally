package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"recally/internal/core/files"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/reader"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) IsBookmarkContentExistByURL(ctx context.Context, tx db.DBTX, url string) (bool, error) {
	return s.dao.IsBookmarkContentExistByURL(ctx, tx, url)
}

func (s *Service) CreateBookmarkContent(ctx context.Context, tx db.DBTX, content *BookmarkContentDTO) (*BookmarkContentDTO, error) {
	// when user save image from url by recally-clipper, we need to upload it to s3 first
	if content.S3Key == "" && content.IsMediaType() {
		// upload image to s3
		file, err := files.DefaultService.CreateFileAndUploadToS3FromUrl(ctx, tx, content.UserID, true, "", content.URL)
		if err != nil {
			return nil, err
		}

		content.S3Key = file.S3Key
	}

	if content.URL == "" {
		content.URL = content.S3Key
	}

	params := content.Dump()

	dbo, err := s.dao.CreateBookmarkContent(ctx, tx, params)
	if err != nil {
		return nil, err
	}

	result := &BookmarkContentDTO{}
	result.Load(&dbo)

	return result, nil
}

func (s *Service) GetBookmarkContentByBookmarkID(ctx context.Context, tx db.DBTX, bookmarkID uuid.UUID) (*BookmarkContentDTO, error) {
	dbo, err := s.dao.GetBookmarkContentByBookmarkID(ctx, tx, bookmarkID)
	if err != nil {
		return nil, err
	}

	result := &BookmarkContentDTO{}
	result.Load(&dbo)

	return result, nil
}

func (s *Service) GetBookmarkContentByURL(ctx context.Context, tx db.DBTX, url string, userID uuid.UUID) (*BookmarkContentDTO, error) {
	dbo, err := s.dao.GetBookmarkContentByURL(ctx, tx, db.GetBookmarkContentByURLParams{
		Url:    url,
		UserID: pgtype.UUID{Bytes: userID, Valid: userID != uuid.Nil},
	})
	if err != nil {
		return nil, err
	}

	result := &BookmarkContentDTO{}
	result.Load(&dbo)

	return result, nil
}

func (s *Service) UpdateBookmarkContent(ctx context.Context, tx db.DBTX, content *BookmarkContentDTO) (*BookmarkContentDTO, error) {
	params := content.DumpToUpdateParams()

	dbo, err := s.dao.UpdateBookmarkContent(ctx, tx, params)
	if err != nil {
		return nil, err
	}

	result := &BookmarkContentDTO{}
	result.Load(&dbo)

	return result, nil
}

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, bookmarkID, userID uuid.UUID, opts fetcher.FetchOptions) (*BookmarkContentDTO, error) {
	bookmarkContent, err := s.GetBookmarkContentByBookmarkID(ctx, tx, bookmarkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark content by id '%s': %w", bookmarkID.String(), err)
	}

	// if content is not empty and force is false, then return
	// if content type is not bookmark, then return
	if (bookmarkContent.Content != "" && !opts.Force) || bookmarkContent.Type != ContentTypeBookmark {
		return bookmarkContent, nil
	}

	webContent, err := s.FetchWebContentWithCache(ctx, bookmarkContent.URL, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	bookmarkContent.FromReaderContent(webContent)

	return s.UpdateBookmarkContent(ctx, tx, bookmarkContent)
}

func (s *Service) FetchWebContentWithCache(ctx context.Context, uri string, opts fetcher.FetchOptions) (*webreader.Content, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid url '%s': %w", uri, err)
	}

	reader, err := reader.New(u.Host, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}

	content, err := cache.RunInCache(ctx, cache.DefaultDBCache, cache.NewCacheKey(fmt.Sprintf("WebReader-%s", opts.FecherType), uri), 24*time.Hour, func() (*webreader.Content, error) {
		return reader.Fetch(ctx, uri)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	return reader.Process(ctx, content)
}
