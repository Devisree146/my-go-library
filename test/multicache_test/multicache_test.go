/*package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"unified/multicache"

	"github.com/stretchr/testify/assert"
)

func TestCacheAPI(t *testing.T) {
	router := multicache.SetupRouter()

	// Test POST /cache endpoint
	t.Run("POST /cache", func(t *testing.T) {
		payload := `{"key": "test_key", "value": 42, "ttl": "10s"}`
		req, _ := http.NewRequest("POST", "/cache", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Key set successfully", response["message"])
	})

	// Test GET /cache endpoint
	t.Run("GET /cache", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache?key=test_key", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test_key", response["key"])
		assert.Equal(t, float64(42), response["value"].(float64))
	})

	// Test DELETE /cache/all endpoint
	t.Run("DELETE /cache/all", func(t *testing.T) {
		// Verify keys before deletion
		reqBeforeDelete, _ := http.NewRequest("GET", "/cache/all", nil)
		wBeforeDelete := httptest.NewRecorder()
		router.ServeHTTP(wBeforeDelete, reqBeforeDelete)
		assert.Equal(t, http.StatusOK, wBeforeDelete.Code)

		var responseBeforeDelete map[string][]string
		err := json.Unmarshal(wBeforeDelete.Body.Bytes(), &responseBeforeDelete)
		assert.NoError(t, err)
		assert.Contains(t, responseBeforeDelete["keys"], "test_key")

		// Perform DELETE /cache/all
		reqDeleteAll, _ := http.NewRequest("DELETE", "/cache/all", nil)
		wDeleteAll := httptest.NewRecorder()
		router.ServeHTTP(wDeleteAll, reqDeleteAll)
		assert.Equal(t, http.StatusOK, wDeleteAll.Code)

		var responseDeleteAll map[string]string
		err = json.Unmarshal(wDeleteAll.Body.Bytes(), &responseDeleteAll)
		assert.NoError(t, err)
		assert.Equal(t, "All keys deleted successfully", responseDeleteAll["message"])

		// Wait for deletion to complete (adjust as needed)
		time.Sleep(100 * time.Millisecond)

		// Verify keys after deletion
		reqAfterDelete, _ := http.NewRequest("GET", "/cache/all", nil)
		wAfterDelete := httptest.NewRecorder()
		router.ServeHTTP(wAfterDelete, reqAfterDelete)
		assert.Equal(t, http.StatusOK, wAfterDelete.Code)

		var responseAfterDelete map[string][]string
		err = json.Unmarshal(wAfterDelete.Body.Bytes(), &responseAfterDelete)
		assert.NoError(t, err)
		assert.NotContains(t, responseAfterDelete["keys"], "test_key", "Expected key 'test_key' to be deleted")
	})

	// Test GET /cache/all endpoint after deletion
	t.Run("GET /cache/all after deletion", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/cache/all", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string][]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotContains(t, response["keys"], "test_key", "Expected key 'test_key' to be deleted")
	})
}
*/