package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func overlap(reads []string, min_overlap int) map[string]map[string]int {
	overlap_graph := make(map[string]map[string]int)

	for _, read := range reads {
		for _, read2 := range reads {
			if read != read2 && strings.HasSuffix(read, read2[:min_overlap]) { // TODO: fix this
				overlap_graph[read][read2] = min_overlap // TODO: needs to be weighted by size of overlap
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

func main() {
	min_overlap := 100
	if len(os.Args) == 3 {
		var err error
		min_overlap, err = strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Usage: go run olc.go <fastq_filename> <min_overlap>")
			os.Exit(1)
		}
	}

	OLCAssembler(os.Args[1], min_overlap)
}
