package main

import "testing"

func TestCanSetCacheWithOneItem(t *testing.T) {
  testCacheItems:= Cache {
    Items: []CacheItem{},
  }

  SetCacheItem("new_cache_item", testCacheItems)
  currentCache = GetCache()


  if(currentCache.Items[0].Value != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", currentCache.Items[0].Value, "new_cache_item")
  }
}


/*func TestCanSetCacheWithMultipleItems(t *testing.T) {
  testCacheItems:= Cache {
    Items: []CacheItem{},
  }

  newTestCache := SetCache("new_cache_item", testCacheItems)
  newTestCache2 := SetCache("new_cache_item2", testCacheItems)

  if(newTestCache.Items[0].Value != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", newTestCache.Items[0].Value, "new_cache_item")
  }

  
  if(newTestCache2.Items[1].Value != "new_cache_item2"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", newTestCache.Items[1].Value, "new_cache_item2")
  }
}*/
