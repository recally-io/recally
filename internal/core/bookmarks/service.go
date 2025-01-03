package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader"
	"recally/internal/pkg/webreader/fetcher"
	"recally/internal/pkg/webreader/processor"
	"recally/internal/pkg/webreader/reader"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	dao DAO
	llm *llms.LLM
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		dao: db.New(),
		llm: llm,
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
	// Clear content and HTML
	dto.HTML = ""
	dto.SummaryEmbedding = nil
	dto.ContentEmbedding = nil
	return &dto, nil
}

// ListBookmarks retrieves a paginated list of bookmarks for a user
func (s *Service) List(ctx context.Context, tx db.DBTX, userID uuid.UUID, limit, offset int32) ([]*BookmarkDTO, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 50 // Default limit
	}
	if offset < 0 {
		offset = 0
	}

	totalCount := int64(0)
	bookmarks, err := s.dao.ListBookmarks(ctx, tx, db.ListBookmarksParams{
		UserID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:  limit,
		Offset: offset,
	})

	dtos := make([]*BookmarkDTO, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		var dto BookmarkDTO
		dto.LoadWithCount(&bookmark)
		dto.HTML = ""
		dto.SummaryEmbedding = nil
		dto.ContentEmbedding = nil
		dtos = append(dtos, &dto)
		totalCount = bookmark.TotalCount
	}
	return dtos, totalCount, err
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

func (s *Service) Refresh(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType fetcher.FecherType, regenerateSummary bool) (*BookmarkDTO, error) {
	var dto *BookmarkDTO
	var err error

	if fetcherType != fetcher.TypeNil {
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

func (s *Service) FetchContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID, fetcherType fetcher.FecherType) (*BookmarkDTO, error) {
	bookmark, err := s.dao.GetBookmarkByUUID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bookmark by id '%s': %w", id.String(), err)
	}

	if bookmark.UserID.Bytes != userID {
		return nil, ErrUnauthorized
	}

	var dto BookmarkDTO
	dto.Load(&bookmark)

	reader, err := reader.New(fetcherType, bookmark.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	content, err := reader.Read(ctx, bookmark.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content: %w", err)
	}

	dto.Content = content.Markwdown
	dto.Title = content.Title
	dto.HTML = content.Html

	// Update metadata
	dto.Metadata.Author = content.Author
	dto.Metadata.SiteName = content.SiteName
	dto.Metadata.Description = content.Description

	dto.Metadata.Cover = content.Cover
	dto.Metadata.Favicon = content.Favicon
	if content.Cover != "" {
		dto.Metadata.Image = content.Cover
	} else {
		dto.Metadata.Image = content.Favicon
	}

	if content.PublishedTime != nil {
		dto.Metadata.PublishedAt = *content.PublishedTime
	}
	return s.Update(ctx, tx, id, userID, &dto)
}

func (s *Service) SummarierContent(ctx context.Context, tx db.DBTX, id, userID uuid.UUID) (*BookmarkDTO, error) {
	user, err := auth.LoadUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

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

	summarier := newSummarier(s.llm, user)

	if err := summarier.Process(ctx, content); err != nil {
		logger.Default.Error("failed to generate summary", "err", err)
	} else {
		dto.Summary = content.Summary
	}

	return s.Update(ctx, tx, id, userID, &dto)
}

func newSummarier(llm *llms.LLM, user *auth.UserDTO) *processor.SummaryProcessor {
	summaryOptions := make([]processor.SummaryOption, 0)
	if user.Settings.SummaryOptions.Prompt != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionPrompt(user.Settings.SummaryOptions.Prompt))
	}
	if user.Settings.SummaryOptions.Model != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionModel(user.Settings.SummaryOptions.Model))
	}
	if user.Settings.SummaryOptions.Language != "" {
		summaryOptions = append(summaryOptions, processor.WithSummaryOptionLanguage(user.Settings.SummaryOptions.Language))
	}

	return processor.NewSummaryProcessor(llm, summaryOptions...)
}
