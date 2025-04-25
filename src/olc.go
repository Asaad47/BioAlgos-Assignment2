package main

import (
	"bufio"
	"log"
	"os"
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

	return consensus
}

func main() {

}
