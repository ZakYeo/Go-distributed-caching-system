package main

import "testing"

func TestCanSetCacheWithOneItem(t *testing.T) {

  AddCacheItem("new_cache_item")
  currentCache := GetCache()

  if(currentCache.Items[0].Value != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", currentCache.Items[0].Value, "new_cache_item")
  }

  ClearCache()
}


func TestCanSetCacheWithMultipleItems(t *testing.T) {

  AddCacheItem("new_cache_item")
  currentCache := GetCache()

  if(currentCache.Items[0].Value != "new_cache_item"){
    t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", currentCache.Items[0].Value, "new_cache_item")
  }

  AddCacheItem("new_cache_item2")
  currentCache = GetCache()
  
  if(currentCache.Items[1].Value != "new_cache_item2"){
    t.Errorf("Unable to successfully set new_cache_item2 into cache. Got: %s, Want: %s", currentCache.Items[1].Value, "new_cache_item2")
  }
  ClearCache()
}
