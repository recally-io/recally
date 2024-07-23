package bots

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/constant"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"
	"vibrain/internal/port/bots/handlers"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Service struct {
	b          *tele.Bot
	handler    *handlers.Handler
	token      string
	webhookUrl string
}

type Option func(*Service)

func WithCache(c *cache.DbCache) Option {
	return func(s *Service) {
		s.handler.Cache = c
	}
}

func WithWebhook(e *echo.Echo, webhookUrl string) Option {
	return func(s *Service) {
		s.webhookUrl = webhookUrl
		u, err := url.Parse(webhookUrl)
		if err != nil {
			logger.Default.Fatal("failed to parse webhook url", "err", err)
		}
		e.POST(u.Path, func(c echo.Context) error {
			if s.token != "" && c.Request().Header.Get("X-Telegram-Bot-Api-Secret-Token") != s.token {
				logger.FromContext(c.Request().Context()).Error("invalid secret token in request")
				return c.String(http.StatusUnauthorized, "invalid secret token")
			}

			var update tele.Update
			if err := json.NewDecoder(c.Request().Body).Decode(&update); err != nil {
				logger.FromContext(c.Request().Context()).Error("cannot decode update", "err", err)
				return c.String(http.StatusBadRequest, fmt.Sprintf("cannot decode update: %s", err))
			}
			s.b.Updates <- update
			return nil
		})
	}
}

func NewServer(token string, pool *db.Pool, opts ...Option) (*Service, error) {
	handler := handlers.New(pool)

	b, err := newBot(token, handler)
	if err != nil {
		return nil, err
	}
	return &Service{
		token: token,
		b:     b,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	if s.webhookUrl != "" {
		// SetWebhook
		params := map[string]string{
			"url": s.webhookUrl,
			// "drop_pending_updates": "true",
			// "secret_token": 	   s.token,
		}
		if _, err := s.b.Raw("setWebhook", params); err != nil {
			logger.FromContext(ctx).Error("failed to set webhook", "err", err)
		}
	} else {
		s.b.Start()
	}
}

func (s *Service) Stop(ctx context.Context) {
	if s.webhookUrl != "" {
		dropPending := true
		// RemoveWebhook
		if _, err := s.b.Raw("deleteWebhook", map[string]bool{
			"drop_pending_updates": dropPending,
		}); err != nil {
			logger.FromContext(ctx).Error("failed to remove webhook", "err", err)
		}
	} else {
		s.b.Stop()
	}
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
