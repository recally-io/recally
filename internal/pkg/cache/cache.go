package cache

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"vibrain/internal/pkg/db"
	"vibrain/internal/pkg/logger"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/patrickmn/go-cache"
)

// var instance *cache.Cache
var (
	instance *cache.Cache
	once     sync.Once
)

type Cache struct {
	MemCache *cache.Cache
	Pool     *db.Pool
}

type Option func(*Cache)

// New creates a new cache instance
func New(opts ...Option) *Cache {
	service := &Cache{}
	once.Do(func() {
		instance = cache.New(24*time.Hour, 36*time.Hour)
	})
	service.MemCache = instance
	for _, opt := range opts {
		opt(service)
	}
	return service
}

func WithDB(pool *db.Pool) Option {
	return func(c *Cache) {
		c.Pool = pool
	}
}

func (c *Cache) getDb() *db.Queries {
	return db.New(c.Pool)
}

func (c *Cache) isUsingDB() bool {
	return c.Pool != nil
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.SetWithContext(context.Background(), key, value, expiration)
}

func (c *Cache) SetWithContext(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	// c.Cache.Set(key, value, expiration)
	if c.isUsingDB() {
		item, ok := c.GetWithContext(ctx, key)
		if !ok {
			jsonValue, err := json.Marshal(value)
			if err != nil {
				logger.FromContext(ctx).Warn("failed to marshal value", "key", key, "err", err)
				return
			}

			params := db.CreateCacheParams{
				Key:       key,
				Value:     pgtype.Text{String: string(jsonValue), Valid: true},
				ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(expiration), Valid: true},
				CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}

			if err := c.getDb().CreateCache(ctx, params); err != nil {
				logger.FromContext(ctx).Warn("failed to create cache", "key", key, "err", err)
			}
		} else {
			params := db.UpdateCacheParams{
				Key:       key,
				Value:     pgtype.Text{String: item.(string), Valid: true},
				ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(expiration), Valid: true},
				UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
			}
			if err := c.getDb().UpdateCache(ctx, params); err != nil {
				logger.FromContext(ctx).Warn("failed to update cache", "key", key, "err", err)
			}
		}
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.GetWithContext(context.Background(), key)
}

func (c *Cache) GetWithContext(ctx context.Context, key string) (interface{}, bool) {
	if val, ok := c.MemCache.Get(key); ok {
		return val, ok
	}

	if c.isUsingDB() {
		item, err := c.getDb().GetCacheByKey(ctx, key)
		if err != nil {
			return nil, false
		}
		c.MemCache.Set(key, item.Value, time.Until(item.ExpiresAt.Time))
		return item.Value, true
	}
	return nil, false
}

func (c *Cache) Delete(key string) {
	c.DeleteWithContext(context.Background(), key)
}

func (c *Cache) DeleteWithContext(ctx context.Context, key string) {
	c.MemCache.Delete(key)
	if c.isUsingDB() {
		if err := c.getDb().DeleteCacheByKey(ctx, key); err != nil {
			logger.FromContext(ctx).Warn("failed to delete cache", "key", key, "err", err)
		}
	}
}

func (c *Cache) DeleteExpired() {
	c.DeleteExpiredWithContext(context.Background())
}

func (c *Cache) DeleteExpiredWithContext(ctx context.Context) {
	c.MemCache.DeleteExpired()
	if c.isUsingDB() {
		if err := c.getDb().DeleteExpiredCache(ctx, pgtype.Timestamp{Time: time.Now()}); err != nil {
			logger.FromContext(ctx).Warn("failed to delete expired cache", "err", err)
		}
	}
}
