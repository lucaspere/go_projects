package log

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var (
	write = []byte("hello world")
	width = uint64(len(write)) + lenWidth
)

// func TestStoreAppendRead(t *testing.T) {

// }

func TestStoreAppendRead(t *testing.T) {
	f, err := ioutil.TempFile("", "store_append_rest_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	store, err := newStore(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Testing Append", func(t *testing.T) {
		t.Helper()

		for i := uint64(1); i < 4; i++ {
			n, pos, err := store.Append(write)
			if err != nil {
				t.Fatal(err)
			}
			if pos+n != width*i {
				t.Errorf("Expect to be equal %d, %d", pos+n, width*i)
			}
		}
	})
	t.Run("Testing Read", func(t *testing.T) {
		t.Helper()
		var pos uint64

		for i := uint64(1); i < 4; i++ {
			read, err := store.Read(pos)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(read, write) {
				t.Errorf("Expect %v to be equal %v", read, write)
			}
		}
	})

	t.Run("Testing ReadAt", func(t *testing.T) {
		t.Helper()
		for i, off := uint64(1), int64(0); i < 4; i++ {
			b := make([]byte, lenWidth)
			n, err := store.ReadAt(b, off)
			if err != nil {
				t.Fatal(err)
			}
			if n != lenWidth {
				t.Errorf("Expect %v to be equal %v", n, lenWidth)
			}
			off += int64(n)

			size := enc.Uint64(b)
			b = make([]byte, size)
			n, err = store.ReadAt(b, off)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(write, b) {
				t.Errorf("Expect %v to be equal %v", write, b)
			}
			if int(size) != n {
				t.Errorf("Expect %v to be equal %v", int(size), lenWidth)
			}
			off += int64(n)
		}
	})
}

func TestStoreClose(t *testing.T) {
	f, err := ioutil.TempFile("", "store_close_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	s, err := newStore(f)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = s.Append(write)
	if err != nil {
		t.Fatal(err)
	}
	f, beforeSize, err := openFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	err = s.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, afterSize, err := openFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if afterSize < beforeSize {
		t.Errorf("Expect %v to be greater than %v", afterSize, beforeSize)
	}
}

func openFile(name string) (file *os.File, size int64, err error) {
	f, err := os.OpenFile(
		name,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, 0, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	return f, fi.Size(), nil
}
