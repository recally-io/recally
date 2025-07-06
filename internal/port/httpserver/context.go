package httpserver

import (
	"context"
	"errors"

	"recally/internal/pkg/auth"
	"recally/internal/pkg/contexts"

	"github.com/jackc/pgx/v5"
)

func loadTx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := contexts.Get[pgx.Tx](ctx, contexts.ContextKeyTx)
	if !ok {
		return nil, errors.New("tx not found")
	}

	return tx, nil
}

func initContext(ctx context.Context) (pgx.Tx, *auth.UserDTO, error) {
	tx, err := loadTx(ctx)
	if err != nil {
		return nil, nil, errors.New("tx not found")
	}

	user, ok := contexts.Get[*auth.UserDTO](ctx, contexts.ContextKeyUser)
	if !ok {
		return tx, nil, errors.New("user not found")
	}

	return tx, user, nil
}
