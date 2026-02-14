package handlers

import (
	"recally/internal/core/bookmarks"
	"recally/internal/core/queue"
	"recally/internal/pkg/auth"
	"recally/internal/pkg/cache"
	"recally/internal/pkg/db"
	"recally/internal/pkg/llms"
)

type Handler struct {
	pool  *db.Pool
	llm   *llms.LLM
	cache cache.Cache
	queue *queue.Queue

	authService     *auth.Service
	bookmarkService *bookmarks.Service
}

func New(pool *db.Pool, llm *llms.LLM, queue *queue.Queue, opts ...Option) *Handler {
	h := &Handler{
		pool:             pool,
		llm:              llm,
		cache:            cache.MemCache,
		queue:            queue,
		authService:     auth.New(),
		bookmarkService: bookmarks.NewService(llm),
	}
	for _, opt := range opts {
		opt(h)
	}

	return h
}

type Option func(*Handler)

func WithCache(c cache.Cache) Option {
	return func(s *Handler) {
		s.cache = c
	}
}
