package bots

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

func contextMiddleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := context.Background()
			start := time.Now()
			ctx = contexts.Set(ctx, contexts.ContextKeyRequestID, uuid.Must(uuid.NewV7()).String())
			ctx = contexts.Set(ctx, contexts.ContextKeyUserID, fmt.Sprintf("%d", c.Sender().ID))
			ctx = contexts.Set(ctx, contexts.ContextKeyUserName, c.Sender().Username)
			logger := logger.FromContext(ctx)
			ctx = contexts.Set(ctx, contexts.ContextKeyLogger, logger)
			c.Set(contexts.ContextKeyContext, ctx)
			defer func() {
				logger.Info("message processed", slog.Duration("duration", time.Since(start)))
			}()
			return next(c)
		}
	}
}

func TransactionMiddleware(pool *db.Pool) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			ctx := c.Get(contexts.ContextKeyContext).(context.Context)
			tx, err := pool.Begin(ctx)
			if err != nil {
				return err
			}
			ctx = contexts.Set(ctx, contexts.ContextKeyTx, tx)
			c.Set(contexts.ContextKeyContext, ctx)
			defer func() {
				if r := recover(); r != nil {
					if err := tx.Rollback(context.Background()); err != nil {
						logger.FromContext(ctx).Error("failed to rollback transaction", "err", err)
					}
					panic(r)
				}
			}()

			if err := next(c); err != nil {
				return tx.Rollback(context.Background())
			}

			return tx.Commit(context.Background())
		}
	}
}
