package memory

import "sync"

type Storage[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewStorage[K comparable, V any]() *Storage[K, V] {
	return &Storage[K, V]{
		data: make(map[K]V),
	}
}

func (s *Storage[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]

	return v, ok
}

func (s *Storage[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
