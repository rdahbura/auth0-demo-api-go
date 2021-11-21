package cache

import (
	"sync"
	"time"
)

type Item struct {
	Key   string
	Value interface{}
}

type ItemPolicy struct {
	AbsoluteExp time.Time
}

type MemoryCache struct {
	cache    map[string]cacheEntry
	cacheMtx sync.RWMutex
	ticker   *time.Ticker
}

type cacheEntry struct {
	value       interface{}
	absoluteExp time.Time
}

func New() *MemoryCache {
	mc := MemoryCache{
		cache:    map[string]cacheEntry{},
		cacheMtx: sync.RWMutex{},
		ticker:   time.NewTicker(time.Minute * 1),
	}

	go mc.startCleaner()

	return &mc
}

func (mc *MemoryCache) Add(item Item, itemPolicy ItemPolicy) bool {
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

func (mc *MemoryCache) Set(item Item, itemPolicy ItemPolicy) {
	mc.cacheMtx.Lock()
	defer mc.cacheMtx.Unlock()

	mc.set(item, itemPolicy)
}

func (mc *MemoryCache) clean() {
	mc.cacheMtx.Lock()
	defer mc.cacheMtx.Unlock()

	now := time.Now()
	for k, v := range mc.cache {
		if now.After(v.absoluteExp) {
			delete(mc.cache, k)
		}
	}
}

func (mc *MemoryCache) get(key string) (interface{}, bool) {
	cacheEntry, ok := mc.cache[key]
	if ok {
		if time.Now().Before(cacheEntry.absoluteExp) {
			return cacheEntry.value, true
		}
	}

	return nil, false
}

func (mc *MemoryCache) set(item Item, itemPolicy ItemPolicy) {
	mc.cache[item.Key] = cacheEntry{
		value:       item.Value,
		absoluteExp: itemPolicy.AbsoluteExp,
	}
}

func (mc *MemoryCache) startCleaner() {
	for {
		<-mc.ticker.C

		mc.clean()
	}
}
