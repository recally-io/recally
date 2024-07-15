package queue

import (
	"github.com/riverqueue/river"
)

func NewDefaultWorkers() *river.Workers {
	workers := river.NewWorkers()
	river.AddWorker(workers, &EmailWorker{})
	return workers
}
