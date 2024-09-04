package handlers

import (
	"vibrain/internal/core/assistants"
	"vibrain/internal/core/queue"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/auth"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
)

type Handler struct {
	pool  *db.Pool
	llm   *llms.LLM
	cache cache.Cache
	queue *queue.Queue

	authService      *auth.Service
	toolService      *workers.Worker
	assistantService *assistants.Service
}

func New(pool *db.Pool, llm *llms.LLM, queue *queue.Queue, opts ...Option) *Handler {
	h := &Handler{
		pool:             pool,
		llm:              llm,
		cache:            cache.MemCache,
		queue:            queue,
		authService:      auth.New(),
		assistantService: assistants.NewService(llm, queue),
	}
	for _, opt := range opts {
		opt(h)
	}
	h.toolService = workers.New(h.cache)

	return h
}

type Option func(*Handler)

func WithCache(c cache.Cache) Option {
	return func(s *Handler) {
		s.cache = c
	}
}
