package bots

import (
	"context"
	"fmt"
	"recally/internal/pkg/config"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"runtime/debug"
	"time"

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
				logger.Info("message processed", "duration", time.Since(start).Milliseconds())
			}()
			return next(c)
		}
	}
}

func recoverErrorHandler(err error, c tele.Context) {
	if config.Settings.Debug {
		debug.PrintStack()
	}
	logger.FromContext(c.Get(contexts.ContextKeyContext).(context.Context)).Error("recovered from panic", "err", err, "trace", string(debug.Stack()))
	_ = c.Reply("Something went wrong, Please try again later")
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
					} else {
						logger.FromContext(ctx).Error("transaction rollbacked after panic")
					}
					panic(r)
				}
			}()

			if err := next(c); err != nil {
				logger.FromContext(ctx).Error("transaction rollbacked after failed to process message", "err", err)
				return tx.Rollback(context.Background())
			}
			err = tx.Commit(context.Background())
			logger.FromContext(ctx).Debug("transaction commited", "err", err)
			return err
		}
	}
}
