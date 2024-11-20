package main

import "fmt"

type Cache struct {
  Items []CacheItem
}

type CacheItem struct {
  Value interface{}
}

func main() {
    fmt.Println("Hello, World!") 
}

func SetCache(item string, cache Cache)(Cache){
  cacheItem:= CacheItem{
    Value: item,
  }
  cache.Items = append(cache.Items, cacheItem)
  return cache
}
