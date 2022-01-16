package main

import (
	"bytes"
	"flag"
	"log"
	"strings"
)

func main() {
	var delimiter = flag.String("delimiter", "\n", "A character used to separate tokens.")
	var order = flag.String("order", "ASC", "Sort order. Supported values: ASC, DESC.")
	var inputFilepath = flag.String("input", "input.txt", "Input file path.")
	var sortedFilepath = flag.String("output", "output.txt", "Output file path.")

	flag.Parse()
	log.SetFlags(0)

	if len(*delimiter) != 1 {
		log.Fatalf("only one character can be specified as delimiter, but %s was given", *delimiter)
	}

	var less func(a, b []byte) bool
	switch {
	case strings.EqualFold(*order, "ASC"):
		less = func(a, b []byte) bool {
			return bytes.Compare(a, b) < 0
		}

	case strings.EqualFold(*order, "DESC"):
		less = func(a, b []byte) bool {
			return bytes.Compare(a, b) > 0
		}

	default:
		log.Fatalf("only ASC and DESC orders are supported, but %s was given", *order)
	}

	delim := (*delimiter)[0]
	inputTokenCount, inputHash, err := parseFile(*inputFilepath, nil, delim)
	if err != nil {
		log.Fatalf("parsing input file failed: %v\n", err)
	}

	sortedTokenCount, sortedHash, err := parseFile(*sortedFilepath, less, delim)
	if err != nil {
		log.Fatalf("parsing sorted file failed: %v\n", err)
	}

	if inputTokenCount != sortedTokenCount {
		log.Fatalf("token count mismatch: input file has %d, sorted file has %d\n", inputTokenCount, sortedTokenCount)
	}

	if inputHash != sortedHash {
		log.Fatalf("hash mismatch: input file has %d, sorted file has %d\n", inputHash, sortedHash)
	}

	log.Println("OK")
}
