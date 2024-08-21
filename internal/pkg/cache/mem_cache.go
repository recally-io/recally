package cache

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

var MemCache = NewMemCache(12*time.Hour, 24*time.Hour)

type memCache struct {
	c *cache.Cache
}

func NewMemCache(defaultExpiration, cleanupInterval time.Duration) *memCache {
	return &memCache{c: cache.New(defaultExpiration, cleanupInterval)}
}

func (m *memCache) Set(key CacheKey, value interface{}, expiration time.Duration) {
	m.c.Set(key.String(), value, expiration)
}

func (m *memCache) SetWithContext(ctx context.Context, key CacheKey, value interface{}, expiration time.Duration) {
	m.Set(key, value, expiration)
}

func (m *memCache) Get(key CacheKey) (any, bool) {
	value, ok := m.c.Get(key.String())
	if !ok {
		return nil, false
	}
	return value, true
}

func (m *memCache) GetWithContext(ctx context.Context, key CacheKey) (any, bool) {
	return m.Get(key)
}

func (m *memCache) Delete(key CacheKey) {
	m.c.Delete(key.String())
}

func (m *memCache) DeleteWithContext(ctx context.Context, key CacheKey) {
	m.Delete(key)
}

func (m *memCache) DeleteExpired() {
	m.c.DeleteExpired()
}

func (m *memCache) DeleteExpiredWithContext(ctx context.Context) {
	m.DeleteExpired()
}
