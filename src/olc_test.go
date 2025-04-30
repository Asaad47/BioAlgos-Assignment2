package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestConstructSuffixTree(t *testing.T) {
	// Test case 1: Basic overlap
	reads := []string{"ACGCG", "CGCGT"}
	minOverlap := 2
	expectedOverlapGraph := map[string]ContigNode{
		"ACGCG": {
			read:     "ACGCG",
			outEdges: map[string]int{"CGCGT": 4},
			inEdges:  map[string]int{},
		},
		"CGCGT": {
			read:     "CGCGT",
			outEdges: map[string]int{},
			inEdges:  map[string]int{"ACGCG": 4},
		},
	}
	overlapGraph := overlap(reads, minOverlap)

	if !reflect.DeepEqual(overlapGraph, expectedOverlapGraph) {
		t.Errorf("Expected graph %v, but got %v", expectedOverlapGraph, overlapGraph)
	}

	// Test case 2: No overlap
	reads = []string{"ACGCG", "CGTGA"}
	minOverlap = 3
	expectedOverlapGraph = map[string]ContigNode{
		"ACGCG": {
			read:     "ACGCG",
			outEdges: map[string]int{},
			inEdges:  map[string]int{},
		},
		"CGTGA": {
			read:     "CGTGA",
			outEdges: map[string]int{},
			inEdges:  map[string]int{},
		},
	}
	overlapGraph = overlap(reads, minOverlap)

	if !reflect.DeepEqual(overlapGraph, expectedOverlapGraph) {
		t.Errorf("Expected graph %v, but got %v", expectedOverlapGraph, overlapGraph)
	}
}

func TestLayout(t *testing.T) {
	reads := []string{"a_long_long",
		"_long_long_",
		"long_long_l",
		"ong_long_lo",
		"ng_long_lon",
		"g_long_long",
		"_long_long_",
		"long_long_t",
		"ong_long_ti",
		"ng_long_tim",
		"g_long_time",
		"_long_time_",
		"long_time_a",
		"ong_time_ag",
		"ng_time_ago"}
	min_overlap := 5
	overlap_graph := overlap(reads, min_overlap)

	layout := layout(overlap_graph)
	expectedLayout := []string{
		"a_long_long",
		"long_long_long",
		"long_long_time_ago",
	} // upto permutation
	sort.Strings(layout)
	sort.Strings(expectedLayout)
	if !reflect.DeepEqual(layout, expectedLayout) {
		t.Errorf("Expected layout %v, but got %v", expectedLayout, layout)
	}
}
