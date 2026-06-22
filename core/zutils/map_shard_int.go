package zutils

import (
	"sync"
)

const shardCount = 64 // 必须 2^n

type shard[V any] struct {
	mu   sync.RWMutex
	data map[uint64]V
}

type ConcurrentMap[V any] struct {
	shards [shardCount]*shard[V]
}

func NewMap[V any]() *ConcurrentMap[V] {
	m := ConcurrentMap[V]{}

	for i := 0; i < shardCount; i++ {
		m.shards[i] = &shard[V]{
			data: make(map[uint64]V),
		}
	}

	return &m
}

func (m *ConcurrentMap[V]) getShard(key uint64) *shard[V] {
	return m.shards[key&(shardCount-1)]
}

func (m *ConcurrentMap[V]) Set(key uint64, value V) {
	s := m.getShard(key)

	s.mu.Lock()
	s.data[key] = value
	s.mu.Unlock()
}

func (m *ConcurrentMap[V]) Get(key uint64) (V, bool) {
	s := m.getShard(key)

	s.mu.RLock()
	v, ok := s.data[key]
	s.mu.RUnlock()

	return v, ok
}

func (m *ConcurrentMap[V]) Delete(key uint64) {
	s := m.getShard(key)

	s.mu.Lock()
	delete(s.data, key)
	s.mu.Unlock()
}

func (m *ConcurrentMap[V]) LoadOrStore(key uint64, value V) (V, bool) {
	s := m.getShard(key)

	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.data[key]; ok {
		return v, true
	}

	s.data[key] = value

	return value, false
}

func (m *ConcurrentMap[V]) Len() int {
	total := 0

	for i := 0; i < shardCount; i++ {
		s := m.shards[i]

		s.mu.RLock()
		total += len(s.data)
		s.mu.RUnlock()
	}

	return total
}

func (m *ConcurrentMap[V]) Range(fn func(key uint64, value V) bool) {
	for i := 0; i < shardCount; i++ {
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

func (m *ConcurrentMap[V]) All() map[uint64]V {
	data := make(map[uint64]V, 1000)

	for i := 0; i < shardCount; i++ {
		s := m.shards[i]

		s.mu.RLock()
		for k, v := range s.data {
			data[k] = v
		}
		s.mu.RUnlock()
	}

	return data
}
