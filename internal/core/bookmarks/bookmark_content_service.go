package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"recally/internal/pkg/webreader/reader"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) IsBookmarkContentExistByURL(ctx context.Context, tx db.DBTX, url string) (bool, error) {
	return s.dao.IsBookmarkContentExistByURL(ctx, tx, url)
}

func (s *Service) CreateBookmarkContent(ctx context.Context, tx db.DBTX, content *BookmarkContentDTO) (*BookmarkContentDTO, error) {
	params := content.Dump()
	dbo, err := s.dao.CreateBookmarkContent(ctx, tx, params)
	if err != nil {
		return nil, err
	}

	result := &BookmarkContentDTO{}
	result.Load(&dbo)
	return result, nil
}

func (s *Service) GetBookmarkContentByID(ctx context.Context, tx db.DBTX, id uuid.UUID) (*BookmarkContentDTO, error) {
	dbo, err := s.dao.GetBookmarkContentByID(ctx, tx, id)
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

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, opts fetcher.FetchOptions) (*BookmarkContentDTO, error) {
	dto, err := s.GetBookmarkContentByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark by id '%s': %w", id.String(), err)
	}
	if dto.Content != "" && !opts.Force {
		return dto, nil
	}
	content, err := s.FetchContentWithCache(ctx, dto.URL, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	dto.FromReaderContent(content)
	return s.UpdateBookmarkContent(ctx, tx, dto)
}

func (s *Service) FetchContentWithCache(ctx context.Context, uri string, opts fetcher.FetchOptions) (*webreader.Content, error) {
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

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*BookmarkContentDTO, error) {
	user, err := auth.LoadUser(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	dto, err := s.GetBookmarkContentByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	content := &webreader.Content{
		Markwdown: dto.Content,
	}

	summarier := processor.NewSummaryProcessor(s.llm, processor.WithSummaryOptionUser(user))

	if len(content.Markwdown) < 1000 {
		logger.FromContext(ctx).Info("content is too short to summarise")
		return dto, nil
	}

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		tags, summary := parseTagsFromSummary(content.Summary)
		if len(tags) > 0 {
			if err := s.linkContentTags(ctx, tx, dto.Tags, tags, id, userID); err != nil {
				return nil, err
			}
		}
		dto.Summary = summary
	}
	return s.UpdateBookmarkContent(ctx, tx, dto)
}
