package migrations

import (
	"context"
	"fmt"
	"recally/internal/pkg/logger"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
)

func Migrate(ctx context.Context, databaseURL string) {
	logger.Default.Info("start migrating database")
	migrateRiver(ctx, databaseURL)

	s := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})
	d, err := bindata.WithInstance(s)
	if err != nil {
		logger.Default.Fatal("Error while creating bindata instance", err, "err")
	}

	m, err := migrate.NewWithSourceInstance("migrations", d, databaseURL)
	if err != nil {
		logger.Default.Fatal("Error while creating migrate instance", err, "err")
	}
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			logger.Default.Info("No migration for normal db needed")
			return
		}
		logger.Default.Fatal("Error while migrating", "err", err)
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

	migrator := rivermigrate.New(riverpgxv5.New(dbPool), nil)
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.Default.Fatal("migrateRiver failed to start transaction", "err", err)
	}

	defer func() {
		if err := tx.Commit(ctx); err != nil {
			logger.Default.Error("Failed to commit transaction", "err", err)
		}
	}()
	res, err := migrator.MigrateTx(ctx, tx, rivermigrate.DirectionUp, &rivermigrate.MigrateOpts{
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
