package migrations

import (
	"context"
	"fmt"
	"recally/internal/pkg/logger"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

// Migrate runs all database migrations (Atlas + River)
func Migrate(ctx context.Context, databaseURL string) {
	logger.Default.Info("start migrating database")

	// Run River migrations first (for background job queue)
	migrateRiver(ctx, databaseURL)

	// Run Atlas migrations for application schema
	if err := MigrateAtlas(ctx, databaseURL); err != nil {
		logger.Default.Fatal("Error while running Atlas migrations", "err", err)
		return
	}

	logger.Default.Info("Migration successful")
}

func migrateRiver(ctx context.Context, databaseURL string) {
	logger.Default.Info("start migrating river")
	dbPool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Default.Fatal("migrateRiver failed to connect to database", "err", err)
	}
	defer dbPool.Close()

	migrator, err := rivermigrate.New(riverpgxv5.New(dbPool), nil)
	if err != nil {
		logger.Default.Fatal("migrateRiver failed to create migrator", "err", err)
	}
	res, err := migrator.Migrate(ctx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{
		TargetVersion: -1,
	})
	if err != nil {
		logger.Default.Fatal("migrateRiver failed to migrate", "err", err)
	}
	for _, version := range res.Versions {
		logger.Default.Info(fmt.Sprintf("Migrated [%s] version %d", strings.ToUpper(string(res.Direction)), version.Version))
	}
	logger.Default.Info("migrate river successful")
}
