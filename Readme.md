# Go In-Memory-Library (lite) #
training project

## Setup ##

```
go get github.com/arturmazanik/in-memory-cache_go-test-library 
```

## Example ## 

```
func main() {
	cache := inmemlib.NewCache()
	wp := inmemlib.NewWorkerPool(3, cache)
	wp.Run()

	// Отправляем задачи
	for i := 1; i <= 10; i++ {
		task := inmemlib.CacheTask{
			Key:   fmt.Sprintf("key-%d", i),
			Value: strconv.Itoa(i * 10),
			TTL:   time.Second * 5,
		}
		wp.Tasks <- task
	}

	// Завершаем работу воркеров
	wp.Wait()

	// Проверка
	for i := 1; i <= 10; i++ {
		if val, ok := cache.Get(fmt.Sprintf("key-%d", i)); ok {
			fmt.Println("Cache value:", val)
		} else {
			fmt.Println("Key expired or missing:", i)
		}
	}
}
```
