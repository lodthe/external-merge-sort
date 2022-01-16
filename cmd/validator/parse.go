package main

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/lodthe/external-merge-sort/pkg/hash"
	"github.com/pkg/errors"
)

const bufferSize = 128 * 1024

// parseFile reads file content and counts tokens and their hash.
// If less is provided, it also checks if less(str[i], str[i + 1]) is true for each i.
func parseFile(filepath string, less func(a, b []byte) bool, delimiter byte) (tokenCount int64, multisetHash int64, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return 0, 0, errors.Wrap(err, "open failed")
	}
	defer func() {
		_ = file.Close()
	}()

	var prevToken *[]byte
	hasher := hash.NewMultiset(17, 1e9+7)

	handleToken := func(token []byte) error {
		tokenCount++
		hasher.Add(hash.Polynomial(token))

		if less == nil {
			return nil
		}

		if prevToken != nil {
			if less(token, *prevToken) {
				log.Printf("wrong order: '%s' goes before '%s'\n", *prevToken, token)
				return errors.New("wrong order")
			}
		}

		prevToken = &token

		return nil
	}

	reader := bufio.NewReaderSize(file, bufferSize)

	for {
		token, err := reader.ReadBytes(delimiter)
		if errors.Is(err, io.EOF) {
			// Handle \n before EOF.
			if len(token) == 0 {
				break
			}

			err = nil
		}
		if err != nil {
			return 0, 0, errors.Wrap(err, "read failed")
		}

		// Remove delimiter if it's written at the end.
		if len(token) > 0 && token[len(token)-1] == delimiter {
			token = token[:len(token)-1]
		}

		err = handleToken(token)
		if err != nil {
			return 0, 0, errors.Wrap(err, "handle failed")
		}
	}

	log.Printf("%s:\n\t%d tokens\n\thash: %d\n\n", filepath, tokenCount, hasher.Hash())

	return tokenCount, hasher.Hash(), nil
}
