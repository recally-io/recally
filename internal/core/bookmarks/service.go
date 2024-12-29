package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	dao       DAO
	llm       LLM
	reader    UrlReader
	summarier Summarier
}

func NewService(llm *llms.LLM) *Service {
	reader, err := NewWebReader(llm)
	if err != nil {
		logger.Default.Fatal("failed to create web reader", "error", err)
	}
	summarier := NewSummarier(llm)
	return &Service{
		dao:       db.New(),
		llm:       llm,
		reader:    reader,
		summarier: summarier,
	}
}

// CreateBookmark creates a new bookmark with content fetching and embedding generation
func (s *Service) Create(ctx context.Context, tx db.DBTX, dto *BookmarkDTO) (*BookmarkDTO, error) {
	// Validate URL
	if _, err := url.ParseRequestURI(dto.URL); err != nil {
		return nil, fmt.Errorf("%w: invalid URL", ErrInvalidInput)
	}

	// Check for existing bookmark
	existing, err := s.dao.GetBookmarkByURL(ctx, tx, db.GetBookmarkByURLParams{
		Url:    dto.URL,
		UserID: pgtype.UUID{Bytes: dto.UserID, Valid: true},
	})
	if err == nil {
		return nil, fmt.Errorf("%w, id: %s", ErrDuplicate, existing.Uuid)
	}

	if !db.IsNotFoundError(err) {
		return nil, fmt.Errorf("failed to check existing bookmark for url '%s': %w", dto.URL, err)
	}

	bookmark, err := s.dao.CreateBookmark(ctx, tx, dto.Dump())
	if err != nil {
		return nil, fmt.Errorf("failed to create bookmark for url '%s': %w", dto.URL, err)
	}
	dto.Load(&bookmark)
	return dto, nil
}

// GetBookmark retrieves a bookmark by ID
func (s *Service) Get(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if bookmark.UserID.Bytes != userID {
		return nil, ErrUnauthorized
	}

	var dto BookmarkDTO
	dto.Load(&bookmark)
	return &dto, nil
}

// ListBookmarks retrieves a paginated list of bookmarks for a user
func (s *Service) List(ctx context.Context, tx db.DBTX, userID uuid.UUID, limit, offset int32) ([]*BookmarkDTO, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	bookmarks, err := s.dao.ListBookmarks(ctx, tx, db.ListBookmarksParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})

	dtos := make([]*BookmarkDTO, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		var dto BookmarkDTO
		dto.Load(&bookmark)
		dtos = append(dtos, &dto)
	}
	return dtos, err
}

// UpdateBookmark updates an existing bookmark
func (s *Service) Update(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, dto *BookmarkDTO) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if bookmark.UserID.Bytes != userID {
		return nil, ErrUnauthorized
	}

	updateParams := dto.DumpToUpdateParams()
	bookmark, err = s.dao.UpdateBookmark(ctx, tx, updateParams)
	if err != nil {
		return nil, err
	}

	dto.Load(&bookmark)
	return dto, nil
}

// DeleteBookmark removes a bookmark
func (s *Service) Delete(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) error {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return err
	}

	if bookmark.UserID.Bytes != userID {
		return ErrUnauthorized
	}

	return s.dao.DeleteBookmark(ctx, tx, db.DeleteBookmarkParams{
		Uuid:   id,
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
	})
}

// DeleteUserBookmarks removes all bookmarks for a user
func (s *Service) DeleteUserBookmarks(ctx context.Context, tx db.DBTX, userID uuid.UUID) error {
	return s.dao.DeleteBookmarksByUser(ctx, tx, pgtype.UUID{Bytes: userID, Valid: true})
}

func (s *Service) Refresh(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType FecherType, regenerateSummary bool) (*BookmarkDTO, error) {
	var dto *BookmarkDTO
	var err error

	if fetcherType != FecherType("") {
		dto, err = s.FetchContent(ctx, tx, id, userID, fetcherType)
		if err != nil {
			return nil, err
		}
	}

	if regenerateSummary {
		dto, err = s.SummarierContent(ctx, tx, id, userID)
		if err != nil {
			return nil, err
		}
	}

	return dto, nil
}

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType FecherType) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if bookmark.UserID.Bytes != userID {
		return nil, ErrUnauthorized
	}

	var dto BookmarkDTO
	dto.Load(&bookmark)

	var readerFetcher webreader.Fetcher
	switch fetcherType {
	case HttpFetcher:
		readerFetcher, err = fetcher.NewHTTPFetcher()
	case JinaFetcher:
		readerFetcher, err = fetcher.NewJinaFetcher()
	case BrowserFetcher:
		readerFetcher, err = fetcher.NewBrowserFetcher()
	default:
		err = fmt.Errorf("invalid fetcher type: %s", fetcherType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create fetcher: %w", err)
	}

	reader := webreader.New(readerFetcher, processor.NewMarkdownProcessor())
	content, err := reader.Read(ctx, bookmark.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	dto.Content = content.Markwdown
	dto.Title = content.Title
	dto.HTML = content.Html
	dto.Metadata.Image = content.Image
	dto.Metadata.Description = content.Description
	if content.PublishedTime != nil {
		dto.Metadata.PublishedAt = *content.PublishedTime
	}
	return s.Update(ctx, tx, id, userID, &dto)
}

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if bookmark.UserID.Bytes != userID {
		return nil, ErrUnauthorized
	}

	var dto BookmarkDTO
	dto.Load(&bookmark)

	content := &webreader.Content{
		Markwdown: dto.Content,
	}

	if err := s.summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		dto.Summary = content.Summary
	}

	return s.Update(ctx, tx, id, userID, &dto)
}
