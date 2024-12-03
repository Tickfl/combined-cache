package lru

import (
	"testing"
	"time"
)

func TestCache_AddAndGet(t *testing.T) {
	cache := NewCache(2, 5*time.Second)
	defer cache.StopCleanUp()

	cache.Add("key1", 1)
	cache.Add("key2", 2)

	value, exists := cache.Get("key1")
	if !exists || value != 1 {
		t.Errorf("expected 1, got %v (exists: %v)", value, exists)
	}

	_, exists = cache.Get("key3")
	if exists {
		t.Errorf("expected not exists for key3")
	}
}

func TestCache_TTLExpiration(t *testing.T) {
	cache := NewCache(2, time.Second)
	defer cache.StopCleanUp()

	cache.Add("key1", 1)
	time.Sleep(2 * time.Second)

	_, exists := cache.Get("key1")
	if exists {
		t.Errorf("expected not exists for key1")
	}
}

func TestCache_CapacityLimit(t *testing.T) {
	cache := NewCache(2, 5*time.Second)
	defer cache.StopCleanUp()

	cache.Add("key1", 1)
	cache.Add("key2", 2)
	cache.Add("key3", 3)

	_, exists := cache.Get("key1")
	if exists {
		t.Errorf("expected not exists for key1")
	}
}
