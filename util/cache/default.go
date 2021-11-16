package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Key   string
	Value interface{}
}

type CacheItemPolicy struct {
	AbsExp time.Time
}

type MemoryCache struct {
	cache    map[string]cacheEntry
	cacheMtx sync.RWMutex
	ticker   *time.Ticker
}

type cacheEntry struct {
	value  interface{}
	absExp time.Time
}

func New() *MemoryCache {
	mc := MemoryCache{
		cache:    map[string]cacheEntry{},
		cacheMtx: sync.RWMutex{},
		ticker:   time.NewTicker(time.Minute * 1),
	}

	go mc.run()

	return &mc
}

func (mc *MemoryCache) Add(item CacheItem, itemPolicy CacheItemPolicy) bool {
	mc.cacheMtx.Lock()
	defer mc.cacheMtx.Unlock()

	_, ok := mc.get(item.Key)
	if ok {
		return false
	}

	mc.set(item, itemPolicy)

	return true
}

func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.cacheMtx.RLock()
	defer mc.cacheMtx.RUnlock()

	cacheEntry, ok := mc.get(key)
	if ok {
		return cacheEntry, true
	}

	return nil, false
}

func (mc *MemoryCache) GetCount() int {
	return len(mc.cache)
}

func (mc *MemoryCache) Set(item CacheItem, itemPolicy CacheItemPolicy) {
	mc.cacheMtx.Lock()
	defer mc.cacheMtx.Unlock()

	mc.set(item, itemPolicy)
}

func (mc *MemoryCache) clean() {
	mc.cacheMtx.Lock()
	defer mc.cacheMtx.Unlock()

	for k, v := range mc.cache {
		if time.Now().After(v.absExp) {
			delete(mc.cache, k)
		}
	}
}

func (mc *MemoryCache) get(key string) (interface{}, bool) {
	cacheEntry, ok := mc.cache[key]
	if ok {
		if time.Now().Before(cacheEntry.absExp) {
			return cacheEntry.value, true
		}
	}

	return nil, false
}

func (mc *MemoryCache) run() {
	for {
		<-mc.ticker.C
		mc.clean()
	}
}

func (mc *MemoryCache) set(item CacheItem, itemPolicy CacheItemPolicy) {
	mc.cache[item.Key] = cacheEntry{
		value:  item.Value,
		absExp: itemPolicy.AbsExp,
	}
}
