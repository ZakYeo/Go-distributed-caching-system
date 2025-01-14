package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	"io"
)

var numberOfShards = 0

/**
This file contains the server that manages the shards
*/

func LaunchServer() {
	/**
	Launch the server that manages the shards
	This server is designed to handle creating / removing shards,
	adding / removing cache from these shards, etc.
	*/
	r := mux.NewRouter()
	r.HandleFunc("/addCacheItem", addCacheItemEndpointWrapper)
	r.HandleFunc("/getCacheItem/{key}", getCacheItemEndpointWrapper)

	fmt.Printf("Caching Server Launched\n")
	err := http.ListenAndServe(":3333", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func addCacheItemEndpointWrapper(w http.ResponseWriter, r *http.Request) {
	/**
	Called when the /addCacheItem of the server is called
	We parse the request body and send the cache item to the appropriate shard's endpoint
	*/
	fmt.Printf("Central server received request: %s\n", time.Now())

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


	shardAddress := "http://localhost:8081"
	callAddCacheItemEndpointOfShard(requestBody.Key, requestBody.Value, getShardNumberToSendTo(requestBody.Key, 1), shardAddress)

	response := map[string]string{
		"status":  "success",
		"message": "Cache item processed and sent to the appropriate shard",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func callAddCacheItemEndpointOfShard(key string, item string, shardNumberToSendTo int, shardAddress string) {
	/**
	Takes the key value pair of the cache item to add,
	Sends it to the specified shard number & that shard's endpoint
	*/

	payload := map[string]string{
		"key":   key,
		"value": item,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	resp, err := http.Post(shardAddress+"/cache/add", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Printf("Failed to send request to shard %d: %v\n", shardNumberToSendTo, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successfully added cache item to shard %d\n", shardNumberToSendTo)
	} else {
		fmt.Printf("Failed to add cache item to shard %d: %s\n", shardNumberToSendTo, resp.Status)
	}
}


func getCacheItemEndpointWrapper(w http.ResponseWriter, r *http.Request) {
	/**
	Called when the /getCacheItem of the server is called
	*/
	fmt.Printf("Central server received request: %s\n", time.Now())

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

    vars := mux.Vars(r)
    key := vars["key"]

    if key == "" {
        http.Error(w, "Key cannot be empty", http.StatusBadRequest)
        return
    }

	callGetCacheItemEndpointOfShard(key, getShardNumberToSendTo(key, numberOfShards))

	response := map[string]string{
		"status":  "success",
		"message": "Gotten cache item",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func callGetCacheItemEndpointOfShard(key string, shardNumberToGetFrom int) {
	/**
	Given a key and the shard number to get the cache item from,
	Sends a request to the specified shard's endpoint to get the cache item
	*/

	shardAddress := fmt.Sprintf("http://localhost:808%d", shardNumberToGetFrom)
	requestURL := fmt.Sprintf("%s/cache/get/%s", shardAddress, key)
	fmt.Printf("Getting cache item from endpoint: %s\n", requestURL)
	resp, err := http.Get(requestURL)
	if err != nil {
		fmt.Printf("Failed to send request to shard %d: %v\n", shardNumberToGetFrom, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successfully got cache item from shard %d\n", shardNumberToGetFrom)
		fmt.Printf("client: status code: %d\n", resp.StatusCode)
		resBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("client: could not read response body: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("client: response body: %s\n", resBody)
	} else {
		fmt.Printf("Failed to get cache item from shard %d: %s\n", shardNumberToGetFrom, resp.Status)
	}
}

func getShardNumberToSendTo(key string, numberOfShards int) int {
	/**
	Hashes the key and applies modulo to get the shard number
	Used to determine which shard to send the cache item to
	*/
	// Compute the SHA-256 hash of the key
	hash := sha256.Sum256([]byte(key))
	// Use the first 8 bytes of the hash to convert it into an integer
	hashValue := binary.BigEndian.Uint64(hash[:8])
	// Apply modulo to get the shard number
	// Plus one to make the shard number 1-indexed
	shardNumberToSendTo := int(hashValue % uint64(numberOfShards)) + 1

	return shardNumberToSendTo
}


func LaunchShard() {
	/**
	Launch a shard and its endpoints
	*/
	r := mux.NewRouter()
	r.HandleFunc("/cache/get/{key}", GetShardCacheEndpointWrapper)
	r.HandleFunc("/cache/add", AddShardCacheItemEndpointWrapper)

	port := fmt.Sprintf(":808%d", numberOfShards+1)

	fmt.Printf("Shard %d launched\n", numberOfShards+1)
	numberOfShards++
	err := http.ListenAndServe(port, r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server on port %s: %s\n", port, err)
		os.Exit(1)
	}
}