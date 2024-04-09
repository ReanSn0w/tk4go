package tools

import "sync"

func NewSafeSlice[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{data: make([]T, 0)}
}

type SafeSlice[T any] struct {
	data []T
	mu   sync.RWMutex
}

func (ss *SafeSlice[T]) Push(v ...T) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.data = append(ss.data, v...)
}

func (ss *SafeSlice[T]) Delete(i int) {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ss.data = append(ss.data[:i], ss.data[i+1:]...)
}

func (ss *SafeSlice[T]) Get(i int) T {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return ss.data[i]
}

func (ss *SafeSlice[T]) Len() int {
	ss.mu.RLock()
	defer ss.mu.RUnlock()
	return len(ss.data)
}
