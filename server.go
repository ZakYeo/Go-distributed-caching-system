package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

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

	fmt.Printf("Caching Server Launched\n")
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func LaunchShard(shard_number string) {
	http.HandleFunc("/get/cache", getCacheEndpointWrapper)

	port := fmt.Sprintf(":808%s", shard_number)

	fmt.Printf("Shard %s launched\n", shard_number)
	err := http.ListenAndServe(port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server on port %s: %s\n", port, err)
		os.Exit(1)
	}
}
