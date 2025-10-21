package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"recally/internal/pkg/logger"

	"ariga.io/atlas-go-sdk/atlasexec"
)

// AtlasMigrator manages database migrations using Atlas
type AtlasMigrator struct {
	client *atlasexec.Client
	dbURL  string
}

// NewAtlasMigrator creates a new Atlas migration manager
func NewAtlasMigrator(databaseURL string) (*AtlasMigrator, error) {
	// Get the working directory to find the database folder
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Create Atlas client
	client, err := atlasexec.NewClient(workDir, "atlas")
	if err != nil {
		return nil, fmt.Errorf("failed to create Atlas client: %w", err)
	}

	return &AtlasMigrator{
		client: client,
		dbURL:  databaseURL,
	}, nil
}

// Migrate runs all pending migrations
func (m *AtlasMigrator) Migrate(ctx context.Context) error {
	logger.Default.Info("Running Atlas migrations")

	// Get the migrations directory path
	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationsDir := filepath.Join(workDir, "database", "migrations")

	// Check if migrations directory exists
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", migrationsDir)
	}

	// Apply migrations
	result, err := m.client.MigrateApply(ctx, &atlasexec.MigrateApplyParams{
		URL:    m.dbURL,
		DirURL: fmt.Sprintf("file://%s", migrationsDir),
	})

	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Log migration results
	if result.Error != "" {
		return fmt.Errorf("migration error: %s", result.Error)
	}

	if len(result.Applied) == 0 {
		logger.Default.Info("No pending migrations to apply")
	} else {
		for _, applied := range result.Applied {
			logger.Default.Info(fmt.Sprintf("Applied migration: %s", applied.Name))
		}
		logger.Default.Info(fmt.Sprintf("Successfully applied %d migration(s)", len(result.Applied)))
	}

	return nil
}

// Status returns the current migration status
func (m *AtlasMigrator) Status(ctx context.Context) (*atlasexec.MigrateStatus, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}
	migrationsDir := filepath.Join(workDir, "database", "migrations")

	status, err := m.client.MigrateStatus(ctx, &atlasexec.MigrateStatusParams{
		URL:    m.dbURL,
		DirURL: fmt.Sprintf("file://%s", migrationsDir),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get migration status: %w", err)
	}

	return status, nil
}

// MigrateAtlas runs Atlas migrations (replaces the go-migrate portion)
func MigrateAtlas(ctx context.Context, databaseURL string) error {
	migrator, err := NewAtlasMigrator(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create Atlas migrator: %w", err)
	}

	if err := migrator.Migrate(ctx); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
