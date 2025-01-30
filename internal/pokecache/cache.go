package pokecache

import (
	"time"
)

type Cache map[string]cacheEntry

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
	(*c)[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
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
		for k, v := range *c {
			if time.Since(v.createdAt) > interval {
				delete(*c, k)
			}
		}
	}
}
