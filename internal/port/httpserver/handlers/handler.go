package handlers

import (
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
)

type Handler struct {
	Pool   *db.Pool
	Cache  *cache.DbCache
	worker *workers.Worker
}

func New(pool *db.Pool) *Handler {
	h := &Handler{
		Pool: pool,
	}

	workerService := workers.New()
	if h.Cache != nil {
		workers.WithCache(h.Cache)(workerService)
	}

	return h
}
