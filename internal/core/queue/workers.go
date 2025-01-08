package queue

import (
	"context"
	"fmt"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

func NewDefaultWorkers(llm *llms.LLM, dbPool *pgxpool.Pool) *river.Workers {
	workers := river.NewWorkers()
	dao := db.New()
	river.AddWorker(workers, NewAttachmentEmbeddingWorker(llm, dao, dbPool))
	river.AddWorker(workers, NewCrawlerWorker(llm, dbPool))
	river.AddWorker(workers, NewSummarierWorker(llm, dbPool))
	return workers
}

func loadAndSetUserContext(ctx context.Context, tx pgx.Tx, userId uuid.UUID) (context.Context, error) {
	dao := db.New()
	dbUser, err := dao.GetUserById(ctx, tx, userId)
	if err != nil {
		return ctx, fmt.Errorf("failed to load user: %w", err)
	}

	user := new(auth.UserDTO)
	user.Load(&dbUser)
	ctx = auth.SetUserToContext(ctx, user)
	return ctx, nil
}

func runInTransaction(ctx context.Context, dbPool *pgxpool.Pool, userId uuid.UUID, f func(context.Context, pgx.Tx) error) error {
	workFunc := func(ctx context.Context, tx pgx.Tx) error {
		ctx, err := loadAndSetUserContext(ctx, tx, userId)
		if err != nil {
			logger.FromContext(ctx).Error("failed to load user context", "err", err)
			return err
		}
		return f(ctx, tx)
	}

	return db.RunInTransaction(ctx, dbPool, workFunc)
}
