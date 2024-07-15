package cache

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// var instance *cache.Cache
var (
	instance *cache.Cache
	once     sync.Once
)

// New creates a new cache instance
func New() *cache.Cache {
	once.Do(func() {
		instance = cache.New(24*time.Hour, 36*time.Hour)
	})

	return instance
}
