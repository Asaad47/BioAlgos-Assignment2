package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ContigNode struct {
	read     string
	outEdges map[string]int
	inEdges  map[string]int
}

func overlap(reads []string, min_overlap int) map[string]ContigNode {
	// naive implementation
	overlap_graph := make(map[string]ContigNode)
	for _, read := range reads {
		overlap_graph[read] = ContigNode{
			read:     read,
			outEdges: make(map[string]int),
			inEdges:  make(map[string]int),
		}
	}

	for _, read := range reads {
		for _, otherRead := range reads {
			if read == otherRead {
				continue
			}

			i := max(1+len(read)-len(otherRead), 1)

			for ; i < len(read)-min_overlap; i++ {
				if read[i:] == otherRead[:len(read)-i] {
					overlap_graph[read].outEdges[otherRead] = len(read) - i
					overlap_graph[otherRead].inEdges[read] = len(read) - i
					break
				}
			}
		}
	}

	return overlap_graph
}

func layout(overlap_graph map[string]ContigNode) []string {
	// first implementation: use reduction of inferrible edges in the overlap graph
	// source: https://www.cs.jhu.edu/~langmea/resources/lecture_notes/assembly_olc.pdf

	toBeRemoved := []string{}
	for read, edges := range overlap_graph {
		for otherRead, overlapLen := range edges.outEdges {
			for furtherRead, otherOverlapLen := range overlap_graph[otherRead].outEdges {
				if overlap_graph[read].outEdges[furtherRead] == overlapLen+otherOverlapLen-len(otherRead) {
					// remove the edge between read and furtherRead as it is inferrible
					toBeRemoved = append(toBeRemoved, read+" "+furtherRead)
				}
			}
		}
	}

	fmt.Println("Removing 1-hop inferrible edges: ", len(toBeRemoved))
	for _, edge := range toBeRemoved {
		parts := strings.Split(edge, " ")
		read := parts[0]
		otherRead := parts[1]
		delete(overlap_graph[read].outEdges, otherRead)
		delete(overlap_graph[otherRead].inEdges, read)
	}

	// identify nodes to be removed which are inferrible by 3 edges away
	toBeRemoved = []string{}
	for read1, edges1 := range overlap_graph {
		for read2, edge2Length := range edges1.outEdges {
			for read3, edge3Length := range overlap_graph[read2].outEdges {
				for read4, edge4Length := range overlap_graph[read3].outEdges {
					if overlap_graph[read1].outEdges[read4] == edge2Length+edge3Length+edge4Length-len(read2)-len(read3) {
						toBeRemoved = append(toBeRemoved, read1+" "+read4)
					}
				}
			}
		}
	}

	fmt.Println("Removing 2-hop inferrible edges: ", len(toBeRemoved))
	for _, edge := range toBeRemoved {
		parts := strings.Split(edge, " ")
		read := parts[0]
		otherRead := parts[1]
		delete(overlap_graph[read].outEdges, otherRead)
		delete(overlap_graph[otherRead].inEdges, read)
	}

	reads := []string{}
	for read := range overlap_graph {
		reads = append(reads, read)
	}

	numSingleOutEdges := 0
	// combine single-edge reads
	for _, read := range reads {
		for {
			if len(overlap_graph[read].outEdges) == 1 {
				numSingleOutEdges++
				otherRead := ""
				for otherRead = range overlap_graph[read].outEdges {
					break
				}
				if len(overlap_graph[otherRead].inEdges) == 1 {
					// combine the two reads
					newRead := read + otherRead[overlap_graph[read].outEdges[otherRead]:]
					newContigNode := ContigNode{
						read:     newRead,
						outEdges: overlap_graph[otherRead].outEdges,
						inEdges:  overlap_graph[read].inEdges,
					}

					for furtherRead, overlapLen := range overlap_graph[otherRead].outEdges {
						overlap_graph[furtherRead].inEdges[newRead] = overlapLen
						delete(overlap_graph[furtherRead].inEdges, otherRead)
					}

					for previousRead := range overlap_graph[read].inEdges {
						overlap_graph[previousRead].outEdges[newRead] = overlap_graph[read].inEdges[previousRead]
						delete(overlap_graph[previousRead].outEdges, read)
					}

					// update the overlap graph
					overlap_graph[newRead] = newContigNode
					delete(overlap_graph, read)
					delete(overlap_graph, otherRead)

					read = newRead
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	fmt.Println("Number of single-edge reads: ", numSingleOutEdges)

	contigs := []string{}
	for read := range overlap_graph {
		// if len(overlap_graph[read].outEdges) <= 1 && len(overlap_graph[read].inEdges) <= 1 {
		// 	contigs = append(contigs, read)
		// }
		contigs = append(contigs, read)
	}

	return contigs
}

func consensus(overlap_graph map[string]ContigNode) []string {
	contigs := []string{}
	visited := make(map[string]bool)

	numNoInEdges := 0
	// walk over the overlap graph greedily
	for read, readNode := range overlap_graph {
		if len(readNode.inEdges) == 0 {
			numNoInEdges++
			if len(readNode.outEdges) == 0 {
				contigs = append(contigs, read)
			} else {
				// start from this read and pick the longest outgoing edge
				max_length := 0
				max_node := ""
				for nextNode, length := range readNode.outEdges {
					if length > max_length {
						max_length = length
						max_node = nextNode
					}
				}

				visited[read] = true
				current := max_node
				path := current

				for {
					if visited[current] || current == "" {
						break
					}
					visited[current] = true
					path += current[max_length+1:]

					max_length := 0
					next := ""
					for k, length := range overlap_graph[current].outEdges {
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
	}

	fmt.Println("Number of reads with no in-edges: ", numNoInEdges)

	return contigs
}

func OLCAssembler(fastq_filename string, min_overlap int) {
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

	overlap_graph := overlap(reads, min_overlap)
	layout(overlap_graph)
	contigs := consensus(overlap_graph)

	// save the assembled contigs to a FASTA file
	fastaFile, err := os.Create(fastq_filename[:len(fastq_filename)-6] + "_olc_" + strconv.Itoa(min_overlap) + ".fasta")
	if err != nil {
		log.Fatalf("Failed to create FASTA file: %v", err)
	}
	defer fastaFile.Close()

	i := 0
	for _, contig := range contigs {
		// capitalize the contig
		contig = strings.ToUpper(contig)
		fastaFile.WriteString(">OLC_assembled_contig_" + strconv.Itoa(i) + "\n")
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

func DebugOLCAssembler() {

	// debugging
	// reads := []string{"a_long_long",
	// 	"_long_long_",
	// 	"long_long_l",
	// 	"ong_long_lo",
	// 	"ng_long_lon",
	// 	"g_long_long",
	// 	"_long_long_",
	// 	"long_long_t",
	// 	"ong_long_ti",
	// 	"ng_long_tim",
	// 	"g_long_time",
	// 	"_long_time_",
	// 	"long_time_a",
	// 	"ong_time_ag",
	// 	"ng_time_ago"}

	// fastq_filename := "../synthetic_dataset/reads/no_error_ont_hq_50x.fastq"
	// fastq_filename := "../synthetic_dataset/reads/no_error_reads_hiseq_5k.fastq"
	fastq_filename := "../synthetic_dataset/reads/ont_hq_50x.fastq"
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

	fmt.Println("------ reads ------", len(reads))

	min_overlap := 3
	overlap_graph := overlap(reads, min_overlap)

	// for read, edges := range overlap_graph {
	// 	fmt.Println(read)
	// 	for otherRead, overlapLen := range edges.outEdges {
	// 		fmt.Println("    ", otherRead, overlapLen)
	// 	}
	// }

	fmt.Println("------ after layout ------")
	layout(overlap_graph)
	fmt.Println("------ overlap graph ------")
	read_to_index := make(map[string]int)
	i := 0
	for read := range overlap_graph {
		// fmt.Println("** READ:", read)
		read_to_index[read] = i
		i++
		// for otherRead, overlapLen := range edges.outEdges {
		// 	fmt.Println("    OTHER READ:", otherRead)
		// 	fmt.Println("    -- OVERLAP LENGTH:", overlapLen)
		// }
	}
	fmt.Println("------ vertex to index mapping ------")
	fmt.Println("i:", i)
	contigs := consensus(overlap_graph)
	fmt.Println("------ contigs ------", len(contigs))
	// for _, contig := range contigs {
	// 	fmt.Println(contig)
	// }
	// for read, edges := range overlap_graph {
	// 	fmt.Println("** READ:", read_to_index[read])
	// 	for otherRead, overlapLen := range edges.outEdges {
	// 		fmt.Println("    OTHER READ:", read_to_index[otherRead])
	// 		fmt.Println("    -- OVERLAP LENGTH:", overlapLen)
	// 	}
	// }

	// readSizeCount := make(map[int]int)
	// for read := range overlap_graph {
	// 	readSizeCount[len(read)]++
	// }

	// for size, count := range readSizeCount {
	// 	fmt.Println("** READ SIZE:", size, "** COUNT:", count)
	// }
}
