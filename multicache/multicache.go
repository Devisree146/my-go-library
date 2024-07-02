package multicache

import (
	"net/http"
	"time"
	"unified/in_memory"
	"unified/redis_cache"

	"github.com/gin-gonic/gin"
)

var cacheInMemory *in_memory.InMemoryCache
var cacheRedis *redis_cache.Cache

func init() {
	cacheInMemory = in_memory.NewInMemoryCache(3)
	cacheRedis = redis_cache.NewRedisCache("localhost:6379", "", 0, 3)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/cache", func(c *gin.Context) {
		var data struct {
			Key   string `json:"key"`
			Value int    `json:"value"`
			TTL   string `json:"ttl"`
		}
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		ttlDuration, err := time.ParseDuration(data.TTL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TTL format"})
			return
		}

		cacheInMemory.Set(data.Key, data.Value, ttlDuration)
		cacheRedis.Set(data.Key, data.Value, ttlDuration)

		c.JSON(http.StatusCreated, gin.H{"message": "Key set successfully"})
	})

	router.GET("/cache", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Key not provided"})
			return
		}

		value, err := cacheInMemory.Get(key)
		if err != nil {
			value, err = cacheRedis.Get(key)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	router.DELETE("/cache", handleDelete)

	router.DELETE("/cache/all", func(c *gin.Context) {
		cacheInMemory.DeleteAll()
		cacheRedis.DeleteAll()
		c.JSON(http.StatusOK, gin.H{"message": "All keys deleted successfully"})
	})

	router.GET("/cache/all", func(c *gin.Context) {
		cachedKeysInMemory := cacheInMemory.GetAllKeys() // Assuming GetAllKeys() retrieves all keys
		cachedKeysRedis := cacheRedis.GetAllKeys()       // Assuming GetAllKeys() retrieves all keys

		allCachedKeys := append(cachedKeysInMemory, cachedKeysRedis...)

		c.JSON(http.StatusOK, gin.H{"keys": allCachedKeys})
	})

	return router
}

func handleDelete(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key not provided"})
		return
	}

	err := cacheInMemory.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting key from in-memory cache"})
		return
	}

	err = cacheRedis.Delete(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting key from Redis cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}
