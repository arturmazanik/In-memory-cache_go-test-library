package inmemlib

import (
	"sync"
	"time"
)

type cacheItem struct {
	value       any
	whenExpired time.Duration
}

type Cache struct {
	mu       sync.RWMutex
	lifeTime time.Duration
	storage  map[string]cacheItem
}

func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		lifeTime: ttl,
		storage:  make(map[string]cacheItem),
	}
}

func (c *Cache) Set(key string, value string, lifeTime time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.storage[key] = cacheItem{
		value:       value,
		whenExpired: lifeTime,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.storage[key]
	c.mu.RUnlock()

	if !exists {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.storage, key)
	c.mu.Unlock()
}
