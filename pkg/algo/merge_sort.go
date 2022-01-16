package algo

import (
	"io"
	"log"
	"os"
	"sort"
	"time"
	"unsafe"

	"github.com/lodthe/external-merge-sort/pkg/buffer"
	"github.com/lodthe/external-merge-sort/pkg/config"
	"github.com/pkg/errors"
)

const MaxUint64 = ^uint64(0)
const MaxInt64 = int64(MaxUint64 / 2)

type mergeSortBlock struct {
	start int64
	end   int64
}

// ExternalMergeSort is an implementation of external merge sort algorithm with 2-way merge.
type ExternalMergeSort struct {
	input  *os.File
	output *os.File

	cfg *config.Config

	tokenCount int64
}

func NewExternalMergeSort(cfg *config.Config) *ExternalMergeSort {
	return &ExternalMergeSort{
		cfg: cfg,
	}
}

// Sort loads data from the input file, sorts it and saves result to the output file.
func (m *ExternalMergeSort) Sort(inputPath, outputPath, tempDir string) error {
	originInput, err := m.createDescriptors(inputPath, tempDir)
	if err != nil {
		return errors.Wrap(err, "failed open basic files")
	}
	defer m.finish(outputPath)

	blocks, err := m.mainMemorySort(originInput, m.input)
	if err != nil {
		return errors.Wrap(err, "failed to sort blocks in RAM")
	}

	err = m.externalSort(blocks)
	if err != nil {
		return errors.Wrap(err, "failed to sort externally")
	}

	return nil
}

func (m *ExternalMergeSort) swapDescriptors() {
	m.input, m.output = m.output, m.input
}

func (m *ExternalMergeSort) createDescriptors(inputPath string, tempDir string) (originInput *os.File, err error) {
	originInput, err = os.Open(inputPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open input file")
	}
	defer func() {
		if err == nil {
			return
		}

		if originInput != nil {
			_ = originInput.Close()
		}
		if m.input != nil {
			_ = m.input.Close()
		}
		if m.output != nil {
			_ = m.output.Close()
		}
	}()

	const tempPattern = "external_merge_sort_*"

	m.input, err = os.CreateTemp(tempDir, tempPattern)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file")
	}

	m.output, err = os.CreateTemp(tempDir, tempPattern)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp file")
	}

	return originInput, nil
}

// finish closes files and renames temp output with the correct name.
func (m *ExternalMergeSort) finish(outputPath string) {
	// After the external merge sort process, the output is stored in m.input.
	// So we swap descriptors here to make the code cleaner.
	m.swapDescriptors()

	err := m.input.Close()
	if err != nil {
		log.Printf("failed to close temp input file: %v\n", err)
	}

	err = os.Remove(m.input.Name())
	if err != nil {
		log.Printf("failed to remove temp input file: %v\n", err)
	}

	err = m.output.Sync()
	if err != nil {
		log.Printf("filaed to sync temp output file: %v\n", err)
	}

	err = m.output.Close()
	if err != nil {
		log.Printf("failed to close temp output file: %v\n", err)
	}

	err = os.Rename(m.output.Name(), outputPath)
	if err != nil {
		log.Printf("failed to rename %s temp file to %s output file: %v\n", m.output.Name(), outputPath, err)
	}
}

// mainMemorySort reads data into buffer of size M and sorts them in main memory.
func (m *ExternalMergeSort) mainMemorySort(input, output *os.File) ([]mergeSortBlock, error) {
	log.Printf("main memory sort started...\n")

	startedAt := time.Now()
	r := buffer.NewReader(input, 0, MaxInt64, m.cfg.BlockSize, m.cfg.Delimiter)
	w := buffer.NewWriter(output, 0, m.cfg.BlockSize, m.cfg.Delimiter)

	var offset int64
	var tokenCapacityTotal int
	var blocks []mergeSortBlock
	var tokens [][]byte

	// writeTokens sorts portion of tokens and writes them.
	writeTokens := func() (err error) {
		if len(tokens) == 0 {
			return nil
		}

		// TODO: a parallel sort can be used here.
		sort.Slice(tokens, func(i, j int) bool {
			return m.cfg.Less(tokens[i], tokens[j])
		})

		var writtenCnt int64
		for _, t := range tokens {
			err = w.Write(t)
			if err != nil {
				return err
			}

			// Length of token + delimiter.
			writtenCnt += int64(len(t)) + 1
		}

		blocks = append(blocks, mergeSortBlock{
			start: offset,
			end:   offset + writtenCnt,
		})
		offset += writtenCnt

		tokens = tokens[:0]
		tokenCapacityTotal = 0

		return nil
	}

	// Read a token while it exists. When current memory usage is too high, sort tokens and write them.
	for !r.EOF() {
		m.tokenCount++
		token, err := r.Next()
		if errors.Is(err, io.EOF) {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to read the next token")
		}

		tokens = append(tokens, token)
		tokenCapacityTotal += cap(token)

		// count total overhead of slices for storing tokens and total capacity of tokens.
		currentUsage := int(unsafe.Sizeof(token))*cap(tokens) + tokenCapacityTotal
		if currentUsage >= m.cfg.MemoryLimit/2 {
			err = writeTokens()
			if err != nil {
				return nil, err
			}
		}
	}

	err := writeTokens()
	if err != nil {
		return nil, err
	}

	err = w.Flush()
	if err != nil {
		return nil, errors.Wrap(err, "final flush failed")
	}

	log.Printf("main memory sort finished in %v\n\n", time.Since(startedAt))

	return blocks, nil
}

func (m *ExternalMergeSort) externalSort(blocks []mergeSortBlock) error {
	startedAt := time.Now()

	log.Printf("external sort started...\n")

	var iterations int
	for len(blocks) > 1 {
		iterations++

		writer := buffer.NewWriter(m.output, 0, m.cfg.BlockSize, m.cfg.Delimiter)
		newBlocks := make([]mergeSortBlock, 0, len(blocks)/2+1)

		if len(blocks)%2 == 1 {
			blocks = append(blocks, mergeSortBlock{
				start: blocks[len(blocks)-1].end,
				end:   blocks[len(blocks)-1].end,
			})
		}

		for i := 0; i+1 < len(blocks); i += 2 {
			err := m.merge(blocks[i], blocks[i+1], writer)
			if err != nil {
				log.Printf("iteration #%d failed: %v\n", iterations, err)
				return errors.Wrap(err, "merge failed")
			}

			newBlocks = append(newBlocks, mergeSortBlock{
				start: blocks[i].start,
				end:   blocks[i+1].end,
			})
		}

		err := writer.Flush()
		if err != nil {
			return errors.Wrap(err, "flush failed")
		}

		blocks = newBlocks
		m.swapDescriptors()

		log.Printf("iteration #%d finished, %d blocks left\n", iterations, len(blocks))
	}

	log.Printf("external sort finished in %d iterations (%v)\n\n", iterations, time.Since(startedAt))

	return nil
}

// merge merges two blocks.
// TODO: use K-way merge, it's much faster.
func (m *ExternalMergeSort) merge(a, b mergeSortBlock, writer *buffer.Writer) error {
	readerA := buffer.NewReader(m.input, a.start, a.end, m.cfg.BlockSize, m.cfg.Delimiter)
	readerB := buffer.NewReader(m.input, b.start, b.end, m.cfg.BlockSize, m.cfg.Delimiter)

	tokenA, err := readerA.Next()
	if err != nil && err != io.EOF {
		return err
	}

	tokenB, err := readerB.Next()
	if err != nil && err != io.EOF {
		return err
	}

	writeAndReadNext := func(token *[]byte, reader *buffer.Reader) error {
		err = writer.Write(*token)
		if err != nil {
			return errors.Wrap(err, "write failed")
		}

		*token, err = reader.Next()
		if err != nil && err != io.EOF {
			return err
		}

		return nil
	}

	for tokenA != nil && tokenB != nil {
		if m.cfg.Less(tokenA, tokenB) {
			err = writeAndReadNext(&tokenA, readerA)
		} else {
			err = writeAndReadNext(&tokenB, readerB)
		}

		if err != nil {
			return err
		}
	}

	for tokenA != nil {
		err = writeAndReadNext(&tokenA, readerA)
		if err != nil {
			return err
		}
	}

	for tokenB != nil {
		err = writeAndReadNext(&tokenB, readerB)
		if err != nil {
			return err
		}
	}

	return nil
}
