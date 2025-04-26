package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// // Node represents a node in the suffix tree
// type TreeNode struct {
// 	children   map[rune]*TreeNode
// 	start      int
// 	end        *int // pointer to allow for leaf nodes to share end position
// 	suffixLink *TreeNode
// 	read       string // stores which read this suffix belongs to
// }

// // SuffixTree represents the suffix tree data structure
// type SuffixTree struct {
// 	root         *TreeNode
// 	activeNode   *TreeNode
// 	activeEdge   rune
// 	activeLength int
// 	remaining    int
// 	end          int
// }

// // NewSuffixTree creates a new suffix tree
// func NewSuffixTree() *SuffixTree {
// 	root := &TreeNode{
// 		children: make(map[rune]*TreeNode),
// 		start:    -1,
// 		end:      new(int),
// 	}
// 	*root.end = -1
// 	return &SuffixTree{
// 		root:         root,
// 		activeNode:   root,
// 		activeEdge:   0,
// 		activeLength: 0,
// 		remaining:    0,
// 		end:          -1,
// 	}
// }

// // Insert adds a string to the tree using Ukkonen's algorithm
// func (st *SuffixTree) Insert(s string) {
// 	s = s + "$"
// 	st.end = -1
// 	st.remaining = 0
// 	st.activeNode = st.root
// 	st.activeLength = 0

// 	for i := 0; i < len(s); i++ {
// 		st.end++
// 		st.remaining++
// 		var lastNewNode *TreeNode

// 		for st.remaining > 0 {
// 			if st.activeLength == 0 {
// 				st.activeEdge = rune(s[i])
// 			}

// 			if child, exists := st.activeNode.children[st.activeEdge]; exists {
// 				if st.walkDown(child) {
// 					continue
// 				}

// 				if s[child.start+st.activeLength] == s[i] {
// 					st.activeLength++
// 					break
// 				}

// 				// Rule 2: Split edge
// 				splitNode := &TreeNode{
// 					children: make(map[rune]*TreeNode),
// 					start:    child.start,
// 					end:      new(int),
// 					read:     s,
// 				}
// 				*splitNode.end = child.start + st.activeLength - 1

// 				child.start = child.start + st.activeLength
// 				splitNode.children[rune(s[child.start])] = child

// 				st.activeNode.children[st.activeEdge] = splitNode

// 				// Create new leaf
// 				leaf := &TreeNode{
// 					children: make(map[rune]*TreeNode),
// 					start:    i,
// 					end:      new(int),
// 					read:     s,
// 				}
// 				*leaf.end = st.end
// 				splitNode.children[rune(s[i])] = leaf

// 				if lastNewNode != nil {
// 					lastNewNode.suffixLink = splitNode
// 				}
// 				lastNewNode = splitNode
// 			} else {
// 				// Rule 1: Create new leaf
// 				leaf := &TreeNode{
// 					children: make(map[rune]*TreeNode),
// 					start:    i,
// 					end:      new(int),
// 					read:     s,
// 				}
// 				*leaf.end = st.end
// 				st.activeNode.children[st.activeEdge] = leaf

// 				if lastNewNode != nil {
// 					lastNewNode.suffixLink = st.activeNode
// 				}
// 				lastNewNode = nil
// 			}

// 			st.remaining--

// 			if st.activeNode == st.root && st.activeLength > 0 {
// 				st.activeLength--
// 				st.activeEdge = rune(s[i-st.remaining+1])
// 			} else if st.activeNode != st.root {
// 				if st.activeNode.suffixLink == nil {
// 					st.activeNode.suffixLink = st.root
// 				}
// 				st.activeNode = st.activeNode.suffixLink
// 			}
// 		}
// 	}
// }

// // walkDown returns true if we need to continue walking down
// func (st *SuffixTree) walkDown(node *TreeNode) bool {
// 	if node == nil || node.end == nil {
// 		return false
// 	}
// 	edgeLength := *node.end - node.start + 1
// 	if st.activeLength >= edgeLength {
// 		st.activeNode = node
// 		st.activeLength -= edgeLength
// 		st.activeEdge = rune(st.end - st.remaining + 1)
// 		return true
// 	}
// 	return false
// }

// // FindOverlaps finds all overlaps between the given string and strings in the tree
// func (st *SuffixTree) FindOverlaps(s string, minOverlap int) map[string]int {
// 	overlaps := make(map[string]int)
// 	originalS := s // Store original string without '$'
// 	s = s + "$"    // Add '$' for tree traversal

// 	current := st.root
// 	length := 0
// 	pos := 0

// 	for pos < len(s) {
// 		// fmt.Println("pos:", pos)
// 		// fmt.Println("length:", length)
// 		if length == 0 {
// 			// Start from root
// 			current = st.root
// 		}

// 		// Try to find the next character in the current node's children
// 		if child, exists := current.children[rune(s[pos])]; exists {
// 			// Calculate how much of this edge we can match
// 			edgeLength := *child.end - child.start + 1
// 			matchLength := 0

// 			// Match as much as possible along this edge
// 			for matchLength < edgeLength && pos+matchLength < len(s) {
// 				if s[pos+matchLength] != s[child.start+matchLength] {
// 					break
// 				}
// 				matchLength++
// 			}

// 			// Update position and length
// 			pos += matchLength
// 			length += matchLength

// 			// If we've matched the entire edge, move to the child node
// 			if matchLength == edgeLength && matchLength != 0 {
// 				current = child
// 			} else {
// 				// Partial match, we can't go further
// 				// TODO: is this correct?
// 				break
// 			}
// 		} else {
// 			// No match found
// 			// TODO: is this correct?
// 			break
// 		}
// 	}
// 	// If we found a sufficient overlap
// 	if length >= minOverlap {
// 		// Collect all reads from the subtree
// 		st.collectReads(current, overlaps, length, originalS)
// 	}

// 	return overlaps
// }

// // collectReads recursively collects all reads from a subtree
// func (st *SuffixTree) collectReads(node *TreeNode, overlaps map[string]int, length int, originalS string) {
// 	if node.read != "" {
// 		// Remove '$' from the stored read
// 		originalRead := node.read[:len(node.read)-1]
// 		// Only add if it's a different read and the overlap is valid
// 		if originalRead != originalS {
// 			overlaps[originalRead] = length
// 		}
// 	}

// 	for _, child := range node.children {
// 		st.collectReads(child, overlaps, length, originalS)
// 	}
// }

// func overlap(reads []string, min_overlap int) map[string]map[string]int {
// 	overlap_graph := make(map[string]map[string]int)

// 	// Initialize the overlap graph
// 	for _, read := range reads {
// 		overlap_graph[read] = make(map[string]int)
// 	}

// 	// Build the suffix tree
// 	st := NewSuffixTree()
// 	for _, read := range reads {
// 		st.Insert(read)
// 	}

// 	fmt.Println("suffix tree:", st)
// 	fmt.Println("root:", st.root)
// 	for _, child := range st.root.children {
// 		fmt.Println("child:", child)
// 	}

// 	// Find overlaps for each read
// 	for _, read := range reads {
// 		overlaps := st.FindOverlaps(read, min_overlap)
// 		for otherRead, overlapLen := range overlaps {
// 			overlap_graph[read][otherRead] = overlapLen
// 		}
// 	}

// 	return overlap_graph
// }

func overlap(reads []string, min_overlap int) map[string]map[string]int {
	// naive implementation
	overlap_graph := make(map[string]map[string]int)
	for _, read := range reads {
		overlap_graph[read] = make(map[string]int)
	}

	for _, read := range reads {
		for _, otherRead := range reads {
			if read == otherRead {
				continue
			}

			for i := 0; i < len(read)-min_overlap; i++ {
				if read[i:] == otherRead[:len(read)-i] {
					overlap_graph[read][otherRead] = len(read) - i
					break
				}
			}
		}
	}

	return overlap_graph
}

func layout(overlap_graph map[string]map[string]int) []string {
	return []string{}
}

func consensus(read_layout []string) string {
	return ""
}

func OLCAssembler(fastq_filename string, min_overlap int) string {
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

	overlap_graph := overlap(reads, min_overlap)
	read_layout := layout(overlap_graph)
	consensus := consensus(read_layout)

	// save the consensus to a FASTA file
	fastaFile, err := os.Create(fastq_filename[:len(fastq_filename)-6] + "_olc.fasta")
	if err != nil {
		log.Fatalf("Failed to create FASTA file: %v", err)
	}
	defer fastaFile.Close()

	fastaFile.WriteString(">OLC_consensus\n")
	fastaFile.WriteString(consensus)

	return consensus
}

// func main() {
// 	min_overlap := 100
// 	if len(os.Args) == 3 {
// 		var err error
// 		min_overlap, err = strconv.Atoi(os.Args[2])
// 		if err != nil {
// 			fmt.Println("Usage: go run olc.go <fastq_filename> <min_overlap>")
// 			os.Exit(1)
// 		}
// 	}

// 	OLCAssembler(os.Args[1], min_overlap)
// }
