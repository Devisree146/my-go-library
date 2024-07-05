package main

import (
	"fmt"
	"net/http"
	"time"
	"unified/redis_cache" // Import your redis_cache package

	"github.com/gin-gonic/gin"
)

type CacheEntry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

const (
	standardTTL = 5 * time.Minute // Standard TTL of 5 minutes
)

func main() {
	// Initialize your Redis cache instance with maxSize of 3
	cache := redis_cache.NewRedisCache("localhost:6379", "", 0, 3)

	// Setup Gin router
	router := gin.Default()

	// POST endpoint to set cache with standard TTL
	router.POST("/cache", func(c *gin.Context) {
		var data CacheEntry
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		err := cache.Set(data.Key, data.Value, redis_cache.StandardTTL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set key"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Key set successfully"})
	})

	// GET endpoint to retrieve cache
	router.GET("/cache", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Key not provided"})
			return
		}

		value, err := cache.Get(key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	})

	// DELETE endpoint to delete cache by key
	router.DELETE("/cache", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Key not provided"})
			return
		}

		err := cache.Delete(key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
	})

	// DELETE endpoint to delete all cache keys
	router.DELETE("/cache/all", func(c *gin.Context) {
		cache.DeleteAll() // Correct usage, no need to assign a value

		c.JSON(http.StatusOK, gin.H{"message": "All keys deleted successfully"})
	})

	// GET endpoint to retrieve all cache keys
	router.GET("/cache/all", func(c *gin.Context) {
		cachedKeys, err := cache.GetAllKeys() // Correct usage, retrieves keys
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get keys"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"keys": cachedKeys})
	})

	// Run the Gin server on port 8080
	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
