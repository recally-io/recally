package queue

import (
	"context"
	"fmt"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

type Queue struct {
	*river.Client[pgx.Tx]
	name         string
	workers      *river.Workers
	periodicJobs []*river.PeriodicJob
	databaseUrl  string
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

func New(databaseUrl string, opts ...Option) (*Queue, error) {
	ctx := context.Background()

	q := &Queue{
		name:         river.QueueDefault,
		workers:      NewDefaultWorkers(),
		periodicJobs: NewDefaultPeriodJobs(),
		databaseUrl:  databaseUrl,
	}
	for _, opt := range opts {
		opt(q)
	}

	dbPool, err := pgxpool.New(ctx, q.databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	cfg := &river.Config{
		Queues: map[string]river.QueueConfig{
			q.name: {MaxWorkers: 100},
		},
		PeriodicJobs: q.periodicJobs,
		Workers:      q.workers,
	}

	riverClient, err := river.NewClient(riverpgxv5.New(dbPool), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create river client: %w", err)
	}

	q.Client = riverClient

	return q, nil
}

type Service struct {
	q *Queue
}

func NewServer() (*Service, error) {
	q, err := New(config.Settings.QueueDatabaseURL)
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
