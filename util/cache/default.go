package cache

import (
	"sync"
)

var (
	memoryCache     *MemoryCache
	memotyCacheOnce sync.Once
)

func GetMemoryCache() *MemoryCache {
	memotyCacheOnce.Do(func() {
		initMemoryCache()
	})

	return memoryCache
}

func initMemoryCache() {
	memoryCache = New()
}
