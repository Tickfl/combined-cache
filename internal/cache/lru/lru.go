package lru

import (
	"log/slog"
	"sync"
	"time"
)

type CacheItem struct {
	key        string
	value      interface{}
	time       time.Time
	prev, next *CacheItem
}

type Cache struct {
	capacity    int
	ttl         time.Duration
	cache       map[string]*CacheItem
	first, last *CacheItem
	mu          sync.Mutex
	logger      *slog.Logger
}

func NewCache(capacity int, ttl time.Duration) *Cache {
	cache := &Cache{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[string]*CacheItem, capacity),
		logger:   slog.Default(),
	}

	go cache.startCleanUp()

	cache.logger.Info("NewCache", slog.Int("capacity", capacity), slog.Duration("ttl", ttl))

	return cache
}

func (c *Cache) Add(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.cache[key]; exists {
		element.value = value
		element.time = time.Now()
		c.moveToFront(element)
		c.logger.Info("Update element", slog.String("key", key), slog.Any("value", value))
		return
	}

	if len(c.cache) >= c.capacity {
		c.logger.Warn("Remove last element")
		c.removeLast()
	}

	newItem := &CacheItem{
		key:   key,
		value: value,
		time:  time.Now(),
	}
	c.cache[key] = newItem
	c.addToFront(newItem)
	c.logger.Info("Add element", slog.String("key", key), slog.Any("value", value))
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.cache[key]; exists {
		if time.Since(element.time) < c.ttl {
			c.moveToFront(element)
			return element.value, true
		}
		c.removeElement(element)
	}

	return nil, false
}

func (c *Cache) addToFront(element *CacheItem) {
	if c.first != nil {
		c.first.prev = element
	}
	if c.last == nil {
		c.last = element
	}
	element.next = c.first
	element.prev = nil
	c.first = element
}

func (c *Cache) moveToFront(element *CacheItem) {
	if c.first == element {
		return
	}
	if c.last == element {
		c.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}
	if c.first != nil {
		c.first.prev = element
	}

	element.next = c.first
	element.prev = nil
	c.first = element
}

func (c *Cache) removeLast() {
	if c.last != nil {
		c.removeElement(c.last)
	}
}

func (c *Cache) removeElement(element *CacheItem) {
	if c.first == element {
		c.first = element.next
	}
	if c.last == element {
		c.last = element.prev
	}
	if element.prev != nil {
		element.prev.next = element.next
	}
	if element.next != nil {
		element.next.prev = element.prev
	}

	c.logger.Info("Remove element", slog.String("key", element.key), slog.Any("value", element.value))
	delete(c.cache, element.key)
}

func (c *Cache) startCleanUp() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.CleanUp()
		}
	}
}

func (c *Cache) CleanUp() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, element := range c.cache {
		if time.Since(element.time) >= c.ttl {
			c.removeElement(element)
		}
	}
}
