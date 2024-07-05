package test

import (
	"testing"
	"unified/redis_cache"
)

func TestSetAndGet(t *testing.T) {
	cache := redis_cache.NewRedisCache("localhost:6379", "", 0, 10)

	value := 12345

	err := cache.Set("test_key", value, redis_cache.StandardTTL)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	got, err := cache.Get("test_key")
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if got != value {
		t.Errorf("Expected %v, got %v", value, got)
	}
}

func TestDelete(t *testing.T) {
	cache := redis_cache.NewRedisCache("localhost:6379", "", 0, 10)

	value := 12345

	err := cache.Set("test_key", value, redis_cache.StandardTTL)
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	err = cache.Delete("test_key")
	if err != nil {
		t.Fatalf("Failed to delete value: %v", err)
	}

	_, err = cache.Get("test_key")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
