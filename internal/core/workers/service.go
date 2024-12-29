package workers

import "recally/internal/pkg/cache"

type Worker struct {
	cache cache.Cache
}

func New(cache cache.Cache) *Worker {
	return &Worker{
		cache: cache,
	}
}
