package tools

import "sync"

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		mutex:  new(sync.RWMutex),
		values: make(map[string]*T, 0),
	}
}

type Cache[T any] struct {
	mutex  *sync.RWMutex
	values map[string]*T
}

func (c *Cache[T]) Get(key string) *T {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	value, ok := c.values[key]
	if !ok {
		return nil
	}

	return value
}

func (c *Cache[T]) Set(key string, value *T) {
	if value == nil {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.values[key] = value
}

func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.values, key)
}
