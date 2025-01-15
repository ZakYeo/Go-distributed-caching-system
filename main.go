package main

import (
	"time"
)

func main() {
	go LaunchShard()
	go LaunchShard()
	go LaunchShard()
	time.Sleep(1 * time.Second) // Allow shards to initialize before starting the server
	LaunchServer()
}
