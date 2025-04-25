package main

import "fmt"

type Node struct {
	Vertex    string
	OutDegree int
	InDegree  int
	Edges     map[string]int
}

func construct_de_bruijn_graph(reads []string, kmer_length int) map[string]Node {
	de_bruijn_graph := make(map[string]Node)

	for _, read := range reads {
		for i := 0; i <= len(read)-kmer_length; i++ {
			kmer := read[i : i+kmer_length]
			prefix := kmer[:kmer_length-1]
			suffix := kmer[1:]

			// Get the node, modify it, then put it back in the map
			prefixNode, ok := de_bruijn_graph[prefix]
			if !ok {
				prefixNode = Node{
					Vertex:    prefix,
					OutDegree: 0,
					InDegree:  0,
					Edges:     make(map[string]int),
				}
			}
			prefixNode.Edges[suffix]++
			prefixNode.OutDegree++
			de_bruijn_graph[prefix] = prefixNode

			// Same for suffix node
			suffixNode, ok := de_bruijn_graph[suffix]
			if !ok {
				suffixNode = Node{
					Vertex:    suffix,
					OutDegree: 0,
					InDegree:  0,
					Edges:     make(map[string]int),
				}
			}
			suffixNode.InDegree++
			de_bruijn_graph[suffix] = suffixNode
		}
	}

	return de_bruijn_graph
}

// check if given graph has a valid degree for an Eulerian path
func check_graph_degree(de_bruijn_graph map[string]Node) bool {
	found_out_more_than_in := 0 // out degree more than in degree by exactly 1
	found_in_more_than_out := 0 // in degree more than out degree by exactly 1

	for _, node := range de_bruijn_graph {
		switch node.OutDegree - node.InDegree {
		case 0:
			continue
		case 1:
			found_out_more_than_in++
		case -1:
			found_in_more_than_out++
		default:
			return false
		}
	}

	return found_out_more_than_in <= 1 && found_in_more_than_out <= 1
}

func find_eulerian_path(de_bruijn_graph map[string]Node) []string {
	// Hierholzer algorithm to find an Eulerian path in the graph

	// find a starting vertex
	// find a path from the starting vertex
	// if the path is not a cycle, find a cycle and add it to the path
	// repeat until no more edges left
	// return the path

	// find a starting vertex
	start_vertex := ""
	for vertex := range de_bruijn_graph {
		if de_bruijn_graph[vertex].OutDegree-de_bruijn_graph[vertex].InDegree == 1 {
			start_vertex = vertex
			break
		}
	}
	if start_vertex == "" {
		for vertex := range de_bruijn_graph {
			start_vertex = vertex
			break
		}
	}
	if start_vertex == "" { // no starting vertex found, should not happen
		return []string{}
	}

	// find a path from the starting vertex
	path := []string{start_vertex}
	short_path := []string{start_vertex}
	index_of_start_vertex := 0

	for len(de_bruijn_graph) > 0 {

		current_vertex := short_path[len(short_path)-1]
		next_vertex := ""
		for edge := range de_bruijn_graph[current_vertex].Edges {
			next_vertex = edge
			break
		}

		if next_vertex == "" {
			// got stuck, no more edges to follow
			if index_of_start_vertex == len(path)-1 {
				path = append(path[:index_of_start_vertex], short_path...)
			} else {
				// problem, should not happen
				fmt.Println("path: ", path)
				fmt.Println("short_path: ", short_path)
				fmt.Println("index_of_start_vertex: ", index_of_start_vertex)
				fmt.Println("de_bruijn_graph: ", de_bruijn_graph)
				return []string{}
			}

			// find a new starting vertex (should make cycle)
			for vertex := range de_bruijn_graph {
				if de_bruijn_graph[vertex].OutDegree > 0 {
					start_vertex = vertex
					for i, v := range path {
						if v == start_vertex {
							index_of_start_vertex = i
							break
						}
					}
					break
				}
			}

			short_path = []string{start_vertex}
			continue
		}

		if current_vertex == start_vertex && len(short_path) > 1 {
			// found a cycle
			short_path = append(short_path[:len(short_path)-1], path[index_of_start_vertex:]...)
			path = append(path[:index_of_start_vertex], short_path...)

			// find new starting vertex
			for vertex := range de_bruijn_graph {
				if de_bruijn_graph[vertex].OutDegree > 0 {
					start_vertex = vertex
					short_path = []string{start_vertex}
					for i, v := range path {
						if v == start_vertex {
							index_of_start_vertex = i
							break
						}
					}
					break
				}
			}
			continue
		}

		// Update the path
		short_path = append(short_path, next_vertex)

		// Update the graph
		// Get the current node and update it
		node := de_bruijn_graph[current_vertex]
		node.Edges[next_vertex]--
		node.OutDegree--
		de_bruijn_graph[current_vertex] = node

		// Get the next node and update it
		next_node := de_bruijn_graph[next_vertex]
		next_node.InDegree--
		de_bruijn_graph[next_vertex] = next_node

		// Check if we need to delete the edge
		if de_bruijn_graph[current_vertex].Edges[next_vertex] == 0 {
			delete(de_bruijn_graph[current_vertex].Edges, next_vertex)
		}
		if de_bruijn_graph[current_vertex].OutDegree == 0 && de_bruijn_graph[current_vertex].InDegree == 0 {
			delete(de_bruijn_graph, current_vertex)
		}
		if de_bruijn_graph[next_vertex].OutDegree == 0 && de_bruijn_graph[next_vertex].InDegree == 0 {
			delete(de_bruijn_graph, next_vertex)
		}
	}
	if index_of_start_vertex == len(path)-1 {
		path = append(path[:index_of_start_vertex], short_path...)
	} else {
		short_path = append(short_path[:len(short_path)-1], path[index_of_start_vertex:]...)
		path = append(path[:index_of_start_vertex], short_path...)
	}

	return path
}

func DBGAssembler(fastq_filename string) string {
	return ""
}
