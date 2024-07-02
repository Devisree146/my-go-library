package main

import (
	"net/http"
	"time"
	"unified/in_memory"

	"github.com/gin-gonic/gin"
)

type CacheEntry struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	TTL   string `json:"ttl"`
}

func main() {
	cache := in_memory.NewInMemoryCache(3)
	router := gin.Default()

	router.POST("/cache", func(c *gin.Context) {
		var data CacheEntry
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		ttlDuration, err := time.ParseDuration(data.TTL)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TTL format"})
			return
		}

		cache.Set(data.Key, data.Value, ttlDuration)
		c.JSON(http.StatusCreated, gin.H{"message": "Key set successfully"})
	})

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

	router.DELETE("/cache/all", func(c *gin.Context) {
		cache.DeleteAll()
		c.JSON(http.StatusOK, gin.H{"message": "All keys deleted successfully"})
	})

	router.GET("/cache/all", func(c *gin.Context) {
		cachedKeys := cache.GetAllKeys()

		c.JSON(http.StatusOK, gin.H{"keys": cachedKeys})
	})

	router.Run(":8080")
}
