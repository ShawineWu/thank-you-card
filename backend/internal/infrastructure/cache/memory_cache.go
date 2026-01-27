package cache

import (
	"sync"
	"time"
)

// CacheItem represents an item in the cache with expiration
type CacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// IsExpired checks if the cache item has expired
func (item CacheItem) IsExpired() bool {
	return time.Now().After(item.Expiration)
}

// MemoryCache is an in-memory cache implementation
type MemoryCache struct {
	items map[string]CacheItem
	mutex sync.RWMutex
	defaultTTL time.Duration
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache(defaultTTL time.Duration) *MemoryCache {
	cache := &MemoryCache{
		items:      make(map[string]CacheItem),
		defaultTTL: defaultTTL,
	}
	
	// Start cleanup goroutine
	go cache.startCleanup()
	
	return cache
}

// Set stores a value in the cache with default TTL
func (c *MemoryCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value in the cache with custom TTL
func (c *MemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
}

// Get retrieves a value from the cache
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return nil, false
	}
	
	if item.IsExpired() {
		// Item expired, remove it
		delete(c.items, key)
		return nil, false
	}
	
	return item.Value, true
}

// Delete removes a value from the cache
func (c *MemoryCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.items, key)
}

// Clear removes all items from the cache
func (c *MemoryCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]CacheItem)
}

// Size returns the number of items in the cache
func (c *MemoryCache) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	return len(c.items)
}

// Keys returns all keys in the cache
func (c *MemoryCache) Keys() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	
	return keys
}

// startCleanup starts a goroutine to periodically clean up expired items
func (c *MemoryCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute) // Cleanup every 5 minutes
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.cleanupExpired()
		}
	}
}

// cleanupExpired removes expired items from the cache
func (c *MemoryCache) cleanupExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	now := time.Now()
	for key, item := range c.items {
		if now.After(item.Expiration) {
			delete(c.items, key)
		}
	}
}

// Cache interface for dependency injection
type Cache interface {
	Set(key string, value interface{})
	SetWithTTL(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
	Clear()
	Size() int
	Keys() []string
}

// Ensure MemoryCache implements Cache interface
var _ Cache = (*MemoryCache)(nil)

// CacheKeys contains commonly used cache key patterns
type CacheKeys struct {
	EmployeePrefix     string
	CompanyValuesKey   string
	Top10EmployeesKey  string
	StatisticsPrefix   string
}

// DefaultCacheKeys returns default cache key patterns
func DefaultCacheKeys() CacheKeys {
	return CacheKeys{
		EmployeePrefix:     "employee:",
		CompanyValuesKey:   "company_values:all",
		Top10EmployeesKey:  "top10_employees",
		StatisticsPrefix:   "stats:",
	}
}

// Helper functions for common cache operations

// GetEmployeeKey generates a cache key for employee data
func GetEmployeeKey(employeeID string) string {
	return DefaultCacheKeys().EmployeePrefix + employeeID
}

// GetStatisticsKey generates a cache key for employee statistics
func GetStatisticsKey(employeeID string) string {
	return DefaultCacheKeys().StatisticsPrefix + employeeID
}