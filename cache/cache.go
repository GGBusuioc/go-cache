package cache

import (
	"errors"
	"sync"
)

var (
	NotFoundErr = errors.New("not found")
)

type Item struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

// Cache in-memory key-value store
type Cache struct {
	data map[string]int
	mu   sync.RWMutex
}

// NewCache creates and initialise a new cache instance
func NewCache() *Cache {
	data := make(map[string]int)
	return &Cache{
		data: data,
	}
}

// Add adds a key-value pair in the cache
func (c *Cache) Add(key string, value int) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

// Get retrieves the value associated with the given key from the cache
func (c *Cache) Get(key string) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.data[key]; ok {
		return val, nil
	}

	return 0, NotFoundErr
}

func (c *Cache) List() (map[string]int, error) {
	return c.data, nil
}

func (c *Cache) Update(key string, value int) error {
	if _, ok := c.data[key]; ok {
		c.data[key] = value
		return nil
	}

	return NotFoundErr
}

// Remove removes all key-value pairs from the cache
func (c *Cache) Remove(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
