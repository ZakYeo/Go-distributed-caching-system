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
)

/**
This file contains the server that manages the shards
*/

func LaunchServer() {
	/**
	Launch the server that manages the shards
	This server is designed to handle creating / removing shards,
	adding / removing cache from these shards, etc.
	*/
	http.HandleFunc("/addCacheItem", addCacheItemEndpointWrapper)

	fmt.Printf("Caching Server Launched\n")
	err := http.ListenAndServe(":3333", nil)
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
	callAddCacheItemEndpointOfShard(requestBody.Key, requestBody.Value, hashAndModulo(requestBody.Key, 1), shardAddress)

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

func hashAndModulo(key string, numberOfShards int) int {
	/**
	Hashes the key and applies modulo to get the shard number
	Used to determine which shard to send the cache item to
	*/
	// Compute the SHA-256 hash of the key
	hash := sha256.Sum256([]byte(key))
	// Use the first 8 bytes of the hash to convert it into an integer
	hashValue := binary.BigEndian.Uint64(hash[:8])
	// Apply modulo to get the bucket index
	bucket := int(hashValue % uint64(numberOfShards))

	return bucket
}

