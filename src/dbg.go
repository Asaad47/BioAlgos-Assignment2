package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct { // Expanded kmer node. Represents a kmer and the edges leading to it and from it
	kmer     string
	OutEdges map[string]int
	InEdges  map[string]int
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
					kmer:     prefix,
					OutEdges: make(map[string]int),
					InEdges:  make(map[string]int),
				}
			}
			prefixNode.OutEdges[suffix]++
			de_bruijn_graph[prefix] = prefixNode

			// Same for suffix node
			suffixNode, ok := de_bruijn_graph[suffix]
			if !ok {
				suffixNode = Node{
					kmer:     suffix,
					OutEdges: make(map[string]int),
					InEdges:  make(map[string]int),
				}
			}
			suffixNode.InEdges[prefix]++
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
		out_degree := 0
		in_degree := 0
		for _, out_edge := range node.OutEdges {
			out_degree += out_edge
		}
		for _, in_edge := range node.InEdges {
			in_degree += in_edge
		}
		switch out_degree - in_degree {
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

// Assumes the graph passes the degree check, which implies the graph is nice
// (i.e. collection of cycles or with a single starting vertex)
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
		out_degree := 0
		in_degree := 0
		for _, out_edge := range de_bruijn_graph[vertex].OutEdges {
			out_degree += out_edge
		}
		for _, in_edge := range de_bruijn_graph[vertex].InEdges {
			in_degree += in_edge
		}
		if out_degree-in_degree == 1 {
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
		for edge := range de_bruijn_graph[current_vertex].OutEdges {
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
				out_degree := 0
				for _, out_edge := range de_bruijn_graph[vertex].OutEdges {
					out_degree += out_edge
				}
				if out_degree > 0 {
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
				out_degree := 0
				for _, out_edge := range de_bruijn_graph[vertex].OutEdges {
					out_degree += out_edge
				}
				if out_degree > 0 {
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
		node.OutEdges[next_vertex]--
		de_bruijn_graph[current_vertex] = node

		// Get the next node and update it
		next_node := de_bruijn_graph[next_vertex]
		next_node.InEdges[current_vertex]--
		de_bruijn_graph[next_vertex] = next_node

		// Check if we need to delete the edge
		if de_bruijn_graph[current_vertex].OutEdges[next_vertex] == 0 {
			delete(de_bruijn_graph[current_vertex].OutEdges, next_vertex)
		}
		if len(de_bruijn_graph[current_vertex].OutEdges) == 0 && len(de_bruijn_graph[current_vertex].InEdges) == 0 {
			delete(de_bruijn_graph, current_vertex)
		}
		if len(de_bruijn_graph[next_vertex].OutEdges) == 0 && len(de_bruijn_graph[next_vertex].InEdges) == 0 {
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

// Doesn't assume the graph passes the degree check. Can be any graph.
func find_all_eulerian_paths(de_bruijn_graph map[string]Node) [][]string {
	// find a starting vertex
	starting_vertices := []string{}
	for vertex := range de_bruijn_graph {
		out_degree := 0
		in_degree := 0
		for _, out_edge := range de_bruijn_graph[vertex].OutEdges {
			out_degree += out_edge
		}
		for _, in_edge := range de_bruijn_graph[vertex].InEdges {
			in_degree += in_edge
		}
		if out_degree-in_degree >= 1 {
			starting_vertices = append(starting_vertices, vertex)
		}
	}
	if len(starting_vertices) == 0 {
		// no starting vertex found, pick any vertex
		for vertex := range de_bruijn_graph {
			starting_vertices = append(starting_vertices, vertex)
		}
	}

	paths := [][]string{}
	path := []string{}
	visited := make(map[string]bool)

	for len(starting_vertices) > 0 {
		starting_vertex := starting_vertices[0]

		stack := []string{starting_vertex}

		for len(stack) > 0 {
			current_vertex := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if visited[current_vertex] {
				// found a cycle
				// TODO: add the cycle to the paths
				path = append(path, current_vertex)
				if len(path) > 1 {
					paths = append(paths, path)
				}
				path = []string{}
				visited = make(map[string]bool)
				continue
			}

			visited[current_vertex] = true
			path = append(path, current_vertex)
			next_vertex := ""
			for edge := range de_bruijn_graph[current_vertex].OutEdges {
				next_vertex = edge
				break
			}

			if next_vertex == "" {
				// got stuck, no more edges to follow
				// TODO: add the path to the paths
				if len(path) > 1 {
					paths = append(paths, path)
				}
				path = []string{}
				visited = make(map[string]bool)
				continue
			}

			stack = append(stack, next_vertex)

			// remove the edge from the graph
			de_bruijn_graph[current_vertex].OutEdges[next_vertex]--

			// Update the vertex properties by creating a new vertex with updated values
			vertex := de_bruijn_graph[current_vertex]
			vertex.OutEdges[next_vertex]--
			de_bruijn_graph[current_vertex] = vertex

			// Update the next vertex properties
			nextVertex := de_bruijn_graph[next_vertex]
			nextVertex.InEdges[current_vertex]--
			de_bruijn_graph[next_vertex] = nextVertex

			if de_bruijn_graph[current_vertex].OutEdges[next_vertex] == 0 {
				delete(de_bruijn_graph[current_vertex].OutEdges, next_vertex)
			}
			if len(de_bruijn_graph[current_vertex].OutEdges) == 0 && len(de_bruijn_graph[current_vertex].InEdges) == 0 {
				delete(de_bruijn_graph, current_vertex)
			}
			if len(de_bruijn_graph[next_vertex].OutEdges) == 0 && len(de_bruijn_graph[next_vertex].InEdges) == 0 {
				delete(de_bruijn_graph, next_vertex)
			}
		}

		starting_vertex = starting_vertices[0]
		start_in_degree_after := de_bruijn_graph[starting_vertex].InEdges
		start_out_degree_after := de_bruijn_graph[starting_vertex].OutEdges
		if len(start_in_degree_after) == 0 && len(start_out_degree_after) == 0 {
			starting_vertices = starting_vertices[1:]
		} else {
			// move the starting vertex to the end of the list
			starting_vertices = append(starting_vertices[1:], starting_vertices[0])
		}
	}

	return paths
}

func reduce_graph(de_bruijn_graph map[string]Node) map[string]Node {
	// reduce nodes that have single in edge and single out edge

	starting_vertices := []string{}
	for vertex := range de_bruijn_graph {
		if len(de_bruijn_graph[vertex].InEdges) != 1 {
			starting_vertices = append(starting_vertices, vertex)
		} else {
			previous_vertex := ""
			for edge := range de_bruijn_graph[vertex].InEdges {
				previous_vertex = edge
				break
			}
			if len(de_bruijn_graph[previous_vertex].OutEdges) != 1 {
				starting_vertices = append(starting_vertices, vertex)
			}
		}
	}
	for _, vertex := range starting_vertices {
		for {
			if len(de_bruijn_graph[vertex].OutEdges) == 1 {
				otherVertex := ""
				for edge := range de_bruijn_graph[vertex].OutEdges {
					otherVertex = edge
					break
				}
				if len(de_bruijn_graph[otherVertex].InEdges) == 1 {
					if len(de_bruijn_graph[otherVertex].OutEdges) == 1 {
						// reduce nodes vertex, otherVertex, furtherVertex to newFirstVertex and newSecondVertex
						furtherVertex := ""
						for edge := range de_bruijn_graph[otherVertex].OutEdges {
							furtherVertex = edge
							break
						}
						if len(de_bruijn_graph[furtherVertex].InEdges) > 1 {
							break
						}

						newFirstVertex := vertex + otherVertex[len(otherVertex)-1:] // this assumes using starting vertices as defined above
						newFirstVertexNode := Node{
							kmer:     newFirstVertex,
							OutEdges: make(map[string]int),
							InEdges:  de_bruijn_graph[vertex].InEdges,
						}

						newSecondVertex := otherVertex + furtherVertex[len(furtherVertex)-1:] // this assumes using starting vertices as defined above
						newSecondVertexNode := Node{
							kmer:     newSecondVertex,
							OutEdges: de_bruijn_graph[furtherVertex].OutEdges,
							InEdges:  make(map[string]int),
						}

						occurances := min(de_bruijn_graph[vertex].OutEdges[otherVertex], de_bruijn_graph[otherVertex].OutEdges[furtherVertex])
						newFirstVertexNode.OutEdges[newSecondVertex] = occurances
						newSecondVertexNode.InEdges[newFirstVertex] = occurances

						for previousVertex := range de_bruijn_graph[vertex].InEdges {
							de_bruijn_graph[previousVertex].OutEdges[newFirstVertex] = de_bruijn_graph[previousVertex].OutEdges[vertex]
							delete(de_bruijn_graph[previousVertex].OutEdges, vertex)
						}
						for nextVertex := range de_bruijn_graph[furtherVertex].OutEdges {
							de_bruijn_graph[nextVertex].InEdges[newSecondVertex] = de_bruijn_graph[nextVertex].InEdges[furtherVertex]
							delete(de_bruijn_graph[nextVertex].InEdges, furtherVertex)
						}

						de_bruijn_graph[newFirstVertex] = newFirstVertexNode
						de_bruijn_graph[newSecondVertex] = newSecondVertexNode
						delete(de_bruijn_graph, vertex)
						delete(de_bruijn_graph, otherVertex)
						delete(de_bruijn_graph, furtherVertex)

						vertex = newFirstVertex

					} else if len(de_bruijn_graph[otherVertex].OutEdges) == 0 {
						// reached end of the path
						// reduce vertex, otherVertex to newVertex
						newVertex := vertex[:1] + otherVertex
						newVertexNode := Node{
							kmer:     newVertex,
							OutEdges: de_bruijn_graph[otherVertex].OutEdges,
							InEdges:  de_bruijn_graph[vertex].InEdges,
						}
						for previousVertex := range de_bruijn_graph[vertex].InEdges {
							de_bruijn_graph[previousVertex].OutEdges[newVertex] = de_bruijn_graph[previousVertex].OutEdges[vertex]
							delete(de_bruijn_graph[previousVertex].OutEdges, vertex)
						}
						de_bruijn_graph[newVertex] = newVertexNode
						delete(de_bruijn_graph, vertex)
						delete(de_bruijn_graph, otherVertex)
						break
					} else {
						break
					}
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	return de_bruijn_graph
}

func DBGAssembler(fastq_filename string, kmer_length int) string {
	fastqFile, err := os.Open(fastq_filename)
	if err != nil {
		log.Fatalf("Failed to open FASTQ file: %v", err)
	}
	defer fastqFile.Close()

	scanner := bufio.NewScanner(fastqFile)

	reads := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@") || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "I") {
			continue
		}

		read := strings.ToLower(strings.TrimSpace(line))
		reads = append(reads, read)
	}

	de_bruijn_graph := construct_de_bruijn_graph(reads, kmer_length)

	paths := find_all_eulerian_paths(de_bruijn_graph)

	// assemble the reads from the path
	longest_assembled_read := ""
	for _, path := range paths {
		assembled_read := ""
		for i, vertex := range path {
			if i == 0 {
				assembled_read += vertex
			} else {
				assembled_read += vertex[len(vertex)-1:]
			}
		}
		if len(assembled_read) > len(longest_assembled_read) {
			longest_assembled_read = assembled_read
		}
	}

	// save the assembled read to a FASTA file
	fastaFile, err := os.Create(fastq_filename[:len(fastq_filename)-6] + "_dbg.fasta")
	if err != nil {
		log.Fatalf("Failed to create FASTA file: %v", err)
	}
	defer fastaFile.Close()

	fastaFile.WriteString(">DBG_assembled_read\n")
	// write on lines with 60 characters
	for i := 0; i < len(longest_assembled_read); i += 60 {
		if i+60 > len(longest_assembled_read) {
			fastaFile.WriteString(strings.ToUpper(longest_assembled_read[i:]) + "\n")
		} else {
			fastaFile.WriteString(strings.ToUpper(longest_assembled_read[i:i+60]) + "\n")
		}
	}

	return longest_assembled_read
}

func main() {
	// kmer_length := 40
	// if len(os.Args) == 3 {
	// 	var err error
	// 	kmer_length, err = strconv.Atoi(os.Args[2])
	// 	if err != nil {
	// 		fmt.Println("Usage: go run dbg.go <fastq_filename> <kmer_length>")
	// 		os.Exit(1)
	// 	}
	// }

	// DBGAssembler(os.Args[1], kmer_length)

	// debugging
	fastqFile, err := os.Open("../toy_dataset/reads_r.fastq")
	if err != nil {
		log.Fatalf("Failed to open FASTQ file: %v", err)
	}
	defer fastqFile.Close()

	scanner := bufio.NewScanner(fastqFile)

	reads := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "@") || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "I") {
			continue
		}

		read := strings.ToLower(strings.TrimSpace(line))
		reads = append(reads, read)
	}

	de_bruijn_graph := construct_de_bruijn_graph(reads, 45)
	for k, v := range de_bruijn_graph {
		fmt.Println(k, v)
	}

	// paths := find_all_eulerian_paths(de_bruijn_graph)
	// longest_path := ""
	// for _, path := range paths {
	// 	// combine the path into a single string
	// 	combined_path := ""
	// 	for i, vertex := range path {
	// 		if i == 0 {
	// 			combined_path += vertex
	// 		} else {
	// 			combined_path += vertex[len(vertex)-1:]
	// 		}
	// 	}
	// 	if len(combined_path) > len(longest_path) {
	// 		longest_path = combined_path
	// 	}
	// }
	// fmt.Println(longest_path)

	fmt.Println("------ after reduction ------")
	reduced_de_bruijn_graph := reduce_graph(de_bruijn_graph)
	i := 0
	vertex_to_index := make(map[string]int)
	for k, v := range reduced_de_bruijn_graph {
		vertex_to_index[k] = i
		i++
		fmt.Println("** vertex:", k)
		for edge, count := range v.OutEdges {
			fmt.Println("    edge:", edge)
			fmt.Println("    -- count:", count)
		}
	}
	fmt.Println("------ vertex to index ------")
	for k, v := range reduced_de_bruijn_graph {
		fmt.Println("** vertex:", vertex_to_index[k])
		for edge, count := range v.OutEdges {
			fmt.Println("    edge:", vertex_to_index[edge])
			fmt.Println("    -- count:", count)
		}
	}

}
