package inmemory

import (
	"sync"

	tsc "github.com/morozovcookie/threadsafecache"
)

type Cache struct {
	dataMutex sync.RWMutex
	data      map[tsc.Key]tsc.Value
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[tsc.Key]tsc.Value),
	}
}

func (cache *Cache) Get(key tsc.Key) (tsc.Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	val, ok := cache.data[key]

	return val, ok
}

func (cache *Cache) GetOrSet(key tsc.Key, valueFn func() tsc.Value) tsc.Value {
	val, ok := cache.Get(key)
	if ok {
		return val
	}

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	val = valueFn()
	cache.data[key] = val

	return val
}
