package cache

import (
	"context"
	"encoding/json"
	"recally/internal/pkg/db"
	"recally/internal/pkg/logger"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

var DefaultDBCache *DbCache

func init() {
	DefaultDBCache = NewDBCache(db.DefaultPool)
}

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
	tx *db.Pool
	db *db.Queries
}

type Option func(*DbCache)

// NewDBCache creates a new cache instance.
func NewDBCache(pool *db.Pool, opts ...Option) *DbCache {
	service := &DbCache{
		tx: pool,
		db: db.New(),
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (c *DbCache) Set(key CacheKey, value interface{}, expiration time.Duration) {
	c.SetWithContext(context.Background(), key, value, expiration)
}

func (c *DbCache) SetWithContext(ctx context.Context, key CacheKey, value interface{}, expiration time.Duration) {
	ok, err := c.db.IsCacheExists(ctx, c.tx, db.IsCacheExistsParams{Domain: key.Domain, Key: key.Key})
	if err != nil {
		logger.FromContext(ctx).Warn("failed to check cache exists", "key", key, "err", err)

		return
	}

	jsonValue, err := Marshaler(value)
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

		if err := c.db.CreateCache(ctx, c.tx, params); err != nil {
			logger.FromContext(ctx).Warn("failed to create cache", "key", key, "err", err)
		}
	} else {
		params := db.UpdateCacheParams{
			Domain:    key.Domain,
			Key:       key.Key,
			Value:     jsonValue,
			ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(expiration), Valid: true},
		}
		if err := c.db.UpdateCache(ctx, c.tx, params); err != nil {
			logger.FromContext(ctx).Warn("failed to update cache", "key", key, "err", err)
		}
	}
}

func (c *DbCache) Get(key CacheKey) (any, bool) {
	return c.GetWithContext(context.Background(), key)
}

func (c *DbCache) GetWithContext(ctx context.Context, key CacheKey) (any, bool) {
	item, err := c.db.GetCacheByKey(ctx, c.tx, db.GetCacheByKeyParams{
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
	if err := c.db.DeleteCacheByKey(ctx, c.tx, db.DeleteCacheByKeyParams{
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
	if err := c.db.DeleteExpiredCache(ctx, c.tx, pgtype.Timestamp{Time: time.Now()}); err != nil {
		logger.FromContext(ctx).Warn("failed to delete expired cache", "err", err)
	}
}

func Marshaler[T any](value T) ([]byte, error) {
	return json.Marshal(value)
}

func MustUnmarshaler[T any](data []byte, value *T) {
	if err := json.Unmarshal(data, value); err != nil {
		logger.Default.Fatal("failed to unmarshal cache value", "err", err)
	}
}
