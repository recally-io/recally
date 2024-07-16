package bots

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Handler func(b *tele.Bot)

func WithHandler(endpoint string, handler tele.HandlerFunc) Handler {
	return func(b *tele.Bot) {
		b.Handle(endpoint, handler)
	}
}

func New(token string, handlers ...Handler) (*tele.Bot, error) {
	b, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create new bot: %w", err)
	}

	if len(handlers) == 0 {
		return nil, fmt.Errorf("no handlers provided")
	}
	registerMiddlewarw(b)

	for _, handler := range handlers {
		handler(b)
	}
	return b, nil
}

func DefaultHandlers() []Handler {
	return []Handler{
		WithHandler("/start", func(c tele.Context) error {
			ctx := c.Get(constant.ContextKeyContext).(context.Context)
			logger.FromContext(ctx).Info("start command")
			return c.Send("Hello! I'm a bot. Ask me anything.")
		}),
	}
}

func registerMiddlewarw(b *tele.Bot) {
	b.Use(contextMiddleware())
	b.Use(middleware.Recover())
}

func contextMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyRequestID), uuid.Must(uuid.NewV7()).String())
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyUserID), c.Sender().ID)
			logger := logger.FromContext(ctx, slog.String("user_name", c.Sender().Username))
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyLogger), logger)
			c.Set(constant.ContextKeyContext, ctx)
			return next(c)
		}
	}
}
