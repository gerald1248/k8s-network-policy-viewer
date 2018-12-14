package main

import (
	"testing"
)

func TestInitializeEdgeMap(t *testing.T) {
	edgeMap := make(map[string][]string)
	namespacePodMap := make(map[string][]string)

	namespacePodMap["default"] = []string{"alice", "eve", "bob"}
	initializeEdgeMap(&edgeMap, &namespacePodMap)

	expected := 2
	if len(edgeMap["alice"]) != expected || len(edgeMap["eve"]) != expected || len(edgeMap["bob"]) != expected {
		t.Errorf("Each pod must have two egress connections: got %v", edgeMap)
	}
}

func TestDeduplicateEdgeMap(t *testing.T) {
	edgeMap := make(map[string][]string)
	edgeMap["default"] = []string{"alice", "bob", "eve", "eve"}

	deduplicateEdgeMap(&edgeMap)

	expected := 3
	count := len(edgeMap["default"])
	if count != expected {
		t.Errorf("Deduplicated slice must have %d elements: got %d", expected, count)
	}
}

func TestUnique(t *testing.T) {
	slice := []string{"alice", "bob", "eve", "eve"}
	deduplicated := unique(slice)

	expected := 3
	count := len(deduplicated)

	if count != expected {
		t.Errorf("Deduplicated slice must have %d elements: got %d", expected, count)
	}
}
