package main

import (
	"strconv"
	"sync"
	"testing"
)

func TestCanSetCacheWithOneItem(t *testing.T) {

	AddCacheItem("key1", "new_cache_item")
	cachedItem := GetCacheItem("key1")

	if cachedItem.Value != "new_cache_item" {
		t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", cachedItem.Value, "new_cache_item")
	}

	ClearCache()
}

func TestCanSetCacheWithMultipleItems(t *testing.T) {

	AddCacheItem("key1", "new_cache_item")
	AddCacheItem("key2", "new_cache_item2")
	cachedItem1 := GetCacheItem("key1")
	cachedItem2 := GetCacheItem("key2")

	if cachedItem1.Value != "new_cache_item" {
		t.Errorf("Unable to successfully set new_cache_item into cache. Got: %s, Want: %s", cachedItem1.Value, "new_cache_item")
	}

	if cachedItem2.Value != "new_cache_item2" {
		t.Errorf("Unable to successfully set new_cache_item2 into cache. Got: %s, Want: %s", cachedItem2.Value, "new_cache_item2")
	}
	ClearCache()
}

func TestCanRemoveCacheItem(t *testing.T) {

	AddCacheItem("key1", "cache_item")
	RemoveCacheItem("key1")
	cache := GetCache()
	if len(cache.Items) > 0 {
		t.Errorf("Cache item %s not removed from cache.", "cache_item")
	}
	ClearCache()

}

func TestWriteFromMultipleThreadsAtOnceToCache(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		itemToAdd := strconv.Itoa(i)
		wg.Add(1)
		go func(itemToAdd string) {
			defer wg.Done()
			AddCacheItem(itemToAdd, itemToAdd)
		}(itemToAdd)
	}

	wg.Wait()

	cache := GetCache()

	if len(cache.Items) != 1000 {
		t.Errorf("Cache items do not equal %s", "1000")
	}
}
