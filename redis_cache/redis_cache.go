package redis_cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache struct {
	client  *redis.Client
	maxSize int
}

func NewRedisCache(address, password string, db, maxSize int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &Cache{
		client:  client,
		maxSize: maxSize,
	}
}

func (c *Cache) Set(key string, value int, ttl time.Duration) error {
	ctx := context.Background()
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}

	// Perform LRU eviction if cache exceeds maxSize
	c.performLRUEviction()

	return nil
}

func (c *Cache) Get(key string) (int, error) {
	ctx := context.Background()
	val, err := c.client.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}

	// Perform LRU eviction if cache exceeds maxSize
	c.performLRUEviction()

	return val, nil
}

func (c *Cache) Delete(key string) error {
	ctx := context.Background()
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	// No need to perform LRU eviction on delete operation

	return nil
}

func (c *Cache) DeleteAll() {
	ctx := context.Background()
	c.client.FlushDB(ctx).Err()

	// No need to perform LRU eviction on delete all operation
}

func (c *Cache) GetAllKeys() []string {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return []string{}
	}

	// Perform LRU eviction if cache exceeds maxSize
	c.performLRUEviction()

	return keys
}

func (c *Cache) performLRUEviction() {
	ctx := context.Background()
	keys, err := c.client.Keys(ctx, "*").Result()
	if err != nil {
		return
	}

	// If cache size exceeds maxSize, perform LRU eviction
	if len(keys) > c.maxSize {
		var oldestKey string
		var oldestTime time.Time

		// Find the oldest accessed key
		for _, key := range keys {
			ts, err := c.client.ObjectIdleTime(ctx, key).Result()
			if err != nil {
				continue
			}
			lastAccessed := time.Now().Add(-time.Second * time.Duration(ts))
			if oldestTime.IsZero() || lastAccessed.Before(oldestTime) {
				oldestTime = lastAccessed
				oldestKey = key
			}
		}

		// Delete the oldest key
		if oldestKey != "" {
			c.client.Del(ctx, oldestKey)
		}
	}
}
