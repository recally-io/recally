package bots

import (
	"context"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"

	"github.com/labstack/echo/v4"
)

type Service struct {
	*Bot
}

func NewServer(botType BotType, cfg config.TelegramConfig, pool *db.Pool, e *echo.Echo, dbCache *cache.DbCache) (*Service, error) {
	var b *Bot
	var err error
	if botType == ReaderBot {
		b, err = NewReaderBot(cfg, pool, e, dbCache)
	} else if botType == ChatBot {
		b, err = NewChatBot(cfg, pool, e, dbCache)
	}
	if err != nil {
		return nil, err
	}
	return &Service{
		Bot: b,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	s.Bot.Start(ctx)
}

func (s *Service) Stop(ctx context.Context) {
	s.Bot.Stop(ctx)
}

func (s *Service) Name() string {
	return s.Bot.cfg.Name
}
