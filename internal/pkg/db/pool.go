package db

import (
	"context"
	"recally/internal/pkg/config"
	"recally/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DefaultPool *Pool

func init() {
	var err error
	DefaultPool, err = NewPool(context.Background(), config.Settings.Database.URL())
	if err != nil {
		logger.Default.Fatal("failed to create default database pool", "err", err)
	}
}

type Pool struct {
	*pgxpool.Pool
}

func NewPool(ctx context.Context, databaseUrl string) (*Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	if config.Settings.Database.DEBUG {
		poolConfig.ConnConfig.Tracer = logger.Default
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Default.Error("failed to connect to database", "err", err)
		return nil, err
	}
	return &Pool{
		Pool: pool,
	}, nil
}
