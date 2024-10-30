package cache

import (
	"sync"
	"time"
)

// CacheItem represents a single item in the cache
type CacheItem struct {
	Value      interface{}
	Expiration int64
}

// Cache represents the cache with a map and a mutex
type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// Exported cache instance
var CacheInstance = NewCache()

// NewCache creates a new Cache
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache with an expiration time in seconds
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).Unix(),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || item.Expiration < time.Now().Unix() {
		return nil, false
	}
	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Cleanup removes expired items from the cache
func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now().Unix()
	for key, item := range c.items {
		if item.Expiration < now {
			delete(c.items, key)
		}
	}
}
