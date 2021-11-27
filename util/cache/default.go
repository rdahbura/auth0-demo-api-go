package cache

import (
	"sync"
)

var (
	memoryCache     *MemoryCache
	memoryCacheOnce sync.Once
)

func GetMemoryCache() *MemoryCache {
	memoryCacheOnce.Do(initMemoryCache)

	return memoryCache
}

func initMemoryCache() {
	memoryCache = New()
}
