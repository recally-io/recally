package workers

import "vibrain/internal/pkg/cache"

type Worker struct {
	cache *cache.DbCache
}

func New() *Worker {
	return &Worker{}
}

type Option func(*Worker)

func WithCache(cache *cache.DbCache) Option {
	return func(w *Worker) {
		w.cache = cache
	}
}
