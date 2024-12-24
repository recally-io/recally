package queue

import (
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

func NewDefaultWorkers(llm *llms.LLM, dbPool *pgxpool.Pool) *river.Workers {
	workers := river.NewWorkers()
	dao := db.New()
	river.AddWorker(workers, NewAttachmentEmbeddingWorker(llm, dao, dbPool))
	river.AddWorker(workers, NewCrawlerWorker(llm, dbPool))
	return workers
}
