package queue

import (
	"context"
	"vibrain/internal/core/bookmarks"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type CrawlerWorkerArgs struct {
	ID          uuid.UUID            `json:"id"`
	UserID      uuid.UUID            `json:"user_id"`
	FetcherName bookmarks.FecherType `json:"fetcher_name"`
}

func (CrawlerWorkerArgs) Kind() string {
	return "web_crawler"
}

func NewCrawlerWorker(llm *llms.LLM, dbPool *pgxpool.Pool) *CrawlerWorker {
	return &CrawlerWorker{
		llm:    llm,
		dbPool: dbPool,
	}
}

type CrawlerWorker struct {
	river.WorkerDefaults[CrawlerWorkerArgs]
	llm    *llms.LLM
	dbPool *pgxpool.Pool
}

func (w *CrawlerWorker) Work(ctx context.Context, job *river.Job[CrawlerWorkerArgs]) error {
	svc := bookmarks.NewService(w.llm)
	tx, err := w.dbPool.Begin(ctx)
	if err != nil {
		logger.FromContext(ctx).Error("failed to start transaction", "error", err)
		return err
	}
	dto, err := svc.FetchContent(ctx, tx, job.Args.ID, job.Args.UserID, job.Args.FetcherName)
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch bookmark", "error", err)
		return err
	}

	if dto.Content != "" {
		dto, err = svc.SummarierContent(ctx, tx, job.Args.ID, job.Args.UserID)
		if err != nil {
			logger.FromContext(ctx).Error("failed to summarise content", "error", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.FromContext(ctx).Error("failed to commit transaction", "error", err)
	}

	logger.FromContext(ctx).Info("fetched bookmark", "id", dto.ID, "title", dto.Title)

	return nil
}
