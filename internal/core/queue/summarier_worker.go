package queue

import (
	"context"
	"recally/internal/core/bookmarks"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type SummarierWorkerArgs struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (SummarierWorkerArgs) Kind() string {
	return "content_summarier"
}

func NewSummarierWorker(llm *llms.LLM, dbPool *pgxpool.Pool) *SummarierWorker {
	return &SummarierWorker{
		llm:    llm,
		dbPool: dbPool,
	}
}

type SummarierWorker struct {
	river.WorkerDefaults[SummarierWorkerArgs]
	llm    *llms.LLM
	dbPool *pgxpool.Pool
}

func (w *SummarierWorker) Work(ctx context.Context, job *river.Job[SummarierWorkerArgs]) error {
	if err := runInTransaction(ctx, w.dbPool, job.Args.UserID, func(ctx context.Context, tx pgx.Tx) error {
		return w.work(ctx, tx, job)
	}); err != nil {
		logger.FromContext(ctx).Error("failed to run job", "err", err, "job", job)
		return err
	}

	return nil
}

func (w *SummarierWorker) work(ctx context.Context, tx pgx.Tx, job *river.Job[SummarierWorkerArgs]) error {
	svc := bookmarks.NewService(w.llm)
	dto, err := svc.SummarierContent(ctx, tx, job.Args.ID, job.Args.UserID)
	if err != nil {
		logger.FromContext(ctx).Error("failed to fetch bookmark", "err", err, "id", job.Args.ID)
		return err
	}
	logger.FromContext(ctx).Info("fetched bookmark", "id", dto.ID, "title", dto.Title, "url", dto.URL)
	return nil
}
