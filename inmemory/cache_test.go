package inmemory

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	tsc "github.com/morozovcookie/threadsafecache"
	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	tests := []struct {
		name string

		cache *Cache

		key tsc.Key

		expectedVal tsc.Value
		expectedOk  bool
	}{
		{
			name: "value exist",

			cache: &Cache{
				data: map[tsc.Key]tsc.Value{
					"test-key": "test-value",
				},
			},

			key: "test-key",

			expectedVal: "test-value",
			expectedOk:  true,
		},
		{
			name: "value does not exist",

			cache: NewCache(),

			key: "test-key",

			expectedVal: "",
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualVal, actualOk := tt.cache.Get(tt.key)

			assert.Equal(t, tt.expectedVal, actualVal)
			assert.Equal(t, tt.expectedOk, actualOk)
		})
	}
}

func TestCache_GetOrSet(t *testing.T) {
	tests := []struct {
		name string

		cache *Cache

		key     tsc.Key
		valueFn func() tsc.Value

		expectedVal tsc.Value
	}{
		{
			name: "value exist",

			cache: &Cache{
				data: map[tsc.Key]tsc.Value{
					"test-key": "test-value",
				},
			},

			key: "test-key",
			valueFn: func() tsc.Value {
				return "test-value"
			},

			expectedVal: "test-value",
		},
		{
			name: "value does not exist, but we created it and now it exist",

			cache: NewCache(),

			key: "test-key",
			valueFn: func() tsc.Value {
				return "test-value"
			},

			expectedVal: "test-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expectedVal, tt.cache.GetOrSet(tt.key, tt.valueFn))
		})
	}
}

func TestValueFnWillCallOnlyOne(t *testing.T) {
	var (
		cache = NewCache()

		key = "test-key"
		i   = 0
		fn  = func() tsc.Value {
			i++
			return "test-value"
		}

		expectedVal = "test-value"
	)

	assert.Equal(t, expectedVal, cache.GetOrSet(key, fn))
	assert.Equal(t, 1, i)

	assert.Equal(t, expectedVal, cache.GetOrSet(key, fn))
	assert.Equal(t, 1, i)
}

func goroutineCacheGetOrSet(
	ctx context.Context,
	t *testing.T,
	idx int,
	wg *sync.WaitGroup,
	cache *Cache,
	key tsc.Key,
	fn func() tsc.Value,
	expectedVal tsc.Value) {
	//t.Logf("starting goroutine #%d", idx)

	for {
		select {
		case <-time.After(time.Millisecond * 10):
			actualVal := cache.GetOrSet(key, fn)
			//t.Logf("[goroutine #%d] GetOrSet(%s, fn) = %s", idx, key, actualVal)
			assert.Equal(t, expectedVal, actualVal)
		case <-ctx.Done():
			//t.Logf("stopping goroutine #%d", idx)
			wg.Done()
			return
		}
	}
}

func goroutineCacheGetOrGetOrSet(
	ctx context.Context,
	t *testing.T,
	idx int,
	wg *sync.WaitGroup,
	cache *Cache,
	key tsc.Key,
	fn func() tsc.Value,
	expectedVal tsc.Value) {
	//t.Logf("starting goroutine #%d", idx)

	arr := []func(t *testing.T, cache *Cache, key tsc.Key, fn func() tsc.Value) tsc.Value{
		func(t *testing.T, cache *Cache, key tsc.Key, _ func() tsc.Value) tsc.Value {
			val, _ := cache.Get(key)

			//t.Logf("[goroutine #%d] Get(%s) = %s", idx, key, val)

			return val
		},
		func(t *testing.T, cache *Cache, key tsc.Key, fn func() tsc.Value) tsc.Value {
			val := cache.GetOrSet(key, fn)

			//t.Logf("[goroutine #%d] GetOrSet(%s, fn) = %s", idx, key, val)

			return val
		},
	}

	for {
		select {
		case <-time.After(time.Millisecond * 10):
			assert.Equal(t, expectedVal, arr[rand.Intn(2)](t, cache, key, fn))
		case <-ctx.Done():
			//t.Logf("stopping goroutine #%d", idx)
			wg.Done()
			return
		}
	}
}

func callGetOrSetConcurrent(t *testing.T, n int, testTime time.Duration) {
	var (
		cache = NewCache()

		key = "test-key"
		fn  = func() tsc.Value {
			return "test-value"
		}

		expectedVal = "test-value"

		wg = sync.WaitGroup{}
	)

	wg.Add(n)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < n; i++ {
		go goroutineCacheGetOrSet(ctx, t, i+1, &wg, cache, key, fn, expectedVal)
	}

	<-time.After(testTime)
	cancel()

	wg.Wait()
}

func callGetOrGetOrSetConcurrent(t *testing.T, n int, testTime time.Duration) {
	var (
		cache = NewCache()

		key = "test-key"
		fn  = func() tsc.Value {
			return "test-value"
		}

		expectedVal = "test-value"

		wg = sync.WaitGroup{}
	)

	wg.Add(n)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < n; i++ {
		go goroutineCacheGetOrGetOrSet(ctx, t, i+1, &wg, cache, key, fn, expectedVal)
	}

	<-time.After(testTime)
	cancel()

	wg.Wait()
}

func TestGetOrSetConcurrentCalls(t *testing.T) {
	callGetOrSetConcurrent(t, 2, time.Second)
}

func TestGetOrSet_1_000_Goroutines(t *testing.T) {
	callGetOrSetConcurrent(t, 1_000, time.Second)
}

func TestGetOrSet_10_000_Goroutines(t *testing.T) {
	callGetOrSetConcurrent(t, 10_000, time.Minute)
}

func TestGetOrSet_100_000_Goroutines(t *testing.T) {
	callGetOrSetConcurrent(t, 100_000, time.Minute)
}

func TestGetOrSet_1_000_000_Goroutines(t *testing.T) {
	callGetOrSetConcurrent(t, 1_000_000, time.Minute)
}

func TestGetOrGetOrSetConcurrentCalls(t *testing.T) {
	callGetOrGetOrSetConcurrent(t, 2, time.Second)
}

func TestGetOrGetOrSet_1_000_Goroutines(t *testing.T) {
	callGetOrGetOrSetConcurrent(t, 1_000, time.Second)
}

func TestGetOrGetOrSet_10_000_Goroutines(t *testing.T) {
	callGetOrGetOrSetConcurrent(t, 10_000, time.Minute)
}

func TestGetOrGetOrSet_100_000_Goroutines(t *testing.T) {
	callGetOrGetOrSetConcurrent(t, 100_000, time.Minute)
}

func TestGetOrGetOrSet_1_000_000_Goroutines(t *testing.T) {
	callGetOrGetOrSetConcurrent(t, 1_000_000, time.Minute)
}
