package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu    sync.Mutex
	table map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		table: make(map[string]cacheEntry),
	}
	c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ent := c.table[key]
	ent.createdAt = time.Now()
	ent.val = data
	c.table[key] = ent
}

func (c *Cache) Get(key string) (data []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ent, found := c.table[key]
	if !found {
		return []byte{}, false
	}
	return ent.val, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.table, key)
}

func (c *Cache) ForEachEntry(cb func(key string, ent cacheEntry)) {
	for k, v := range c.table {
		cb(k, v)
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	// defer ticker.Stop()
	go func() {
		for {
			t := <-ticker.C
			c.ForEachEntry(func(key string, ent cacheEntry) {
				deadline := ent.createdAt.Add(interval)
				if t.After(deadline) {
					c.Delete(key)
				}
			})
		}
	}()
}
