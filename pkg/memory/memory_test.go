package memory

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_SetAndGet(t *testing.T) {
	store := NewStorage[string, int]()
	store.Set("key1", 42)

	value, ok := store.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 42, value)
}

func TestStorage_GetNonExistentKey(t *testing.T) {
	store := NewStorage[string, int]()

	value, ok := store.Get("missing")
	assert.False(t, ok)
	assert.Zero(t, value)
}

func TestStorage_ConcurrentAccess(t *testing.T) {
	store := NewStorage[int, int]()
	var wg sync.WaitGroup
	n := 1000

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store.Set(i, i*2)
		}(i)
	}

	wg.Wait()

	for i := 0; i < n; i++ {
		value, ok := store.Get(i)
		assert.True(t, ok)
		assert.Equal(t, i*2, value, i)
	}
}
