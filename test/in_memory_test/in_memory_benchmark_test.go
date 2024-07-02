package test

import (
	"fmt"
	"testing"
	"time"
	"unified/in_memory"
)

func BenchmarkInMemoryCacheSet(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 5*time.Minute)
	}
}

func BenchmarkInMemoryCacheGet(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key%d", i))
	}
}

func BenchmarkInMemoryCacheDelete(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	for i := 0; i < b.N; i++ {
		cache.Delete(fmt.Sprintf("key%d", i))
	}
}

func BenchmarkInMemoryCacheSetParallel(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Set(fmt.Sprintf("key%d", i), "value", 5*time.Minute)
			i++
		}
	})
}

func BenchmarkInMemoryCacheGetParallel(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			cache.Get(fmt.Sprintf("key%d", i))
			i++
		}
	})
}

func BenchmarkInMemoryCacheDeleteAll(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 5*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.DeleteAll()
		for i := 0; i < 1000; i++ {
			cache.Set(fmt.Sprintf("key%d", i), "value", 5*time.Minute)
		}
	}
}

func BenchmarkInMemoryCacheGetAll(b *testing.B) {
	cache := in_memory.NewInMemoryCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), "value", 5*time.Minute)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.GetAllKeys() // Use GetAllKeys instead of Keys
	}
}
