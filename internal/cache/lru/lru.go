package lru

import (
	"container/list"
	"time"
)

type CacheItem struct {
	key   string
	value interface{}
	time  time.Time
}

type Cache struct {
	capacity int
	ttl      time.Duration
	cache    map[string]*list.Element
	list     *list.List
}

func NewCache(capacity int, ttl time.Duration) *Cache {
	return &Cache{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *Cache) Add(key string, value interface{}) {
	if element, exists := c.cache[key]; exists {
		element.Value.(*CacheItem).value = value
		element.Value.(*CacheItem).time = time.Now()
		c.list.MoveToFront(element)
		return
	}

	if c.list.Len() >= c.capacity {
		c.removeLast()
	}

	newItem := &CacheItem{
		key:   key,
		value: value,
		time:  time.Now(),
	}
	listElement := c.list.PushFront(newItem)
	c.cache[key] = listElement
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if element, exists := c.cache[key]; exists {
		if time.Since(element.Value.(*CacheItem).time) < c.ttl {
			c.list.MoveToFront(element)
			return element.Value.(*CacheItem).value, true
		}
		c.removeElement(element)
	}

	return nil, false
}

func (c *Cache) removeLast() {
	element := c.list.Back()
	if element != nil {
		c.removeElement(element)
	}
}

func (c *Cache) removeElement(element *list.Element) {
	c.list.Remove(element)
	delete(c.cache, element.Value.(*CacheItem).key)
}
