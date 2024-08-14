package handlers

import (
	"vibrain/internal/core/assistants"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
)

type Handler struct {
	pool  *db.Pool
	llm   *llms.LLM
	cache cache.Cache

	toolService      *workers.Worker
	assistantService *assistants.Service
}

func New(pool *db.Pool, llm *llms.LLM, opts ...Option) *Handler {
	h := &Handler{
		pool:  pool,
		llm:   llm,
		cache: cache.MemCache,
	}
	for _, opt := range opts {
		opt(h)
	}
	h.toolService = workers.New(h.cache)
	h.assistantService = assistants.NewService(h.llm)

	return h
}

type Option func(*Handler)

func WithCache(c cache.Cache) Option {
	return func(s *Handler) {
		s.cache = c
	}
}
