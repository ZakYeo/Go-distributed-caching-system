package main

import "fmt"

type Cache struct {
  Items string
}

func main() {
    fmt.Println("Hello, World!") 
}

func SetCache(item string, cache Cache)(Cache){
  cache.Items = item
  return cache
}
