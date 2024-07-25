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
)

func contextMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()
			start := time.Now()
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyRequestID), uuid.Must(uuid.NewV7()).String())
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyUserID), fmt.Sprintf("%d", c.Sender().ID))
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyUserName), c.Sender().Username)
			logger := logger.FromContext(ctx)
			ctx = context.WithValue(ctx, constant.ContextKey(constant.ContextKeyLogger), logger)
			c.Set(constant.ContextKeyContext, ctx)
			defer func() {
				logger.Info("message processed", slog.Duration("duration", time.Since(start)))
			}()
			return next(c)
		}
	}
}
