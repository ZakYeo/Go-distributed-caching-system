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
		{"key1", 10, 2},    // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 10 -> plus 1 -> shard 2
		{"key2", 5, 4},     // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 5 -> plus 1 -> shard 4
		{"", 200, 53},      // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 200 -> plus 1 -> shard 53
		{"long_key", 3, 2}, // Precomputed: SHA-256 hash -> first 8 bytes -> modulo 3 -> plus 1 -> shard 2
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
