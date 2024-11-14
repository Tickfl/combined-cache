package main

import (
	"combined-cache/internal/cache/lru"
	"fmt"
	"time"
)

func main() {
	cache := lru.NewCache(3, time.Minute)
	cache.Add("key1", 1)

	value, exists := cache.Get("key1")
	fmt.Println(value, exists)

	time.Sleep(time.Minute)

	value, exists = cache.Get("key1")
	fmt.Println(value, exists)
}
