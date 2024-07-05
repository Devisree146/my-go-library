package test

import (
	"fmt"
	"testing"
	"unified/in_memory"
)

func BenchmarkSet(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "benchmarkKey"
	value := 42

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("%s%d", key, i), value)
	}
}

func BenchmarkGet(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "benchmarkKey"
	value := 42
	cache.Set(key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(key)
	}
}

func BenchmarkDelete(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "benchmarkKey"
	value := 42
	cache.Set(key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(key)
	}
}

func BenchmarkExists(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "benchmarkKey"
	value := 42
	cache.Set(key, value)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Exists(key)
	}
}

func BenchmarkDeleteAll(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.DeleteAll()
	}
}

func BenchmarkNegativeGet(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "nonExistentKey"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(key)
	}
}

func BenchmarkNegativeDelete(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	key := "nonExistentKey"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Delete(key)
	}
}

func BenchmarkGetAll(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000, 0)

	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAllKeys()
	}
}
