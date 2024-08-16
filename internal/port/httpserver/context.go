package httpserver

import (
	"context"
	"errors"
	"vibrain/internal/pkg/contexts"
	"vibrain/internal/pkg/db"
)

func loadTx(ctx context.Context) (db.DBTX, error) {
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return nil, errors.New("tx not found")
	}
	return tx, nil
}

func loadUserId(ctx context.Context) (string, error) {
	userId, ok := contexts.Get[string](ctx, contexts.ContextKeyUserID)
	if !ok {
		return "", errors.New("user not found")
	}
	return userId, nil
}
