package inmemory

import (
	"testing"

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

//func TestThousandGoroutines(t *testing.T) {
//
//}
