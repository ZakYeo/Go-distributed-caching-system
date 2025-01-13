package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
	"os"
	"github.com/gorilla/mux"
)

/**
This file contains the endpoints and server functions for an individual shard
A single shard is a server that contains its own cache
Managed by the central server
*/

func LaunchShard(shard_number string) {
	/**
	Launch a shard and its endpoints
	*/
	r := mux.NewRouter()
	r.HandleFunc("/cache/get/{key}", GetShardCacheEndpointWrapper)
	r.HandleFunc("/cache/add", AddShardCacheItemEndpointWrapper)

	port := fmt.Sprintf(":808%s", shard_number)

	fmt.Printf("Shard %s launched\n", shard_number)
	err := http.ListenAndServe(port, r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server on port %s: %s\n", port, err)
		os.Exit(1)
	}
}

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