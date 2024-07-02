package redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client  *redis.Client
	maxSize int
}

func NewRedisCache(addr, password string, db int, maxSize int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Cache{
		client:  client,
		maxSize: maxSize,
	}
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	// Manage cache size and eviction
	c.evictIfMaxSizeExceeded()

	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Cache) Delete(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) DeleteAll() error {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		err := c.client.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) GetAllKeys() []string {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil
	}

	sort.Strings(keys)
	return keys
}

func (c *Cache) evictIfMaxSizeExceeded() {
	keys := c.GetAllKeys()
	numKeys := len(keys)
	if numKeys > c.maxSize {
		// Calculate number of keys to delete to maintain max size
		keysToDelete := numKeys - c.maxSize
		for i := 0; i < keysToDelete; i++ {
			err := c.client.Del(context.Background(), keys[i]).Err()
			if err != nil {
				fmt.Printf("Error deleting key %s: %v\n", keys[i], err)
			}
		}
	}
}
