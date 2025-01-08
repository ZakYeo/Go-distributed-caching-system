package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var numberOfShards = 0

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func LaunchServer() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	http.HandleFunc("/addCacheItem", addCacheItemEndpoint)

	fmt.Printf("Caching Server Launched\n")
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
func addCacheItemEndpoint(w http.ResponseWriter, r *http.Request) {
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

	// Call AddCacheItemServer to process the request
	AddCacheItemServer(requestBody.Key, requestBody.Value)

	// Respond with success
	response := map[string]string{
		"status":  "success",
		"message": "Cache item processed and sent to the appropriate shard",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func AddCacheItemServer(key string, item string) {
	// Calculate where this cache item should go
	// Then send it to the right shard
	shardNumberToSendTo := hashAndModulo(key, numberOfShards)
	shardAddress := "https://localhost:8081"

	// Prepare the payload
	payload := map[string]string{
		"key":   key,
		"value": item,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Failed to marshal payload: %v\n", err)
		return
	}

	// Send the POST request to the appropriate shard
	resp, err := http.Post(shardAddress+"/addCacheItem", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Printf("Failed to send request to shard %d: %v\n", shardNumberToSendTo, err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Successfully added cache item to shard %d\n", shardNumberToSendTo)
	} else {
		fmt.Printf("Failed to add cache item to shard %d: %s\n", shardNumberToSendTo, resp.Status)
	}

}

func hashAndModulo(key string, numberOfShards int) int {
	// Compute the SHA-256 hash of the key
	hash := sha256.Sum256([]byte(key))

	// Use the first 8 bytes of the hash to convert it into an integer
	hashValue := binary.BigEndian.Uint64(hash[:8])

	// Apply modulo to get the bucket index
	bucket := int(hashValue % uint64(numberOfShards))

	return bucket
}

func LaunchShard(shard_number string) {
	http.HandleFunc("/cache/get", GetCacheEndpointWrapper)
	http.HandleFunc("/cache/add", AddCacheItemEndpointWrapper)

	port := fmt.Sprintf(":808%s", shard_number)

	fmt.Printf("Shard %s launched\n", shard_number)
	numberOfShards += 1
	err := http.ListenAndServe(port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server on port %s: %s\n", port, err)
		os.Exit(1)
	}
}
