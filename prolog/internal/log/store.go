package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8
)

// Store represents a log file that can be appended record to and read record from.
//
// Embedding the `*os.File` pointer in the struct allows the `store` struct to access all the fields and methods of the `os.File` type, effectively inheriting its functionaly.
// This makes it easy to use function from the `os` package to interact with the file, while also adding additional functionality for appending and reading data.
//
// The `sync.Mutex` is used to ensure that only one goroutine can access the file at a time, preventing concurrent writes from interfering with each other.
// The `*bufio.Writer` is used to buffer writes to the file, which can improve performance by reducing the number of system calls needed to write data.
// The `size` fiel represents the current size of the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

// The `newStore` creates a `store` based on `*os.File` which represents a file that has already been opened.
//
// It returns a pointer to a `store` or a error if occurs.
func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	size := uint64(fi.Size())

	return &store{
		File: f,
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil
}

// The `Append` method is used to write data to the end of the store.
// It takes a byte slice `p` as its argument, which represents the data to be written.
// It returns the number of bytes written (`n`), the position in the file where the data was written (`pos`), and any errors that occurred during the write (`err`)/
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size
	if err = binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return
	}

	nn, err := s.buf.Write(p)
	if err != nil {
		return n, pos, err
	}

	n = lenWidth + uint64(nn)
	s.size += uint64(n)

	return
}

// Read method reads data from the store at the specified position (`pos`).
// It returns the buffer that holds the data. If there is an error reading the data, the method returns it.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}
	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return b, nil
}

// The ReadAt reads data from the store at the specified offset.
// It takes a byte slice p and an off argument as the offset from which to read the data.
// It returns the number of bytes read and any error that occurred.
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}

// The Close method closes the file that is represented by the store.
// Before closing it, the method flushes the buffer to ensure that any data in the buffer is written to the fie.
// It returns an error if occurs during the closed.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return err
	}

	return s.File.Close()
}
