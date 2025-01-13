# Go-distributed-caching-system

My attempt at making a distributed In-Memory caching system built in Go.

The following is a list of each file and their use:

- cache.go:
  - Contains the raw implementation of the caching system
  - e.g Add cache item, Delete cache item
- shard_server.go:
  - Contains the server & endpoint wrappers to interact with the caching system
  - Each "shard" is an individual cache, with endpoints used to interact with the cache using cache.go functions
- central_shard_server_manager.go:
  - The central server used to interact with each individual shard through the shard_server.go endpoints
  - Is intended to use to interact with each shard, remove shards, add shards, etc.
- main.go:
  - Main file to run to launch a server and shard(s)

# Try it out

- Run in the root: `go run .`
- Now in a Terminal, add a cache item via

````
curl -v -X POST http://localhost:3333/addCacheItem -H "Content-Type: application/json" -d "{\"key\":\"sampleKey\",\"value\":\"sampleValue\"}"```
- Now in a Terminal, get the cache item using its key
````

```
curl -X GET http://localhost:3333/getCacheItem/sampleKey

```
