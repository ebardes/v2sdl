package static

// go:generate make

// #include "static.h"
import "C"
import (
	"bytes"
	"io"
	"os"
	"time"
	"unsafe"
)

type Item struct {
	Name   string
	Data   []byte
	Length int
	Time   time.Time
}

func HasItem(fn string) bool {
	cstr := C.CString(fn)
	defer C.free(unsafe.Pointer(cstr))
	i := C.hasitem(cstr)
	return i != 0
}

func GetItem(fn string) (item Item, err error) {
	cstr := C.CString(fn)
	defer C.free(unsafe.Pointer(cstr))

	i := C.finditem(cstr)
	if i == nil {
		err = os.ErrNotExist
		return
	}
	len := int(i.len)
	item = Item{
		Name:   fn,
		Length: len,
		Data:   C.GoBytes(i.data, C.int(i.len)),
		Time:   time.Unix(int64(i.time), 0),
	}
	return
}

func (i *Item) Reader() io.Reader {
	return bytes.NewReader(i.Data)
}
