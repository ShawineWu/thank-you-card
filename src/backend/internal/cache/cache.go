package cache

import (
	"sync"
	"time"
)

// CacheItem represents a cached value with expiration
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// InMemoryCache provides simple in-memory caching with TTL
type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	cache := &InMemoryCache{
		items: make(map[string]CacheItem),
	}

	// Start cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves a value from cache
func (c *InMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.Expiration) {
		return nil, false
	}

	return item.Value, true
}

// Set stores a value in cache with TTL
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from cache
func (c *InMemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// InvalidateAll clears all cached items
func (c *InMemoryCache) InvalidateAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]CacheItem)
}

// cleanupExpired removes expired items periodically
func (c *InMemoryCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.After(item.Expiration) {
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}
