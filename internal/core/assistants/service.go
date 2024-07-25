package assistants

import (
	"context"
	"vibrain/internal/pkg/db"

	"github.com/google/uuid"
)

type Service struct {
	db Repository
}

func NewService(db *db.Pool) (*Service, error) {
	s := &Service{
		db: NewRepository(db),
	}

	return s, nil
}

func (s *Service) CreateAssistant(ctx context.Context, assistant *Assistant) error {
	return s.db.CreateAssistant(ctx, assistant)
}

func (s *Service) GetAssistant(ctx context.Context, id string) (*Assistant, error) {
	return s.db.GetAssistant(ctx, uuid.MustParse(id))
}

func (s *Service) CreateThread(ctx context.Context, thread *Thread) error {
	return s.db.CreateThread(ctx, thread)
}

func (s *Service) GetThread(ctx context.Context, id string) (*Thread, error) {
	return s.db.GetThread(ctx, uuid.MustParse(id))
}
