package buffer

import (
	"io"
	"os"
)

type Reader struct {
	file      *os.File
	metEOF    bool
	offset    int64
	endOffset int64

	bufIndex int
	bufLen   int
	buf      []byte

	delimiter byte
}

func NewReader(f *os.File, offset, endOffset int64, capacity int, delimiter byte) *Reader {
	return &Reader{
		file:      f,
		metEOF:    false,
		offset:    offset,
		endOffset: endOffset,
		buf:       make([]byte, capacity), // TODO: sync.Pool
		delimiter: delimiter,
	}
}

// Next reads bytes and stops when it finds the delimiter or EOF.
// If EOF has occurred and no data can be read, (nil, io.EOF) is returned.
// Otherwise, Next returns a non-nil slice (without the delimiter).
func (r *Reader) Next() (token []byte, err error) {
	if r.EOF() {
		return nil, io.EOF
	}

	data := make([]byte, 0, 8)

	for {
		if r.bufIndex == r.bufLen {
			err := r.read()
			if err != nil {
				return nil, err
			}

			if r.EOF() {
				if len(data) == 0 {
					return nil, io.EOF
				}

				break
			}
		}

		c := r.buf[r.bufIndex]
		r.bufIndex++

		if c == r.delimiter {
			break
		}

		data = append(data, c)
	}

	return data, nil
}

func (r *Reader) EOF() bool {
	return r.metEOF && r.bufIndex == r.bufLen
}

func (r *Reader) read() error {
	r.bufIndex = 0
	r.bufLen = 0

	for r.bufLen < len(r.buf) && r.offset < r.endOffset {
		maxN := int64(len(r.buf) - r.bufLen)
		if maxN > r.endOffset-r.offset {
			maxN = r.endOffset - r.offset
		}

		n, err := r.file.ReadAt(r.buf[r.bufLen:r.bufLen+int(maxN)], r.offset)
		r.bufLen += n
		r.offset += int64(n)

		if err == io.EOF {
			r.metEOF = true
			break
		}
		if err != nil {
			return err
		}
	}

	if r.offset >= r.endOffset {
		r.metEOF = true
	}

	return nil
}
