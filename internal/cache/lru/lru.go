package lru

import (
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
}

func NewCache(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[string]*CacheItem, capacity),
	}
}

func (c *Cache) Add(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.cache[key]; exists {
		element.value = value
		element.time = time.Now()
		c.moveToFront(element)
		return
	}

	if len(c.cache) >= c.capacity {
		c.removeLast()
	}

	newItem := &CacheItem{
		key:   key,
		value: value,
		time:  time.Now(),
	}
	c.cache[key] = newItem
	c.addToFront(newItem)
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

	delete(c.cache, element.key)
}
