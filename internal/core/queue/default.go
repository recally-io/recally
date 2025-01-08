package queue

import (
	"fmt"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
	"recally/internal/pkg/logger"
)

var DefaultQueue *Queue

func init() {
	var err error
	DefaultQueue, err = NewDefaultQueue()
	if err != nil {
		logger.Default.Error("failed to create default queue", "err", err)
	}
}

func NewDefaultQueue() (*Queue, error) {
	// start queue service
	q, err := New(db.DefaultPool, llms.DefaultLLM)
	if err != nil {
		return nil, fmt.Errorf("failed to create new queue service: %w", err)
	}
	return q, nil
}
