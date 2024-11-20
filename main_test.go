package main

import "testing"

func TestSetCache(t *testing.T) {
  testCacheItems:= Cache {
    Items: "",
  }

  newTestCache := SetCache("new_cache_item", testCacheItems)

  if(newTestCache.Items != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", newTestCache.Items, "new_cache_item")
  }
}
