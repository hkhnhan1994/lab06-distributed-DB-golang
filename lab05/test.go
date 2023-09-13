package main

import (
	"testing"
)

func TestShardedDB(t *testing.T) {
	maxShard := 4
	db := NewShardedDB(maxShard)

	// Test adding a node
	db.AddNode()
	if len(db.Shards) != maxShard+1 {
		t.Errorf("Expected %d shards after adding a node, got %d", maxShard+1, len(db.Shards))
	}

	// Test removing a node
	nodeToRemove := 1
	db.RemoveNode(nodeToRemove)
	if db.Shards[nodeToRemove].Active {
		t.Errorf("Expected Node %d to be inactive after removal", nodeToRemove)
	}

	// Test getting node ID for a key
	key := "key4"
	expectedNodeID := len(key) % (maxShard + 1) // Node count after adding a node
	nodeID := db.GetNodeID(key)
	if nodeID != expectedNodeID {
		t.Errorf("Expected key '%s' to belong to Node %d, got Node %d", key, expectedNodeID, nodeID)
	}
}

func TestShardedDB_MaxNodes(t *testing.T) {
	maxShard := 2 // Set a small number of max nodes for testing
	db := NewShardedDB(maxShard)

	// Attempt to add more nodes than the maximum allowed
	for i := 0; i < maxShard+1; i++ {
		db.AddNode()
	}

	if len(db.Shards) != maxShard {
		t.Errorf("Expected only %d shards after adding more nodes than allowed, got %d", maxShard, len(db.Shards))
	}
}

func TestShardedDB_InvalidNodeID(t *testing.T) {
	maxShard := 3
	db := NewShardedDB(maxShard)

	// Attempt to remove an invalid node
	invalidNodeID := -1
	db.RemoveNode(invalidNodeID)
	if !db.Shards[0].Active {
		t.Errorf("Expected Node 0 to be active after attempting to remove an invalid node")
	}
}
