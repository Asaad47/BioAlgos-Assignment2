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
