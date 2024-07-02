package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"unified/multicache"

	"github.com/gin-gonic/gin"
)

// setupRouter initializes the Gin router with multicache routes.
func setupRouter() *gin.Engine {
	return multicache.SetupRouter()
}

// BenchmarkSetCache benchmarks the performance of setting a cache entry.
func BenchmarkSetCache(b *testing.B) {
	router := setupRouter()

	payload := map[string]interface{}{
		"key":   "benchmark_key_set",
		"value": 42,
		"ttl":   "10s",
	}

	payloadBytes, _ := json.Marshal(payload)

	b.ResetTimer() // Reset the timer before starting benchmark

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", "/cache", bytes.NewReader(payloadBytes))
		if err != nil {
			b.Fatalf("could not create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}

// BenchmarkGetCache benchmarks the performance of retrieving a cache entry.
func BenchmarkGetCache(b *testing.B) {
	router := setupRouter()

	payload := map[string]interface{}{
		"key":   "benchmark_key_get",
		"value": 42,
		"ttl":   "10s",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Set the cache entry before benchmarking
	req, err := http.NewRequest("POST", "/cache", bytes.NewReader(payloadBytes))
	if err != nil {
		b.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	b.ResetTimer() // Reset the timer before starting benchmark

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", "/cache?key=benchmark_key_get", nil)
		if err != nil {
			b.Fatalf("could not create request: %v", err)
		}
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}

// BenchmarkDeleteCache benchmarks the performance of deleting a cache entry.
func BenchmarkDeleteCache(b *testing.B) {
	router := setupRouter()

	payload := map[string]interface{}{
		"key":   "benchmark_key_delete",
		"value": 42,
		"ttl":   "10s",
	}

	payloadBytes, _ := json.Marshal(payload)

	// Set the cache entry before benchmarking
	req, err := http.NewRequest("POST", "/cache", bytes.NewReader(payloadBytes))
	if err != nil {
		b.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	b.ResetTimer() // Reset the timer before starting benchmark

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("DELETE", "/cache?key=benchmark_key_delete", nil)
		if err != nil {
			b.Fatalf("could not create request: %v", err)
		}
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
	}
}
