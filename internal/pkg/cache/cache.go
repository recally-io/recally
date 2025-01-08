package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(key CacheKey, value interface{}, expiration time.Duration)
	SetWithContext(ctx context.Context, key CacheKey, value interface{}, expiration time.Duration)
	Get(key CacheKey) (any, bool)
	GetWithContext(ctx context.Context, key CacheKey) (any, bool)
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

	b, ok := data.([]byte)
	if !ok {
		return data.(*T), true
	}

	var value T
	MustUnmarshaler(b, &value)
	return &value, true
}

func RunInCache[T any](ctx context.Context, c Cache, key CacheKey, expiration time.Duration, f func() (*T, error)) (*T, error) {
	data, ok := Get[T](ctx, c, key)
	if ok && data != nil {
		return data, nil
	}

	data, err := f()
	if err != nil {
		return nil, err
	}

	c.SetWithContext(ctx, key, data, expiration)
	return data, nil
}
