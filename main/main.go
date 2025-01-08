package main

import (
  "fmt"
  "sync"
	"net/http"
	"os"
	"errors"
  "encoding/json"
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

func main() {
	http.HandleFunc("/get/cache", getCacheEndpointWrapper)

	err := http.ListenAndServe(":8080", nil)
  if errors.Is(err, http.ErrServerClosed) {
      fmt.Printf("server closed\n")
    } else if err != nil {
      fmt.Printf("error starting server: %s\n", err)
      os.Exit(1)
    }
}

func AddCacheItem(key string, item string){
  cacheItem:= CacheItem{
    Value: item,
  }
  currentCache.mu.Lock()
  defer currentCache.mu.Unlock()
  currentCache.Items[key] = cacheItem
}

func getCacheEndpointWrapper(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /getCache request\n")

    jsonBytes, err := json.Marshal(GetCache())
    if err != nil {
        http.Error(w, "Failed to encode cache as JSON", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    w.Write(jsonBytes)
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

func GetCacheItem(key string)(CacheItem){
  currentCache.mu.Lock()
  defer currentCache.mu.Unlock()
  return currentCache.Items[key]
}

func ClearCache(){
  currentCache.mu.Lock()
  defer currentCache.mu.Unlock()
  currentCache.Items = make(map[string]CacheItem)
}

func RemoveCacheItem(key string){
  currentCache.mu.Lock()
  defer currentCache.mu.Unlock()
  delete(currentCache.Items, key)
}
