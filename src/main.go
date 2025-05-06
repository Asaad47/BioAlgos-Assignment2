package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <assembler_type> <fastq_filename> <parameter>")
		fmt.Println("assembler_type: olc or dbg")
		fmt.Println("parameter: min_overlap for olc, kmer_length for dbg")
		fmt.Println("OR for debugging, go run . <assembler_type> debug")
		os.Exit(1)
	}

	assemblerType := strings.ToLower(os.Args[1])
	fastq_filename := os.Args[2]

	if fastq_filename == "debug" {
		if assemblerType == "olc" {
			DebugOLCAssembler()
		} else {
			DebugDBGAssembler()
		}
		os.Exit(0)
	}

	length_value, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Usage: go run . <assembler_type> <fastq_filename> <parameter>")
		fmt.Println("parameter: min_overlap for olc, kmer_length for dbg")
		os.Exit(1)
	}

	switch assemblerType {
	case "olc":
		OLCAssembler(fastq_filename, length_value)
	case "dbg":
		DBGAssembler(fastq_filename, length_value)
	default:
		fmt.Println("Unknown assembler type. Use 'olc' or 'dbg'")
		os.Exit(1)
	}
}
