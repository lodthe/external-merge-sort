package buffer

import (
	"os"
)

type Writer struct {
	file   *os.File
	offset int64

	bufIndex int
	buf      []byte

	delimiter byte
}

func NewWriter(file *os.File, offset int64, capacity int, delimiter byte) *Writer {
	return &Writer{
		file:      file,
		offset:    offset,
		buf:       make([]byte, capacity), // TODO: sync.Pool
		delimiter: delimiter,
	}
}

func (w *Writer) Write(data []byte) (err error) {
	for _, b := range data {
		err = w.write(b)
		if err != nil {
			return err
		}
	}

	return w.write(w.delimiter)
}

func (w *Writer) write(b byte) error {
	if w.bufIndex == len(w.buf) {
		err := w.Flush()
		if err != nil {
			return err
		}
	}

	w.buf[w.bufIndex] = b
	w.bufIndex++

	return nil
}

func (w *Writer) Flush() error {
	var written int
	for written < w.bufIndex {
		n, err := w.file.WriteAt(w.buf[:w.bufIndex], w.offset)
		w.offset += int64(n)
		written += n

		if err != nil {
			return err
		}
	}

	w.bufIndex = 0

	return nil
}
