package bots

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/port/bots/handlers"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Service struct {
	b       *tele.Bot
	handler *handlers.Handler
}

type Option func(*Service)

func WithCache(c *cache.DbCache) Option {
	return func(s *Service) {
		s.handler.Cache = c
	}
}

func NewServer(token string, pool *db.Pool, opts ...handlers.Option) (*Service, error) {
	handler := handlers.New(pool, opts...)
	b, err := newBot(token, handler)
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

func newBot(token string, handler *handlers.Handler) (*tele.Bot, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}
	registerMiddlewarw(b)
	registerhandlers(b, handler)
	return b, nil
}

func registerhandlers(b *tele.Bot, handler *handlers.Handler) {
	b.Handle(tele.OnText, handler.TextHandler)
	b.Handle("/start", func(c tele.Context) error {
		ctx := c.Get(constant.ContextKeyContext).(context.Context)
		logger.FromContext(ctx).Info("start command")
		return c.Send("Hello! I'm a bot. Ask me anything.")
	})
}

func registerMiddlewarw(b *tele.Bot) {
	b.Use(contextMiddleware())
	b.Use(middleware.Recover())
}

func contextMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()
			start := time.Now()
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyRequestID), uuid.Must(uuid.NewV7()).String())
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyUserID), c.Sender().ID)
			logger := logger.FromContext(ctx, slog.String("user_name", c.Sender().Username))
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyLogger), logger)
			c.Set(constant.ContextKeyContext, ctx)
			defer func() {
				logger.Info("message processed", slog.Duration("duration", time.Since(start)))
			}()
			return next(c)
		}
	}
}
