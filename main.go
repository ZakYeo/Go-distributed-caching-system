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

func AddCacheItem(item string){
  cacheItem:= CacheItem{
    Value: item,
  }
  currentCache.Items = append(currentCache.Items, cacheItem)
}

func GetCache()(*Cache){
    copy := currentCache // Make a shallow copy of the cache
    return &copy         // Return a pointer to the copy
}

func ClearCache(){
  currentCache = Cache{
    Items: []CacheItem{},
  }
}
