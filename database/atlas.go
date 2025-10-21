package migrations

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"recally/internal/pkg/logger"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AtlasMigrator manages database migrations using Atlas.
type AtlasMigrator struct {
	client *atlasexec.Client
	dbURL  string
}

// NewAtlasMigrator creates a new Atlas migration manager.
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

// tableExists checks if a table exists in the database.
func tableExists(ctx context.Context, dbURL, tableName string) (bool, error) {
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return false, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	var exists bool

	query := `SELECT EXISTS (
		SELECT FROM information_schema.tables
		WHERE table_schema = 'public'
		AND table_name = $1
	)`

	err = pool.QueryRow(ctx, query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if table exists: %w", err)
	}

	return exists, nil
}

// atlasTrackingExists checks if Atlas migration tracking exists.
func (m *AtlasMigrator) atlasTrackingExists(ctx context.Context) (bool, error) {
	pool, err := pgxpool.New(ctx, m.dbURL)
	if err != nil {
		return false, fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	var exists bool
	// Atlas creates a schema named 'atlas_schema_revisions' with a table of the same name
	query := `SELECT EXISTS (
		SELECT FROM information_schema.schemata
		WHERE schema_name = 'atlas_schema_revisions'
	)`

	err = pool.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if atlas schema exists: %w", err)
	}

	return exists, nil
}

// needsBaseline checks if database needs to be baselined.
func (m *AtlasMigrator) needsBaseline(ctx context.Context) (bool, error) {
	// Check if the users table exists (indicator of go-migrate having run)
	usersExists, err := tableExists(ctx, m.dbURL, "users")
	if err != nil {
		return false, fmt.Errorf("failed to check if users table exists: %w", err)
	}

	// If users table doesn't exist, this is a fresh database - no baseline needed
	if !usersExists {
		logger.Default.Info("Fresh database detected, no baseline needed")

		return false, nil
	}

	// Check if Atlas migration table exists
	// Note: Atlas creates its tracking table in a schema with the same name
	atlasTableExists, err := m.atlasTrackingExists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check if atlas table exists: %w", err)
	}

	// If both users table and atlas table exist, migrations are already tracked
	if atlasTableExists {
		logger.Default.Info("Atlas migration tracking already exists")

		return false, nil
	}

	// Users table exists but Atlas table doesn't - need to baseline
	logger.Default.Info("Existing database detected, will baseline with Atlas")

	return true, nil
}

// Migrate runs all pending migrations.
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

	// Check if database needs to be baselined (for go-migrate -> Atlas transition)
	needsBaseline, err := m.needsBaseline(ctx)
	if err != nil {
		return fmt.Errorf("failed to check baseline status: %w", err)
	}

	// Prepare migration apply parameters
	params := &atlasexec.MigrateApplyParams{
		URL:    m.dbURL,
		DirURL: fmt.Sprintf("file://%s", migrationsDir),
		// Allow dirty mode to handle databases with extension schemas (like paradedb)
		AllowDirty: true,
	}

	// If baseline is needed, set the baseline version to the initial migration
	if needsBaseline {
		params.BaselineVersion = "20251021000000"

		logger.Default.Info("Baselining database with initial migration version")
	}

	// Apply migrations
	result, err := m.client.MigrateApply(ctx, params)
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

// Status returns the current migration status.
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

// MigrateAtlas runs Atlas migrations (replaces the go-migrate portion).
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
