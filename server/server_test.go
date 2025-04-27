package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZivojinMilutinovic/golang-coding-test/store_api"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	api := &API{
		Store: store_api.NewStore(),
	}
	r := gin.Default()
	api.registerRoutes(r)
	return r
}

func TestSet(t *testing.T) {
	r := setupRouter()

	body := map[string]interface{}{
		"value": "testValue",
		"ttl":   60,
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/set/testKey", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate(t *testing.T) {
	r := setupRouter()

	// First, set a value
	body := map[string]interface{}{
		"value": "initialValue",
		"ttl":   60,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/set/testKey", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Now, update the value
	updateBody := map[string]interface{}{
		"value": "updatedValue",
	}
	updateBodyBytes, _ := json.Marshal(updateBody)

	req, _ = http.NewRequest("POST", "/update/testKey", bytes.NewBuffer(updateBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify that the value was updated
	req, _ = http.NewRequest("GET", "/get/testKey", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updatedValue")
}

func TestRemove(t *testing.T) {
	r := setupRouter()

	// First, set a value
	body := map[string]interface{}{
		"value": "testValue",
		"ttl":   60,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/set/testKey", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Now, remove the value
	req, _ = http.NewRequest("DELETE", "/remove/testKey", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify that the value is removed
	req, _ = http.NewRequest("GET", "/get/testKey", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPush(t *testing.T) {
	r := setupRouter()

	// First, set a list
	body := map[string]interface{}{
		"value": []string{},
		"ttl":   60,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/set/myList", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Now, push a value to the list
	pushBody := map[string]interface{}{
		"value": "newItem",
	}
	pushBodyBytes, _ := json.Marshal(pushBody)

	req, _ = http.NewRequest("POST", "/push/myList", bytes.NewBuffer(pushBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Verify that the value was pushed
	req, _ = http.NewRequest("GET", "/get/myList", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "newItem")
}

func TestPopEmptyList(t *testing.T) {
	r := setupRouter()

	// First, set an empty list
	body := map[string]interface{}{
		"value": []string{},
		"ttl":   60,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/set/emptyList", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Now, pop a value from the empty list
	req, _ = http.NewRequest("POST", "/pop/emptyList", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
