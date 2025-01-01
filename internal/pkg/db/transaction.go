package db

import (
	"context"
	"fmt"
	"recally/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RunInTransaction executes the provided function within a database transaction.
// It handles the transaction lifecycle including begin, commit and rollback operations.
//
// Parameters:
//   - ctx: The context.Context for the transaction
//   - dbPool: A connection pool to the PostgreSQL database
//   - f: A function that takes a context and transaction, and returns an error
//
// Returns:
//   - error: Returns nil on successful commit, or an error if transaction operations fail
//
// The transaction will be rolled back if the provided function returns an error.
// Otherwise, the transaction will be committed.
func RunInTransaction(ctx context.Context, dbPool *pgxpool.Pool, f func(context.Context, pgx.Tx) error) error {
	tx, err := dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to load transaction: %w", err)
	}

	if err := f(ctx, tx); err != nil {
		logger.FromContext(ctx).Error("rollback transaction", "err", err)
		if err := tx.Rollback(ctx); err != nil {
			logger.FromContext(ctx).Error("failed to rollback transaction", "err", err)
		}
		return err
	}

	return tx.Commit(ctx)
}
