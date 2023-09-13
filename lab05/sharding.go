package main

import (
	"fmt"
	"sync"
)

// Shard represents a shard/node in the distributed key-value store.
type Shard struct {
	ID     int
	Data   map[string]string
	Lock   sync.Mutex
	Active bool
}

// ShardedDB represents a sharded key-value database.
type ShardedDB struct {
	Shards   []*Shard
	MaxShard int
}

// NewShardedDB creates a new sharded key-value database.
func NewShardedDB(maxShard int) *ShardedDB {
	db := &ShardedDB{
		MaxShard: maxShard,
	}

	// Initialize the shards
	db.Shards = make([]*Shard, maxShard)
	for i := 0; i < maxShard; i++ {
		db.Shards[i] = &Shard{
			ID:     i,
			Data:   make(map[string]string),
			Active: true,
		}
	}

	return db
}

// AddNode adds a new shard/node to the database.
func (db *ShardedDB) AddNode() {
	if len(db.Shards) < db.MaxShard {
		newShard := &Shard{
			ID:     len(db.Shards),
			Data:   make(map[string]string),
			Active: true,
		}
		db.Shards = append(db.Shards, newShard)
		fmt.Printf("Added Node %d\n", newShard.ID)
	} else {
		fmt.Println("Cannot add more nodes. Max nodes reached.")
	}
}

// RemoveNode removes a shard/node from the database by ID.
func (db *ShardedDB) RemoveNode(nodeID int) {
	if nodeID >= 0 && nodeID < len(db.Shards) {
		db.Shards[nodeID].Active = false
		fmt.Printf("Removed Node %d\n", nodeID)
	} else {
		fmt.Println("Invalid node ID.")
	}
}

// GetNodeID returns the node ID (shard ID) responsible for a given key.
func (db *ShardedDB) GetNodeID(key string) int {
	// Simple hashing to determine the shard/node ID for a key
	hash := len(key) % len(db.Shards)
	return hash
}

func main() {
	maxShard := 4 // Maximum number of nodes (shards)
	db := NewShardedDB(maxShard)

	// Initial data insertion
	db.Shards[0].Data["key1"] = "value1"
	db.Shards[1].Data["key2"] = "value2"
	db.Shards[2].Data["key3"] = "value3"

	// Adding a new node
	db.AddNode()

	// Removing a node
	db.RemoveNode(1)

	// Getting node ID for a key
	key := "key4"
	nodeID := db.GetNodeID(key)
	fmt.Printf("Key '%s' belongs to Node %d\n", key, nodeID)
}
