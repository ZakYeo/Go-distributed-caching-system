package main

import "fmt"

type Cache struct {
  Items map[string]CacheItem
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
  currentCache.Items[key] = cacheItem
}

func GetCache()(*Cache){
    copy := currentCache // Make a shallow copy of the cache
    return &copy         // Return a pointer to the copy
}

func ClearCache(){
  currentCache = Cache{
    Items: map[string]CacheItem{},
  }
}
