package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

/**
This file contains the endpoints and server functions for an individual shard
A single shard is a server that contains its own cache
Managed by the central server
*/


func AddShardCacheItemEndpointWrapper(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /addCacheItem request on shard\n")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	if requestBody.Key == "" || requestBody.Value == "" {
		http.Error(w, "Key and value cannot be empty", http.StatusBadRequest)
		return
	}

	AddCacheItem(requestBody.Key, requestBody.Value)

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: "Cache item added successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
	}
}

func GetShardCacheEndpointWrapper(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /cache/get/{key} request on shard\n")

    vars := mux.Vars(r)
    key := vars["key"]

    if key == "" {
        http.Error(w, "Key cannot be empty", http.StatusBadRequest)
        return
    }

    if key == "all"{
        GetAllShardCache(w, r)
        return
    }

	cacheItem := GetCacheItem(key)
    if cacheItem.Value == nil {
        http.Error(w, "Cache item not found", http.StatusNotFound)
        return
    }

	// TODO:
	// For now, assume the values for cache items are only strings
	value, ok := cacheItem.Value.(string)
    if !ok {
        http.Error(w, "Cache item value is not a valid string", http.StatusInternalServerError)
        return
    }
	fmt.Printf("Got cache item: %s\n", value)

    response := struct {
        Status  string `json:"status"`
        Key     string `json:"key"`
        Value   string `json:"value"`
    }{
        Status: "success",
        Key:    key,
        Value:  value,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
    }
}


func GetAllShardCache(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("got /cache/get/all request on shard\n")

	cacheItem := GetCache()
    if cacheItem.Items == nil {
        http.Error(w, "Cache item not found", http.StatusNotFound)
        return
    }

   items := make(map[string]string)
    for key, cacheItem := range cacheItem.Items {
		// TODO: For now, assume all cache item values are strings
        if strValue, ok := cacheItem.Value.(string); ok {
            items[key] = strValue
        } else {
            items[key] = fmt.Sprintf("%v", cacheItem.Value)
        }
    }

    response := struct {
        Status string            `json:"status"`
        Items  map[string]string `json:"value"`
    }{
        Status: "success",
        Items:  items,
    } 

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
    }
}