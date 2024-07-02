package test

import (
	"reflect"
	"sort"
	"testing"
	"time"
	"unified/in_memory"
)

func TestInMemoryCache(t *testing.T) {
	cache := in_memory.NewInMemoryCache(2)

	// Test Set and Get
	cache.Set("key1", "value1", 1*time.Second)
	value, err := cache.Get("key1")
	if err != nil || value != "value1" {
		t.Fatalf("expected 'value1', got '%v'", value)
	}

	// Test TTL expiration
	time.Sleep(2 * time.Second)
	_, err = cache.Get("key1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Test LRU eviction
	cache.Set("key2", "value2", 1*time.Second)
	cache.Set("key3", "value3", 1*time.Second)
	_, err = cache.Get("key1")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Test GetAllKeys
	keys := cache.GetAllKeys()
	expectedKeys := []string{"key2", "key3"}
	sort.Strings(keys)
	sort.Strings(expectedKeys)
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Fatalf("expected keys '%v', got '%v'", expectedKeys, keys)
	}

	// Test Delete
	cache.Delete("key2")
	_, err = cache.Get("key2")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Test DeleteAll
	cache.Set("key4", "value4", 1*time.Second)
	cache.Set("key5", "value5", 1*time.Second)
	cache.DeleteAll()
	keys = cache.GetAllKeys()
	if len(keys) != 0 {
		t.Fatalf("expected cache to be empty, got '%v' keys", len(keys))
	}
}
