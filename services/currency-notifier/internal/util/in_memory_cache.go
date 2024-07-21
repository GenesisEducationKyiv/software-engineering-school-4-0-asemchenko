package util

import (
	"sync"
	"time"
)

type InMemoryCache[T any] struct {
	updatedAt   time.Time
	cachedValue T
	cacheMutex  sync.Mutex
	ttl         time.Duration
}

func NewInMemoryCache[T any](ttl time.Duration) *InMemoryCache[T] {
	return &InMemoryCache[T]{
		ttl: ttl,
	}
}

func (c *InMemoryCache[T]) Get() (T, bool) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	if time.Since(c.updatedAt) < c.ttl {
		return c.cachedValue, true
	}
	return c.cachedValue, false
}

func (c *InMemoryCache[T]) Set(value T) {
	c.cacheMutex.Lock()
	defer c.cacheMutex.Unlock()
	c.cachedValue = value
	c.updatedAt = time.Now()
}
