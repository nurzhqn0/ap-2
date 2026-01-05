package store

import (
	"sync"
)

// Repository generic
type Repository[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// NewRepository new generic repository
func NewRepository[K comparable, V any]() *Repository[K, V] {
	return &Repository[K, V]{
		data: make(map[K]V),
	}
}

// Set stores a value by key
func (r *Repository[K, V]) Set(key K, value V) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

// Get retrieves a value by key
func (r *Repository[K, V]) Get(key K) (V, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	val, ok := r.data[key]
	return val, ok
}

// GetAll returns all
func (r *Repository[K, V]) GetAll() []V {
	r.mu.RLock()
	defer r.mu.RUnlock()

	values := make([]V, 0, len(r.data))
	for _, v := range r.data {
		values = append(values, v)
	}
	return values
}

// Count returns the number of items
func (r *Repository[K, V]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.data)
}
