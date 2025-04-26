package main

import (
	"reflect"
	"testing"
)

func TestConstructDeBruijnGraph(t *testing.T) {
	reads := []string{"ACGCG", "CGCGT"}
	kmerLength := 3
	expectedGraph := map[string]Node{
		"AC": {Vertex: "AC", OutDegree: 1, InDegree: 0, Edges: map[string]int{"CG": 1}},
		"CG": {Vertex: "CG", OutDegree: 3, InDegree: 3, Edges: map[string]int{"GC": 2, "GT": 1}},
		"GC": {Vertex: "GC", OutDegree: 2, InDegree: 2, Edges: map[string]int{"CG": 2}},
		"GT": {Vertex: "GT", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}

	graph := construct_de_bruijn_graph(reads, kmerLength)

	if !reflect.DeepEqual(graph, expectedGraph) {
		t.Errorf("Expected graph %v, but got %v", expectedGraph, graph)
	}
}

func TestCheckGraphDegree(t *testing.T) {
	graph := map[string]Node{
		"AC": {Vertex: "AC", OutDegree: 1, InDegree: 0, Edges: map[string]int{"CG": 1}},
		"CG": {Vertex: "CG", OutDegree: 3, InDegree: 3, Edges: map[string]int{"GC": 2, "GT": 1}},
		"GC": {Vertex: "GC", OutDegree: 2, InDegree: 2, Edges: map[string]int{"CG": 2}},
		"GT": {Vertex: "GT", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}

	if !check_graph_degree(graph) {
		t.Errorf("Expected graph to have a valid degree, but got invalid")
	}
}

func TestFindEulerianPath(t *testing.T) {
	graph := map[string]Node{
		"AC": {Vertex: "AC", OutDegree: 1, InDegree: 0, Edges: map[string]int{"CG": 1}},
		"CG": {Vertex: "CG", OutDegree: 3, InDegree: 3, Edges: map[string]int{"GC": 2, "GT": 1}},
		"GC": {Vertex: "GC", OutDegree: 2, InDegree: 2, Edges: map[string]int{"CG": 2}},
		"GT": {Vertex: "GT", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}
	expectedPath := []string{"AC", "CG", "GC", "CG", "GC", "CG", "GT"}

	path := find_eulerian_path(graph)

	if !reflect.DeepEqual(path, expectedPath) {
		t.Errorf("Expected path %v, but got %v", expectedPath, path)
	}
}

func TestGraphWithoutEulerianPath(t *testing.T) {
	graph := map[string]Node{
		"A": {Vertex: "A", OutDegree: 2, InDegree: 0, Edges: map[string]int{"B": 2}},
		"B": {Vertex: "B", OutDegree: 0, InDegree: 2, Edges: map[string]int{}},
		"C": {Vertex: "C", OutDegree: 1, InDegree: 0, Edges: map[string]int{"D": 1}},
		"D": {Vertex: "D", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}

	if check_graph_degree(graph) {
		t.Errorf("Expected graph to be invalid for Eulerian path, but it was valid")
	}
}

func TestEmptyGraph(t *testing.T) {
	graph := map[string]Node{}

	if !check_graph_degree(graph) {
		t.Errorf("Expected empty graph to be valid (trivial Eulerian path)")
	}

	path := find_eulerian_path(graph)
	if len(path) != 0 {
		t.Errorf("Expected empty path for empty graph, got %v", path)
	}
}

func TestSingleNodeSelfLoop(t *testing.T) {
	graph := map[string]Node{
		"A": {Vertex: "A", OutDegree: 1, InDegree: 1, Edges: map[string]int{"A": 1}},
	}

	if !check_graph_degree(graph) {
		t.Errorf("Expected graph with single self-loop to have a valid Eulerian path")
	}

	expectedPath := []string{"A", "A"}
	path := find_eulerian_path(graph)
	if !reflect.DeepEqual(path, expectedPath) {
		t.Errorf("Expected path %v, got %v", expectedPath, path)
	}
}

func TestDisconnectedGraph(t *testing.T) {
	graph := map[string]Node{
		"A": {Vertex: "A", OutDegree: 1, InDegree: 0, Edges: map[string]int{"B": 1}},
		"B": {Vertex: "B", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
		"C": {Vertex: "C", OutDegree: 1, InDegree: 0, Edges: map[string]int{"D": 1}},
		"D": {Vertex: "D", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}

	if check_graph_degree(graph) {
		t.Errorf("Expected disconnected graph to be invalid for Eulerian path")
	}
}

func TestFindAllEulerianPaths(t *testing.T) {
	// Test case 1: Graph with multiple valid paths
	graph1 := map[string]Node{
		"A": {Vertex: "A", OutDegree: 2, InDegree: 0, Edges: map[string]int{"B": 1, "C": 1}},
		"B": {Vertex: "B", OutDegree: 1, InDegree: 1, Edges: map[string]int{"D": 1}},
		"C": {Vertex: "C", OutDegree: 1, InDegree: 1, Edges: map[string]int{"D": 1}},
		"D": {Vertex: "D", OutDegree: 0, InDegree: 2, Edges: map[string]int{}},
	}
	expectedPaths1 := [][]string{
		{"A", "B", "D"},
		{"A", "C", "D"},
	}
	paths1 := find_all_eulerian_paths(graph1)
	if !reflect.DeepEqual(paths1, expectedPaths1) {
		t.Errorf("Expected paths %v, but got %v", expectedPaths1, paths1)
	}

	// Test case 2: Disconnected graph
	graph2 := map[string]Node{
		"A": {Vertex: "A", OutDegree: 1, InDegree: 0, Edges: map[string]int{"B": 1}},
		"B": {Vertex: "B", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
		"C": {Vertex: "C", OutDegree: 1, InDegree: 0, Edges: map[string]int{"D": 1}},
		"D": {Vertex: "D", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}
	expectedPaths2 := [][]string{
		{"A", "B"},
		{"C", "D"},
	}
	paths2 := find_all_eulerian_paths(graph2)
	if !reflect.DeepEqual(paths2, expectedPaths2) {
		t.Errorf("Expected paths %v, but got %v", expectedPaths2, paths2)
	}

	// Test case 3: Graph with a single valid path
	graph3 := map[string]Node{
		"A": {Vertex: "A", OutDegree: 1, InDegree: 0, Edges: map[string]int{"B": 1}},
		"B": {Vertex: "B", OutDegree: 1, InDegree: 1, Edges: map[string]int{"C": 1}},
		"C": {Vertex: "C", OutDegree: 0, InDegree: 1, Edges: map[string]int{}},
	}
	expectedPaths3 := [][]string{
		{"A", "B", "C"},
	}
	paths3 := find_all_eulerian_paths(graph3)
	if !reflect.DeepEqual(paths3, expectedPaths3) {
		t.Errorf("Expected paths %v, but got %v", expectedPaths3, paths3)
	}

	// Test case 4: Empty graph
	graph4 := map[string]Node{}
	paths4 := find_all_eulerian_paths(graph4)
	if len(paths4) != 0 {
		t.Errorf("Expected no paths for empty graph, but got %v", paths4)
	}

	// Test case 5: Graph with a cycle
	graph5 := map[string]Node{
		"A": {Vertex: "A", OutDegree: 1, InDegree: 1, Edges: map[string]int{"B": 1}},
		"B": {Vertex: "B", OutDegree: 1, InDegree: 1, Edges: map[string]int{"C": 1}},
		"C": {Vertex: "C", OutDegree: 1, InDegree: 1, Edges: map[string]int{"A": 1}},
	}
	expectedPaths5 := [][]string{
		{"A", "B", "C", "A"},
	}
	paths5 := find_all_eulerian_paths(graph5)
	if !reflect.DeepEqual(paths5, expectedPaths5) {
		t.Errorf("Expected paths %v, but got %v", expectedPaths5, paths5)
	}

	// Test case 6: Graph with multiple edges between same vertices
	graph6 := map[string]Node{
		"A": {Vertex: "A", OutDegree: 2, InDegree: 0, Edges: map[string]int{"B": 2}},
		"B": {Vertex: "B", OutDegree: 0, InDegree: 2, Edges: map[string]int{}},
	}
	expectedPaths6 := [][]string{
		{"A", "B"},
		{"A", "B"},
	}
	paths6 := find_all_eulerian_paths(graph6)
	if !reflect.DeepEqual(paths6, expectedPaths6) {
		t.Errorf("Expected paths %v, but got %v", expectedPaths6, paths6)
	}
}
