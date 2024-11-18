package bookmarks

import (
	"context"
	"fmt"
	"net/url"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	dao    DAO
	llm    LLM
	reader UrlReader
}

func NewService(llm *llms.LLM) *Service {
	reader, err := NewWebReader(llm)
	if err != nil {
		logger.Default.Fatal("failed to create web reader", "error", err)
	}
	return &Service{
		dao:    db.New(),
		llm:    llm,
		reader: reader,
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

	if dto.Content == "" {
		readerResult, err := s.reader.Read(ctx, dto.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch content for url '%s': %w", dto.URL, err)
		}
		dto.Content = readerResult.Markwdown
		dto.Title = readerResult.Title
		dto.Summary = readerResult.Summary
		dto.HTML = readerResult.Html
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
