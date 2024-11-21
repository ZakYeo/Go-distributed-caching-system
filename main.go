package main

import (
  "fmt"
  "sync"
)

type Cache struct {
  Items map[string]CacheItem
  mu sync.Mutex
}

type CacheItem struct {
  Value interface{}
}

var currentCache = Cache{
  Items: map[string]CacheItem{},
}

func main() {
    fmt.Println("Hello, World!") 
}

func AddCacheItem(key string, item string){
  cacheItem:= CacheItem{
    Value: item,
  }
  currentCache.mu.Lock()
  defer currentCache.mu.Unlock()
  currentCache.Items[key] = cacheItem
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
