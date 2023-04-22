package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth
)

// Index represents a log file that can be appended record indexes to and read from.
//
// The `file` pointer is the file that will be used to write and read.
// It `mmap` is a Memory mapping that allows us to efficiently i/o operations in the file that it map to.
// `size` is the amount of bytes in the `file`.
type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

// `newIndex` creates a `index` struct with the `f` file pointer and trucates the its size basedo on `c.Segment.MaxIndexBytes`.
// It also create a Memory mapping, and returns a index pointer or an error.
func newIndex(f *os.File, c Config) (*index, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx := &index{
		file: f,
		size: uint64(fi.Size()),
	}

	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}

	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_NONE|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}

	return idx, nil
}

// `Close` method closes the file index.
// Before it closed, it flush the `mmap` and `file` data to the disk.
// It also truncate the file size with the current `index.size`.
func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}

	return i.file.Close()
}

// `Read method` reads the file data from the memory mapping based on the `in` argument.
// It returns the offset (`out`), position of the record (`pos`), and a error if it occurred.
func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if in == -1 {
		out = uint32((i.size / entWidth) - 1)
	} else {
		out = uint32(in)
	}
	pos = uint64(out) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF
	}
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])
	return out, pos, nil
}

// `Write` method writes the record's offset (`pos`) and its position (`pos`) in index file.
// It returns immediately a error if it occurred.
func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}

	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entWidth], pos)
	i.size += uint64(entWidth)

	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}
