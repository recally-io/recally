package httpserver

import (
	"context"
	"errors"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"
	"recally/internal/pkg/db"

	"github.com/google/uuid"
)

func loadTx(ctx context.Context) (db.DBTX, error) {
	tx, ok := contexts.Get[db.DBTX](ctx, contexts.ContextKeyTx)
	if !ok {
		return nil, errors.New("tx not found")
	}
	return tx, nil
}

func initContext(ctx context.Context) (db.DBTX, *auth.UserDTO, error) {
	tx, err := loadTx(ctx)
	if err != nil {
		return nil, nil, errors.New("tx not found")
	}

	userId, ok := contexts.Get[uuid.UUID](ctx, contexts.ContextKeyUserID)
	if !ok {
		return tx, nil, errors.New("user not found")
	}

	user, err := auth.New().GetUserById(ctx, tx, userId)
	if err != nil {
		return tx, nil, errors.New("user not found")
	}

	return tx, user, nil
}
