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

func (c *Cache) Set(key string, value string, lifeTime time.Duration) error {
	c.mu.Lock()

	c.storage[key] = cacheItem{
		value:       value,
		whenExpired: lifeTime,
	}

	c.mu.Unlock()

	return nil
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()

	item, found := c.storage[key]

	if !found {
		return nil, false
	}

	c.mu.Unlock()

	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.storage, key)
	c.mu.Unlock()
}
