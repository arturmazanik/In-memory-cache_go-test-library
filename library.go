package inmemlib

import (
	"fmt"
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

// Worker Pool
type CacheTask struct {
	Key   string
	Value string
	TTL   time.Duration
}

type WorkerPool struct {
	Tasks      chan CacheTask
	cache      *Cache
	wg         sync.WaitGroup
	numWorkers int
}

// Создать новый пул
func NewWorkerPool(numWorkers int, cache *Cache) *WorkerPool {
	return &WorkerPool{
		Tasks:      make(chan CacheTask, 100),
		cache:      cache,
		numWorkers: numWorkers,
	}
}

// Запустить воркеров
func (wp *WorkerPool) Run() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Worker функция
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for task := range wp.Tasks {
		// Здесь добавляем элемент в кеш
		wp.cache.Set(task.Key, task.Value, task.TTL)
		fmt.Printf("Worker %d set key=%s\n", id, task.Key)
	}
}

// Ожидание завершения всех воркеров
func (wp *WorkerPool) Wait() {
	close(wp.Tasks)
	wp.wg.Wait()
}
