package workers

import "vibrain/internal/pkg/cache"

type Worker struct {
	cache *cache.Cache
}

func New() *Worker {
	return &Worker{}
}

type Option func(*Worker)

func WithCache(cache *cache.Cache) Option {
	return func(w *Worker) {
		w.cache = cache
	}
}
