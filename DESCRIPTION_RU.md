Тестовое задание


Дано:
	InMemoryCache - потоко-безопасная реализация Key-Value кэша, хранящая данные в оперативной памяти
Задача:
	1. Реализовать метод GetOrSet, предоставив следующие гарантии:
		- Значение каждого ключа будет вычислено ровно 1 раз
		- Конкурентные обращения к существующим ключам не блокируют друг друга
	2. Покрыть его тестами, проверить метод 1000+ горутинами
*/

// ----------------------------------------------

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

// ----------------------------------------------

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn, сохраняет его в кэш и возвращает это значение.
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
    // TODO
}

