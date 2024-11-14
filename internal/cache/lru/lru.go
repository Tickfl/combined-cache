package lru

import (
	"container/list"
)

type CacheItem struct {
	key   string
	value interface{}
}

type Cache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *Cache) Add(key string, value interface{}) {
	if element, exists := c.cache[key]; exists {
		element.Value.(*CacheItem).value = value
		c.list.MoveToFront(element)
		return
	}

	if c.list.Len() >= c.capacity {
		c.removeLast()
	}

	newItem := &CacheItem{key: key, value: value}
	listElement := c.list.PushFront(newItem)
	c.cache[key] = listElement
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if element, exists := c.cache[key]; exists {
		c.list.MoveToFront(element)
		return element.Value.(*CacheItem).value, true
	}

	return nil, false
}

func (c *Cache) removeLast() {
	item := c.list.Back()
	if item != nil {
		c.list.Remove(item)
		delete(c.cache, item.Value.(*CacheItem).key)
	}
}
