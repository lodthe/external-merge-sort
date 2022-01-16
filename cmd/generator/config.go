package main

import (
	"os"
)

type config struct {
	output *os.File

	count int

	minLength         int
	maxLength         int
	equalPrefixLength int

	delimiter byte
	alphabet  []byte
}
