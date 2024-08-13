package handlers

import (
	"vibrain/internal/core/assistants"
	"vibrain/internal/core/workers"
	"vibrain/internal/pkg/cache"
	"vibrain/internal/pkg/db"
)

type Handler struct {
	Pool      *db.Pool
	Cache     cache.Cache
	worker    *workers.Worker
	assistant *assistants.Service
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
