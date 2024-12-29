package db

import (
	"context"
	"recally/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, databaseUrl string) (*Pool, error) {
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		logger.Default.Error("failed to connect to database", "err", err)
		return nil, err
	}
	return &Pool{
		Pool: pool,
	}, nil
}
