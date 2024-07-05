package test

import (
	"testing"
	"unified/redis_cache"
)

func BenchmarkSet(b *testing.B) {
	cache := redis_cache.NewRedisCache("localhost:6379", "", 0, 10)
	value := 12345

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("benchmark_key", value, redis_cache.StandardTTL)
	}
}

func BenchmarkGet(b *testing.B) {
	cache := redis_cache.NewRedisCache("localhost:6379", "", 0, 10)
	value := 12345
	cache.Set("benchmark_key", value, redis_cache.StandardTTL)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("benchmark_key")
	}
}
