package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct { // Expanded kmer node. Represents a kmer and the edges leading to it and from it
	kmer     string
	OutEdges map[string]int
	InEdges  map[string]int
}

func constructDeBruijnGraph(reads []string, kmer_length int) map[string]Node {
	deBruijnGraph := make(map[string]Node)

	for _, read := range reads {
		for i := 0; i <= len(read)-kmer_length; i++ {
			kmer := read[i : i+kmer_length]
			prefix := kmer[:kmer_length-1]
			suffix := kmer[1:]

			// Get the node, modify it, then put it back in the map
			prefixNode, ok := deBruijnGraph[prefix]
			if !ok {
				prefixNode = Node{
					kmer:     prefix,
					OutEdges: make(map[string]int),
					InEdges:  make(map[string]int),
				}
			}
			prefixNode.OutEdges[suffix]++
			deBruijnGraph[prefix] = prefixNode

			// Same for suffix node
			suffixNode, ok := deBruijnGraph[suffix]
			if !ok {
				suffixNode = Node{
					kmer:     suffix,
					OutEdges: make(map[string]int),
					InEdges:  make(map[string]int),
				}
			}
			suffixNode.InEdges[prefix]++
			deBruijnGraph[suffix] = suffixNode
		}
	}

	return deBruijnGraph
}

func walkGraph(deBruijnGraph map[string]Node) []string {
	contigs := []string{}
	visited := make(map[string]bool)

	for nodeID, node := range deBruijnGraph {
		// start only from nodes that don't have an in-edge
		if len(node.InEdges) == 0 && len(node.OutEdges) > 0 {
			// pick the longest outgoing edge
			max_length := 0
			max_node := ""
			for nextNode, length := range node.OutEdges {
				if length > max_length {
					max_length = length
					max_node = nextNode
				}
			}

			visited[nodeID] = true
			current := max_node
			path := current

			for {
				if visited[current] || current == "" {
					break
				}
				visited[current] = true
				path += current[len(current)-1:] // append last base of kmer

				max_length := 0
				next := ""
				for k, length := range deBruijnGraph[current].OutEdges {
					if length > max_length {
						max_length = length
						next = k
					}
				}
				current = next
			}

			contigs = append(contigs, path)
		}
	}

	return contigs
}

func exportToGFA(graph map[string]Node, filename string) {
	gfaFile, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create GFA file: %v", err)
	}
	defer gfaFile.Close()

	// Write segments (S lines)
	for kmer := range graph {
		gfaFile.WriteString("S\t" + kmer + "\t*\n")
	}

	// Write links (L lines)
	for from, node := range graph {
		for to := range node.OutEdges {
			// All overlaps are of length k-1
			gfaFile.WriteString("L\t" + from + "\t+\t" + to + "\t+\t" + strconv.Itoa(len(from)-1) + "M\n")
		}
	}
}

func DBGAssembler(fastq_filename string, kmer_length int) {
	fastqFile, err := os.Open(fastq_filename)
	if err != nil {
		log.Fatalf("Failed to open FASTQ file: %v", err)
	}
	defer fastqFile.Close()

	scanner := bufio.NewScanner(fastqFile)

	reads := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "A") && !strings.HasPrefix(line, "T") && !strings.HasPrefix(line, "C") && !strings.HasPrefix(line, "G") {
			continue
		}

		read := strings.ToLower(strings.TrimSpace(line))
		reads = append(reads, read)
	}

	deBruijnGraph := constructDeBruijnGraph(reads, kmer_length)

	// exportToGFA(deBruijnGraph, fastq_filename[:len(fastq_filename)-6]+"_dbg_k_"+strconv.Itoa(kmer_length)+".gfa")

	// reducedGraph := reduceGraph(deBruijnGraph)
	// paths := walkReducedGraph(reducedGraph)

	contigs := walkGraph(deBruijnGraph)

	// save the assembled contigs to a FASTA file
	fastaFile, err := os.Create(fastq_filename[:len(fastq_filename)-6] + "_dbg_k_" + strconv.Itoa(kmer_length) + ".fasta")
	if err != nil {
		log.Fatalf("Failed to create FASTA file: %v", err)
	}
	defer fastaFile.Close()

	i := 0
	for _, contig := range contigs {
		// capitalize the contig
		contig = strings.ToUpper(contig)
		fastaFile.WriteString(">DBG_assembled_contig_" + strconv.Itoa(i) + "\n")
		// write on lines with 60 characters
		for j := 0; j < len(contig); j += 60 {
			if j+60 > len(contig) {
				fastaFile.WriteString(contig[j:] + "\n")
			} else {
				fastaFile.WriteString(contig[j:j+60] + "\n")
			}
		}
		i++
	}
}

func DebugDBGAssembler() {
	kmer_length := 45
	fastq_filename := "../synthetic_dataset/reads/reads_hiseq_5k.fastq"

	// DBGAssembler(fastq_filename, kmer_length)

	fastqFile, err := os.Open(fastq_filename)
	if err != nil {
		log.Fatalf("Failed to open FASTQ file: %v", err)
	}
	defer fastqFile.Close()

	scanner := bufio.NewScanner(fastqFile)

	reads := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "A") && !strings.HasPrefix(line, "T") && !strings.HasPrefix(line, "C") && !strings.HasPrefix(line, "G") {
			continue
		}

		read := strings.ToLower(strings.TrimSpace(line))
		reads = append(reads, read)
	}

	deBruijnGraph := constructDeBruijnGraph(reads, kmer_length)

	contigs := walkGraph(deBruijnGraph)

	for i, contig := range contigs {
		fmt.Println("contig ", i, ": ", contig)
	}
}
