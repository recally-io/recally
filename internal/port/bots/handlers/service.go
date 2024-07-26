package handlers

import (
	"vibrain/internal/core/assistants"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/config"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/llms"
	"vibrain/internal/pkg/logger"
)

type Handler struct {
	Cache     *cache.DbCache
	worker    *workers.Worker
	assistant *assistants.Service
}

func New(pool *db.Pool, opts ...Option) *Handler {
	h := &Handler{
		worker: workers.New(),
	}
	llm := llms.New(config.Settings.OpenAI.BaseURL, config.Settings.OpenAI.ApiKey)
	ass, err := assistants.NewService(llm)
	if err != nil {
		logger.Default.Fatal("failed to create assistant service", "err", err)
	}
	h.assistant = ass

	for _, opt := range opts {
		opt(h)
	}

	if h.Cache != nil {
		workers.WithCache(h.Cache)(h.worker)
	}

	return h
}

type Option func(*Handler)

func WithCache(c *cache.DbCache) Option {
	return func(s *Handler) {
		s.Cache = c
	}
}
