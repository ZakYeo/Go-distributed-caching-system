package main

func main() {
	go LaunchShard()
	go LaunchShard()
	go LaunchShard()
	LaunchServer()
}
