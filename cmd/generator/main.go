package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	var filepath = flag.String("output", "input.txt", "Where the generated data will be saved.")
	var count = flag.Int("count", 1000, "How many tokens must be generated?.")
	var minLength = flag.Int("min-length", 128, "Minimum allowed token length.")
	var maxLength = flag.Int("max-length", 128, "Maximum allowed token length.")
	var equalPrefixLength = flag.Int("equal-prefix-length", 0, "All tokens will begin with the same prefix, and you can specify the length of this prefix.")
	var delimiter = flag.String("delimiter", "\n", "A character used to separate tokens.")
	var alphabetName = flag.String("alphabet", "lower", `You can specify one of the supported sets of characters used for generation:
binary - 01;
lower - abc..xyz;
upper - ABC...XYZ;
numbers - 012345689;
alnum - abc...xyz0123456789;
hex - 0123456789ABCDEF;
non-space - ASCII code in range (32, 128).

If the specified value doesn't match with any of the predefined alphabet names, this value will be used as the set of characters.`)

	flag.Parse()
	log.SetFlags(0)

	if len(*delimiter) != 1 {
		log.Fatalf("only one character can be specified as delimiter, but '%s' was given\n", *delimiter)
	}

	if *filepath == "" {
		log.Fatalf("empty filepath\n")
	}

	if *count <= 0 {
		log.Fatalf("count must be > 0\n")
	}

	if *maxLength <= 0 {
		log.Fatalf("max-length must be >= 0\n")
	}

	if *minLength <= 0 {
		log.Fatalf("min-length must be >= 0\n")
	}

	if *minLength > *maxLength {
		*minLength = *maxLength
	}

	if *equalPrefixLength > *minLength {
		*equalPrefixLength = *minLength
	}

	if *alphabetName == "" {
		log.Fatalf("empty alphabet\n")
	}

	alphabet, exists := alphabets[strings.ToLower(*alphabetName)]
	if !exists {
		alphabet = *alphabetName
	}

	file, err := os.OpenFile(*filepath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("cannot open file: %v\n", err)
	}
	defer func() {
		_ = file.Close()
	}()

	cfg := &config{
		output:            file,
		count:             *count,
		minLength:         *minLength,
		maxLength:         *maxLength,
		equalPrefixLength: *equalPrefixLength,
		delimiter:         (*delimiter)[0],
		alphabet:          []byte(alphabet),
	}

	// TODO: allow users to specify seed.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	gen := newGenerator(cfg, rnd)

	log.Printf(`configuration:

output file: %s
count: %d

length: [%d; %d]
equal prefix length: %d

delimiter: %#v
alphabet: %s`, *filepath, cfg.count, cfg.minLength, cfg.maxLength, cfg.equalPrefixLength, string(cfg.delimiter), string(cfg.alphabet))

	err = gen.generate()
	if err != nil {
		log.Fatalf("generation failed: %v\n", err)
	}

	log.Println("\n\nfinished")
}
