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
	"github.com/jackc/pgx/v5"
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
	if bookmarkContent.Content != "" && !opts.Force {
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

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, bookmarkID, userID uuid.UUID) (*BookmarkContentDTO, error) {
	user, err := auth.LoadUser(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	bookmarkContent, err := s.GetBookmarkContentByBookmarkID(ctx, tx, bookmarkID)
	if err != nil {
		return nil, err
	}

	content := &webreader.Content{
		Markwdown: bookmarkContent.Content,
	}

	summarier := processor.NewSummaryProcessor(s.llm, processor.WithSummaryOptionUser(user))

	if len(content.Markwdown) < 1000 {
		logger.FromContext(ctx).Info("content is too short to summarise")
		return bookmarkContent, nil
	}

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		summary, tags := s.ProcessSummaryTags(ctx, tx, bookmarkID, userID, content.Summary)
		bookmarkContent.Summary = summary
		if len(tags) > 0 {
			bookmarkContent.Tags = tags
		}
	}
	return s.UpdateBookmarkContent(ctx, tx, bookmarkContent)
}

func (s *Service) ProcessSummaryTags(ctx context.Context, tx db.DBTX, bookmarkID, userID uuid.UUID, summary string) (string, []string) {
	tags, summary := parseTagsFromSummary(summary)
	if len(tags) > 0 {
		// link tags in background
		newUserCtx := auth.SetUserToContextByUserID(context.Background(), userID)
		go func() {
			if err := db.RunInTransaction(newUserCtx, db.DefaultPool.Pool, func(ctx context.Context, tx pgx.Tx) error {
				return s.linkContentTags(ctx, tx, tags, tags, bookmarkID, userID)
			}); err != nil {
				logger.Default.Error("failed to link content tags", "err", err, "bookmark_id", bookmarkID)
			}
		}()
	}
	return summary, tags
}
