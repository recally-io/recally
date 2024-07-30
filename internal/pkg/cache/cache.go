package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(key CacheKey, value interface{}, expiration time.Duration)
	SetWithContext(ctx context.Context, key CacheKey, value interface{}, expiration time.Duration)
	Get(key CacheKey) ([]byte, bool)
	GetWithContext(ctx context.Context, key CacheKey) ([]byte, bool)
	Delete(key CacheKey)
	DeleteWithContext(ctx context.Context, key CacheKey)
	DeleteExpired()
	DeleteExpiredWithContext(ctx context.Context)
}

func Get[T any](ctx context.Context, c Cache, key CacheKey) (*T, bool) {
	data, ok := c.GetWithContext(ctx, key)
	if !ok {
		return nil, false
	}

	var value T
	MustUnmarshaler(data, &value)
	return &value, true
}
