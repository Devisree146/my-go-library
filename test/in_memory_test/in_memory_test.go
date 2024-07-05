package test

import (
	"testing"
	"time"

	"unified/in_memory"
)

func TestSetAndGet(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Positive test case: Set and Get a value
	key := "testKey"
	value := 42
	err := cache.Set(key, value)
	if err != nil {
		t.Errorf("Error setting value: %v", err)
	}

	result, err := cache.Get(key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if result != value {
		t.Errorf("Expected %d but got %d", value, result)
	}

	// Negative test case: Get non-existent key
	_, err = cache.Get("nonExistentKey")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
}

func TestDelete(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Set a value
	key := "testKey"
	value := 42
	cache.Set(key, value)

	// Positive test case: Delete existing key
	err := cache.Delete(key)
	if err != nil {
		t.Errorf("Error deleting key: %v", err)
	}

	// Verify the key is deleted
	_, err = cache.Get(key)
	if err == nil {
		t.Error("Expected error for deleted key")
	}

	// Negative test case: Delete non-existent key
	err = cache.Delete("nonExistentKey")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}
}

func TestExists(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Set a value
	key := "testKey"
	value := 42
	cache.Set(key, value)

	// Positive test case: Check if key exists
	if !cache.Exists(key) {
		t.Errorf("Expected key %s to exist", key)
	}

	// Negative test case: Check if non-existent key exists
	if cache.Exists("nonExistentKey") {
		t.Error("Expected non-existent key to not exist")
	}
}

func TestCleanupExpiredEntries(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 1*time.Second)

	// Set a value with TTL
	key := "testKey"
	value := 42
	cache.Set(key, value)

	// Wait for TTL to expire
	time.Sleep(2 * time.Second)

	// Verify the key is deleted due to expiration
	_, err := cache.Get(key)
	if err == nil {
		t.Error("Expected error for expired key")
	}
}

func TestDeleteAll(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Set some values
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)

	// Delete all keys
	cache.DeleteAll()

	// Verify all keys are deleted
	keys := cache.GetAllKeys()
	if len(keys) != 0 {
		t.Errorf("Expected all keys to be deleted, found %d keys", len(keys))
	}
}

func TestLRUEviction(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Set 3 values
	cache.Set("key1", 1)
	cache.Set("key2", 2)
	cache.Set("key3", 3)

	// Add another value to trigger LRU eviction
	cache.Set("key4", 4)

	// Verify the first key was evicted
	_, err := cache.Get("key1")
	if err == nil {
		t.Error("Expected key1 to be evicted")
	}

	// Verify other keys still exist
	if _, err := cache.Get("key2"); err != nil {
		t.Error("Expected key2 to exist")
	}
	if _, err := cache.Get("key3"); err != nil {
		t.Error("Expected key3 to exist")
	}
	if _, err := cache.Get("key4"); err != nil {
		t.Error("Expected key4 to exist")
	}
}

func TestSetAndGetWithTTL(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 2*time.Second)

	// Set a value with TTL
	key := "testKey"
	value := 42
	cache.Set(key, value)

	// Verify the key exists before TTL expiry
	result, err := cache.Get(key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if result != value {
		t.Errorf("Expected %d but got %d", value, result)
	}

	// Wait for TTL to expire
	time.Sleep(3 * time.Second)

	// Verify the key is deleted due to expiration
	_, err = cache.Get(key)
	if err == nil {
		t.Error("Expected error for expired key")
	}
}

func TestSetUpdate(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Set a value
	key := "testKey"
	value := 42
	cache.Set(key, value)

	// Update the value
	newValue := 84
	cache.Set(key, newValue)

	// Verify the updated value
	result, err := cache.Get(key)
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}
	if result != newValue {
		t.Errorf("Expected %d but got %d", newValue, result)
	}
}

func TestNegativeSetWithEmptyKey(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Negative test case: Set with empty key
	err := cache.Set("", 42)
	if err == nil {
		t.Error("Expected error for empty key")
	} else {
		t.Logf("Expected error for empty key: %v", err)
	}
}

func TestNegativeSetWithNilValue(t *testing.T) {
	cache := in_memory.NewInMemoryCache(3, 5*time.Second)

	// Negative test case: Set with nil value
	err := cache.Set("testKey", nil)
	if err == nil {
		t.Error("Expected error for nil value")
	} else {
		t.Logf("Expected error for nil value: %v", err)
	}
}
