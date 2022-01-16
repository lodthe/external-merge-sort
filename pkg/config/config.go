package config

type Config struct {
	// Size of one block is bytes.
	BlockSize int

	// How much memory the program can waste.
	MemoryLimit int

	// Delimiter separates one token from another.
	Delimiter byte

	// Less determines whether the first token must be presented earlier than the second one.
	Less func(a, b []byte) bool
}
