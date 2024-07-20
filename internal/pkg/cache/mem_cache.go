package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var MemCache = cache.New(12*time.Hour, 24*time.Hour)
