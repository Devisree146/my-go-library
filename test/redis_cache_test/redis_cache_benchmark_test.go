package redis_cache

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var benchmarkCtx = context.Background()

func setupBenchmarkCache(maxSize int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	client.FlushDB(benchmarkCtx) // Start with empty Redis
	return NewRedisCache("localhost:6379", "", 0, maxSize)
}

func BenchmarkSet(b *testing.B) {
	cache := setupBenchmarkCache(10000) // Adjust maxSize for your test

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i)
		value := map[string]interface{}{
			"value": i,
		}
		err := cache.Set(key, value, 10*time.Second)
		if err != nil {
			b.Fatalf("Error setting key %s: %v", key, err)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	cache := setupBenchmarkCache(10000) // Adjust maxSize for your test

	// Populate cache with initial data
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := map[string]interface{}{
			"value": i,
		}
		err := cache.Set(key, value, 10*time.Second)
		if err != nil {
			b.Fatalf("Error setting key %s: %v", key, err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%10000) // Reuse existing keys for benchmark
		_, err := cache.Get(key)
		if err != nil {
			b.Fatalf("Error getting key %s: %v", key, err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cache := setupBenchmarkCache(10000) // Adjust maxSize for your test

	// Populate cache with initial data
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key%d", i)
		value := map[string]interface{}{
			"value": i,
		}
		err := cache.Set(key, value, 10*time.Second)
		if err != nil {
			b.Fatalf("Error setting key %s: %v", key, err)
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key%d", i%10000) // Reuse existing keys for benchmark
		err := cache.Delete(key)
		if err != nil {
			b.Fatalf("Error deleting key %s: %v", key, err)
		}
	}
}
