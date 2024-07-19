package handlers

import (
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
)

type Handler struct {
	Pool   *db.Pool
	Cache  *cache.Cache
	worker *workers.Worker
}

func New(pool *db.Pool, opts ...Option) *Handler {
	h := &Handler{
		Pool: pool,
	}
	for _, opt := range opts {
		opt(h)
	}

	workerService := workers.New()
	if h.Cache != nil {
		workers.WithCache(h.Cache)(workerService)
	}

	return h
}

type Option func(*Handler)

func WithCache(c *cache.Cache) Option {
	return func(s *Handler) {
		s.Cache = c
	}
}
