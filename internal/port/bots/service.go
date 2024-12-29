package bots

import (
	"context"
	"recally/internal/core/queue"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/config"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"

	"github.com/labstack/echo/v4"
)

type Service struct {
	*Bot
}

func NewServer(botType BotType, cfg config.TelegramConfig, pool *db.Pool, e *echo.Echo, cacheService cache.Cache, llm *llms.LLM, queue *queue.Queue) (*Service, error) {
	var b *Bot
	var err error
	if botType == ReaderBot {
		b, err = NewReaderBot(cfg, pool, e, cacheService, llm, queue)
	} else if botType == ChatBot {
		b, err = NewChatBot(cfg, pool, e, cacheService, llm, queue)
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
