package cache

import (
	"context"
	"encoding/json"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgtype"
)

type CacheKey struct {
	Domain string
	Key    string
}

func (c *CacheKey) String() string {
	return c.Domain + ":" + c.Key
}

func NewCacheKey(domain, key string) CacheKey {
	return CacheKey{Domain: domain, Key: key}
}

type DbCache struct {
	Pool *db.Pool
}

type Option func(*DbCache)

// NewDBCache creates a new cache instance
func NewDBCache(pool *db.Pool, opts ...Option) *DbCache {
	service := &DbCache{
		Pool: pool,
	}

	for _, opt := range opts {
		opt(service)
	}
	return service
}

func (c *DbCache) getConn() *db.Queries {
	return db.New(c.Pool)
}

func (c *DbCache) Set(key CacheKey, value interface{}, expiration time.Duration) {
	c.SetWithContext(context.Background(), key, value, expiration)
}

func (c *DbCache) SetWithContext(ctx context.Context, key CacheKey, value interface{}, expiration time.Duration) {
	ok, err := c.getConn().IsCacheExists(ctx, db.IsCacheExistsParams{Domain: key.Domain, Key: key.Key})
	if err != nil {
		logger.FromContext(ctx).Warn("failed to check cache exists", "key", key, "err", err)
		return
	}
	jsonValue, err := json.Marshal(value)
	if err != nil {
		logger.FromContext(ctx).Warn("failed to marshal value", "key", key, "err", err)
		return
	}
	if !ok {
		params := db.CreateCacheParams{
			Domain:    key.Domain,
			Key:       key.Key,
			Value:     jsonValue,
			ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(expiration), Valid: true},
		}

		if err := c.getConn().CreateCache(ctx, params); err != nil {
			logger.FromContext(ctx).Warn("failed to create cache", "key", key, "err", err)
		}
	} else {
		params := db.UpdateCacheParams{
			Domain:    key.Domain,
			Key:       key.Key,
			Value:     jsonValue,
			ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(expiration), Valid: true},
		}
		if err := c.getConn().UpdateCache(ctx, params); err != nil {
			logger.FromContext(ctx).Warn("failed to update cache", "key", key, "err", err)
		}
	}
}

func (c *DbCache) Get(key CacheKey) (any, bool) {
	return c.GetWithContext(context.Background(), key)
}

func (c *DbCache) GetWithContext(ctx context.Context, key CacheKey) (any, bool) {
	item, err := c.getConn().GetCacheByKey(ctx, db.GetCacheByKeyParams{
		Domain: key.Domain,
		Key:    key.Key,
	})
	if err != nil {
		logger.FromContext(ctx).Warn("failed to get cache", "key", key, "err", err)
		return nil, false
	}
	// c.MemCache.Set(key.String(), item.Value, time.Until(item.ExpiresAt.Time))
	return item.Value, true
}

func (c *DbCache) Delete(key CacheKey) {
	c.DeleteWithContext(context.Background(), key)
}

func (c *DbCache) DeleteWithContext(ctx context.Context, key CacheKey) {
	if err := c.getConn().DeleteCacheByKey(ctx, db.DeleteCacheByKeyParams{
		Key:    key.Key,
		Domain: key.Domain,
	}); err != nil {
		logger.FromContext(ctx).Warn("failed to delete cache", "key", key, "err", err)
	}
}

func (c *DbCache) DeleteExpired() {
	c.DeleteExpiredWithContext(context.Background())
}

func (c *DbCache) DeleteExpiredWithContext(ctx context.Context) {
	if err := c.getConn().DeleteExpiredCache(ctx, pgtype.Timestamp{Time: time.Now()}); err != nil {
		logger.FromContext(ctx).Warn("failed to delete expired cache", "err", err)
	}
}
