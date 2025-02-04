package pokecache

import (
	"sync"
	"time"
)

type Cache map[string]cacheEntry

var mux sync.RWMutex

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{}
	go cache.ReadLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	mux.Lock()
	defer mux.Unlock()

	(*c)[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	mux.RLock()
	defer mux.RUnlock()

	entry, ok := (*c)[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) ReadLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		mux.Lock()
		for k, v := range *c {
			if time.Since(v.createdAt) > interval {
				delete(*c, k)
			}
		}
		mux.Unlock()
	}
}
