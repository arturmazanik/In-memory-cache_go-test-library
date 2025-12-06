package inmemlib

import (
	"fmt"
	"sync"
	"time"
)

// Элемент кеша
type cacheItem struct {
	value      interface{}
	expiration int64
}

// Кеш с защитой для горутин
type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

// Новый кеш
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

// Установить значение с TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	exp := int64(0)
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = cacheItem{
		value:      value,
		expiration: exp,
	}
	c.mu.Unlock()
}

// Получить значение
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		c.Delete(key)
		return nil, false
	}

	return it.value, true
}

// Удалить ключ
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// Очистить кеш
func (c *Cache) Clear() {
	c.mu.Lock()
	c.items = make(map[string]cacheItem)
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
