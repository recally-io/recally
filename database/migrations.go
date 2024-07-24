package migrations

import (
	"context"
	"vibrain/internal/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

func Migrate(ctx context.Context, databaseURL string) {
	logger.Default.Info("Migrating database")
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
			logger.Default.Info("No migration needed")
			return
		}
		logger.Default.Fatal("Error while migrating", err, "err")
	}
	logger.Default.Info("Migration successful")
}
