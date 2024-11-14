package main

import (
	"combined-cache/internal/cache/lru"
	"fmt"
)

func main() {
	cache := lru.NewCache(3)
	cache.Add("key1", 1)
	value, exists := cache.Get("key1")

	fmt.Println(value, exists)
}
