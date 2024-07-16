package bots

import (
	"context"

	tele "gopkg.in/telebot.v3"
)


type Service struct {
	b *tele.Bot
}


func NewServer(token string, handlers ...Handler) (*Service, error) {
	b, err := New(token, handlers...)
	if err != nil {
		return nil, err
	}
	return &Service{
		b: b,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	s.b.Start()
}

func (s *Service) Stop(ctx context.Context) {
	s.b.Stop()
}

func (s *Service) Name() string {
	return "telegram bot"
}
