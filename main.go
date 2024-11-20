package main

import "fmt"

type Cache struct {
  Items []CacheItem
}

type CacheItem struct {
  Value interface{}
}

var currentCache = Cache{
  Items: []CacheItem{},
}

func main() {
    fmt.Println("Hello, World!") 
}

func SetCacheItem(item string, cache Cache){
  cacheItem:= CacheItem{
    Value: item,
  }
  currentCache.Items = append(cache.Items, cacheItem)
}

func GetCache()(*Cache){
    copy := currentCache // Make a shallow copy of the cache
    return &copy         // Return a pointer to the copy
}
