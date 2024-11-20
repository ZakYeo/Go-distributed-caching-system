package main

import "testing"

func TestCanSetCacheWithOneItem(t *testing.T) {
  testCacheItems:= Cache {
    Items: []CacheItem{},
  }

  newTestCache := SetCache("new_cache_item", testCacheItems)


  if(newTestCache.Items[0].Value != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", newTestCache.Items[0].Value, "new_cache_item")
  }
}
