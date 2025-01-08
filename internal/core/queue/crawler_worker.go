package queue

import (
	"context"
	"recally/internal/core/bookmarks"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
	"recally/internal/pkg/webreader/fetcher"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type CrawlerWorkerArgs struct {
	ID          uuid.UUID          `json:"id"`
	UserID      uuid.UUID          `json:"user_id"`
	FetcherName fetcher.FecherType `json:"fetcher_name"`
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
	if err := runInTransaction(ctx, w.dbPool, job.Args.UserID, func(ctx context.Context, tx pgx.Tx) error {
		return w.work(ctx, tx, job)
	}); err != nil {
		logger.FromContext(ctx).Error("failed to run job", "err", err, "job", job)
		return err
	}

	return nil
}

func (w *CrawlerWorker) work(ctx context.Context, tx pgx.Tx, job *river.Job[CrawlerWorkerArgs]) error {
	svc := bookmarks.NewService(w.llm)
	dto, err := svc.FetchContent(ctx, tx, job.Args.ID, job.Args.UserID, job.Args.FetcherName)
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch bookmark", "err", err, "id", job.Args.ID, "fetcher", job.Args.FetcherName)
		return err
	}

	go func() {
		ctx := logger.CopyContext(ctx)
		if dto.Content != "" && dto.Summary == "" {
			dto, err = svc.SummarierContent(ctx, tx, job.Args.ID, job.Args.UserID)
			if err != nil {
				logger.FromContext(ctx).Error("failed to summarise content", "err", err)
			}
		}
	}()

	logger.FromContext(ctx).Info("fetched bookmark", "id", dto.ID, "title", dto.Title, "url", dto.URL, "fetcher", job.Args.FetcherName)
	return nil
}
