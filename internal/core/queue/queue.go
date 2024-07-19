package queue

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type Queue struct {
	*river.Client[pgx.Tx]
	name         string
	workers      *river.Workers
	periodicJobs []*river.PeriodicJob
}

type Option func(q *Queue)

func WithQueueName(queueName string) Option {
	return func(q *Queue) {
		q.name = queueName
	}
}

func WithWorkers(workers *river.Workers) Option {
	return func(q *Queue) {
		q.workers = workers
	}
}

func WithPeriodicJobs(jobs []*river.PeriodicJob) Option {
	return func(q *Queue) {
		q.periodicJobs = jobs
	}
}

func New(pool *db.Pool, opts ...Option) (*Queue, error) {
	q := &Queue{
		name:         river.QueueDefault,
		workers:      NewDefaultWorkers(),
		periodicJobs: NewDefaultPeriodJobs(),
	}
	for _, opt := range opts {
		opt(q)
	}
	cfg := &river.Config{
		Queues: map[string]river.QueueConfig{
			q.name: {MaxWorkers: 100},
		},
		PeriodicJobs: q.periodicJobs,
		Workers:      q.workers,
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pool.Pool), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	q.Client = riverClient

	return q, nil
}

type Service struct {
	q *Queue
}

func NewServer(pool *db.Pool) (*Service, error) {
	q, err := New(pool)
	if err != nil {
		return nil, err
	}
	return &Service{
		q: q,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	if err := s.q.Start(ctx); err != nil {
		logger.Default.Fatal("failed to start", "service", s.Name(), "error", err)
	}
}

func (s *Service) Stop(ctx context.Context) {
	if err := s.q.Stop(ctx); err != nil {
		logger.Default.Fatal("failed to stop", "service", s.Name(), "error", err)
	}
}

func (s *Service) Name() string {
	return "river queue server"
}
