package main

import (
	"reflect"
	"testing"
)

func TestConstructSuffixTree(t *testing.T) {
	// Test case 1: Basic overlap
	reads := []string{"ACGCG", "CGCGT"}
	minOverlap := 2
	expectedOverlapGraph := map[string]map[string]int{
		"ACGCG": {"CGCGT": 4},
		"CGCGT": {},
	}
	overlapGraph := overlap(reads, minOverlap)

	if !reflect.DeepEqual(overlapGraph, expectedOverlapGraph) {
		t.Errorf("Expected graph %v, but got %v", expectedOverlapGraph, overlapGraph)
	}

	// Test case 2: No overlap
	reads = []string{"ACGCG", "CGTGA"}
	minOverlap = 3
	expectedOverlapGraph = map[string]map[string]int{
		"ACGCG": {},
		"CGTGA": {},
	}
	overlapGraph = overlap(reads, minOverlap)

	if !reflect.DeepEqual(overlapGraph, expectedOverlapGraph) {
		t.Errorf("Expected graph %v, but got %v", expectedOverlapGraph, overlapGraph)
	}
}
