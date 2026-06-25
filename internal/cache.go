package internal

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mutex   sync.RWMutex
}

func NewCache(lifetime time.Duration) *Cache {
	cache := Cache{
		entries: map[string]cacheEntry{},
	}

	go cache.reapLoop(lifetime)
	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	return entry.val, exists
}

func (c *Cache) reapLoop(lifetime time.Duration) {
	ticker := time.NewTicker(lifetime)

	for {
		now := <-ticker.C

		c.mutex.Lock()
		for key, entry := range c.entries {
			if now.Sub(entry.createdAt) > lifetime {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()
	}
}
