package bookmarks

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type Service struct {
	dao      DAO
	embedder Embedder
	fetcher  URLFetcher
}

func NewService(embedder Embedder) *Service {
	return &Service{
		dao:      db.New(),
		embedder: embedder,
		fetcher:  NewJinaFetcher(),
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
		dto.Load(&existing)
		return dto, nil
	}

	if dto.Content == "" {
		readerResult, err := s.fetcher.Fetch(ctx, dto.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch content: %w", err)
		}
		dto.Content = readerResult.Content
		dto.Title = readerResult.Title
	}

	bookmark, err := s.dao.CreateBookmark(ctx, tx, dto.Dump())
	if err != nil {
		return nil, err
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

	metadata, _ := json.Marshal(dto.Metadata)
	updateParams := db.UpdateBookmarkParams{
		Uuid:       id,
		UserID:     pgtype.UUID{Bytes: userID, Valid: true},
		Title:      pgtype.Text{String: dto.Title, Valid: dto.Title != ""},
		Summary:    pgtype.Text{String: dto.Summary, Valid: dto.Summary != ""},
		Content:    pgtype.Text{String: dto.Content, Valid: dto.Content != ""},
		Html:       pgtype.Text{String: dto.HTML, Valid: dto.HTML != ""},
		Screenshot: pgtype.Text{String: dto.Screenshot, Valid: dto.Screenshot != ""},
		Metadata:   metadata,
	}

	if len(dto.ContentEmbedding) > 0 {
		updateParams.ContentEmbeddings = pgvector.NewVector(dto.ContentEmbedding)
	}
	if len(dto.SummaryEmbedding) > 0 {
		updateParams.SummaryEmbeddings = pgvector.NewVector(dto.SummaryEmbedding)
	}

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
