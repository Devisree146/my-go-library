package in_memory

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Entry represents a cache entry with key, value, and TTL.
type Entry struct {
	Key   string
	Value interface{}
	TTL   time.Time
}

// InMemoryCache represents an in-memory cache with LRU eviction.
type InMemoryCache struct {
	maxSize int
	cache   map[string]*list.Element
	lruList *list.List
	lock    sync.Mutex
}

// NewInMemoryCache initializes a new cache with a given maximum size.
func NewInMemoryCache(maxSize int) *InMemoryCache {
	return &InMemoryCache{
		maxSize: maxSize,
		cache:   make(map[string]*list.Element),
		lruList: list.New(),
	}
}

// Set adds or updates a key-value pair in the cache and handles LRU eviction.
func (c *InMemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// If the key already exists, update the value and TTL, and move it to the front.
	if element, exists := c.cache[key]; exists {
		c.lruList.MoveToFront(element)
		element.Value.(*Entry).Value = value
		element.Value.(*Entry).TTL = time.Now().Add(ttl)
		return nil
	}

	// If the cache is at its maximum size, evict the least recently used element.
	if len(c.cache) >= c.maxSize {
		c.evict()
	}

	// Add the new key-value pair to the cache.
	newEntry := &Entry{
		Key:   key,
		Value: value,
		TTL:   time.Now().Add(ttl),
	}
	element := c.lruList.PushFront(newEntry)
	c.cache[key] = element

	return nil
}

// Get fetches the value from the cache and moves the entry to the front of the LRU list.
func (c *InMemoryCache) Get(key string) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the key exists in the cache.
	if element, exists := c.cache[key]; exists {
		// Check if the entry has expired.
		if element.Value.(*Entry).TTL.After(time.Now()) {
			c.lruList.MoveToFront(element)
			return element.Value.(*Entry).Value, nil
		}
		// If the entry has expired, remove it.
		c.removeElement(element)
	}

	return nil, fmt.Errorf("key not found")
}

// Delete removes an entry from the cache.
func (c *InMemoryCache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, exists := c.cache[key]; exists {
		c.removeElement(element)
		return nil
	}

	return fmt.Errorf("key not found")
}

// DeleteAll removes all entries from the cache.
func (c *InMemoryCache) DeleteAll() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.lruList.Init()
	c.cache = make(map[string]*list.Element)
}

// evict removes the least recently used entry from the cache.
func (c *InMemoryCache) evict() {
	element := c.lruList.Back()
	if element != nil {
		c.removeElement(element)
	}
}

// removeElement removes a specific element from the linked list and hash map.
func (c *InMemoryCache) removeElement(element *list.Element) {
	c.lruList.Remove(element)
	delete(c.cache, element.Value.(*Entry).Key)
}

// Exists checks if a key is present in the cache.
func (c *InMemoryCache) Exists(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, exists := c.cache[key]
	return exists
}

// GetAllKeys returns a slice of all keys in the cache.
func (c *InMemoryCache) GetAllKeys() []string {
	c.lock.Lock()
	defer c.lock.Unlock()

	keys := make([]string, 0, len(c.cache))
	for key := range c.cache {
		keys = append(keys, key)
	}

	return keys
}
