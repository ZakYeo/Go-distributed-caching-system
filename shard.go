package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Cache struct {
	Items map[string]CacheItem `json:"items"`
	mu    sync.Mutex           `json:"-"`
}

type CacheItem struct {
	Value interface{}
}

var currentCache = Cache{
	Items: map[string]CacheItem{},
}

func AddCacheItem(key string, item string) {
	cacheItem := CacheItem{
		Value: item,
	}
	currentCache.mu.Lock()
	defer currentCache.mu.Unlock()
	currentCache.Items[key] = cacheItem
}

func GetCacheEndpointWrapper(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /getCache request\n")

	jsonBytes, err := json.Marshal(GetCache())
	if err != nil {
		http.Error(w, "Failed to encode cache as JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonBytes)
}

func AddCacheItemEndpointWrapper(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /addCacheItem request\n")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if requestBody.Key == "" || requestBody.Value == "" {
		http.Error(w, "Key and value cannot be empty", http.StatusBadRequest)
		return
	}

	AddCacheItem(requestBody.Key, requestBody.Value)

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Cache item added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
	}
}

func GetCache() *Cache {
	currentCache.mu.Lock()
	defer currentCache.mu.Unlock()

	// Deep copy the Items map
	copy := Cache{
		Items: make(map[string]CacheItem),
	}

	for key, value := range currentCache.Items {
		copy.Items[key] = value
	}

	return &copy
}

func GetCacheItem(key string) CacheItem {
	currentCache.mu.Lock()
	defer currentCache.mu.Unlock()
	return currentCache.Items[key]
}

func ClearCache() {
	currentCache.mu.Lock()
	defer currentCache.mu.Unlock()
	currentCache.Items = make(map[string]CacheItem)
}

func RemoveCacheItem(key string) {
	currentCache.mu.Lock()
	defer currentCache.mu.Unlock()
	delete(currentCache.Items, key)
}
