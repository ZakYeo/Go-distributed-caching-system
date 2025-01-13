package main

import (
	"testing"
)

func TestGetShardNumberToSendTo(t *testing.T) {
	testCases := []struct {
		key            string
		numberOfShards int
		expectedShard  int
	}{
		{"key1", 10, 1},    // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 10 -> shard 2
		{"key2", 5, 3},     // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 5 -> shard 3
		{"", 200, 52},      // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 200 -> shard 52
		{"long_key", 3, 1}, // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 3 -> shard 1
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			result := getShardNumberToSendTo(tc.key, tc.numberOfShards)
			if result != tc.expectedShard {
				t.Errorf("For key %q and %d shards: expected %d, got %d", tc.key, tc.numberOfShards, tc.expectedShard, result)
			}
		})
	}
}
