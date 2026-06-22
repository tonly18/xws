package zutils

import (
	"sync"

	"github.com/cespare/xxhash/v2"
)

const shardCountStr = 64 // 必须 2^n

type shardStr[V any] struct {
	mu   sync.RWMutex
	data map[string]V
}

type ConcurrentMapStr[V any] struct {
	shards [shardCountStr]*shardStr[V]
}

func NewMapStr[V any]() *ConcurrentMapStr[V] {
	m := &ConcurrentMapStr[V]{}

	for i := 0; i < shardCountStr; i++ {
		m.shards[i] = &shardStr[V]{
			data: make(map[string]V),
		}
	}

	return m
}

func (m *ConcurrentMapStr[V]) getShard(key string) *shardStr[V] {
	return m.shards[xxhash.Sum64String(key)&(shardCountStr-1)]
}

func (m *ConcurrentMapStr[V]) Set(key string, value V) {
	s := m.getShard(key)

	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
}

func (m *ConcurrentMapStr[V]) Get(key string) (V, bool) {
	s := m.getShard(key)

	s.mu.RLock()
	v, ok := s.data[key]
	s.mu.RUnlock()

	return v, ok
}

func (m *ConcurrentMapStr[V]) Delete(key string) {
	s := m.getShard(key)

	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
}

func (m *ConcurrentMapStr[V]) LoadOrStore(key string, value V) (V, bool) {
	s := m.getShard(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.data[key]; ok {
		return v, true
	}

	s.data[key] = value

	return value, false
}

func (m *ConcurrentMapStr[V]) Len() int {
	total := 0

	for i := 0; i < shardCountStr; i++ {
		s := m.shards[i]

		s.mu.RLock()
		total += len(s.data)
		s.mu.RUnlock()
	}

	return total
}

func (m *ConcurrentMapStr[V]) Range(fn func(key string, value V) bool) {
	for i := 0; i < shardCountStr; i++ {
		s := m.shards[i]

		s.mu.RLock()
		for k, v := range s.data {
			if !fn(k, v) {
				s.mu.RUnlock()
				return
			}
		}
		s.mu.RUnlock()
	}
}

func (m *ConcurrentMapStr[V]) All() map[string]V {
	data := make(map[string]V, 1000)

	for i := 0; i < shardCountStr; i++ {
		s := m.shards[i]

		s.mu.RLock()
		for k, v := range s.data {
			data[k] = v
		}
		s.mu.RUnlock()
	}

	return data
}
