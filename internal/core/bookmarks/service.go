package bookmarks

import (
	"context"
	"fmt"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Service struct {
	dao *db.Queries
	llm *llms.LLM
}

func NewService(llm *llms.LLM) *Service {
	return &Service{
		dao: db.New(),
		llm: llm,
	}
}

type TagDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (s *Service) ListTags(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]TagDTO, error) {
	tags, err := s.dao.ListBookmarkTagsByUser(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags for user '%s': %w", userID.String(), err)
	}

	tagsList := make([]TagDTO, 0, len(tags))
	for _, tag := range tags {
		tagsList = append(tagsList, TagDTO{
			Name:  tag.Name,
			Count: tag.Cnt,
		})
	}
	return tagsList, nil
}

type DomainDTO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (s *Service) ListDomains(ctx context.Context, tx db.DBTX, userID uuid.UUID) ([]DomainDTO, error) {
	domains, err := s.dao.ListBookmarkDomains(ctx, tx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list domains for user '%s': %w", userID.String(), err)
	}
	tags := make([]DomainDTO, 0, len(domains))
	for _, domain := range domains {
		tags = append(tags, DomainDTO{
			Name:  domain.Domain.String,
			Count: domain.Cnt,
		})
	}
	return tags, nil
}
