package main

import (
	"bytes"
	"flag"
	"log"
	"strings"

	"github.com/lodthe/external-merge-sort/pkg/algo"
	"github.com/lodthe/external-merge-sort/pkg/config"
)

func main() {
	var blockSize = flag.Int("blocksize", 1024*1024, "Size of one block (in bytes).")
	var memoryLimit = flag.Int("memory", 512*1024*1024, "The algorithm will use at most O(memory) main memory.")
	var delimiter = flag.String("delimiter", "\n", "A character used to separate tokens.")
	var order = flag.String("order", "ASC", "Sort order. Supported values: ASC, DESC.")
	var inputFilepath = flag.String("input", "input.txt", "Input file path.")
	var outputFilepath = flag.String("output", "output.txt", "Output file path.")
	var tempDir = flag.String("tempdir", ".", "Where temporary files can be created. If you use /tmp, make sure there is enough space for two copies of the input file.")

	flag.Parse()
	log.SetFlags(0)

	if *blockSize < 0 {
		log.Fatalf("blocksize must be positive, but %d was given", *blockSize)
	}

	if *memoryLimit / *blockSize < 3 {
		log.Fatalf("'memory' must be at least three times larger than 'blocksize'")
	}

	if len(*delimiter) != 1 {
		log.Fatalf("only one character can be specified as delimiter, but %s was given", *delimiter)
	}

	cfg := &config.Config{
		BlockSize:   *blockSize,
		MemoryLimit: *memoryLimit,
		Delimiter:   (*delimiter)[0],
	}

	switch {
	case strings.EqualFold(*order, "ASC"):
		cfg.Less = func(a, b []byte) bool {
			return bytes.Compare(a, b) < 0
		}

	case strings.EqualFold(*order, "DESC"):
		cfg.Less = func(a, b []byte) bool {
			return bytes.Compare(a, b) > 0
		}

	default:
		log.Fatalf("only ASC and DESC orders are supported, but %s was given", *order)
	}

	msort := algo.NewExternalMergeSort(cfg)
	err := msort.Sort(*inputFilepath, *outputFilepath, *tempDir)
	if err != nil {
		log.Fatalf("sort failed: %v\n", err)
	}
}
