package video

import (
	"fmt"
	"unsafe"
)

// #cgo CFLAGS: -I ../../av/include
// #cgo LDFLAGS: -L ../../av/lib -l bz2 -l z -l avdevice -l avutil -l avcodec -l avresample -l avfilter -l avformat -l swscale
// #include "av.h"
import "C"

type AVMediaPlayer struct {
	Debug bool
	afctx *C.AVFormatContext
}

func init() {
	C.av_register_all()
}

func NewAVMediaPlayer() (amp *AVMediaPlayer, err error) {
	amp = &AVMediaPlayer{}
	amp.afctx = C.avformat_alloc_context()
	return
}

func (amp *AVMediaPlayer) Dispose() {
	C.avformat_free_context(amp.afctx)
}

func (amp *AVMediaPlayer) Open(fn string) (err error) {
	cfn := C.CString(fn)
	defer C.free(unsafe.Pointer(cfn))

	status := int(C.avformat_open_input(&amp.afctx, cfn, nil, nil))
	if status != 0 {
		err = amp.geterror(status)
		return
	}

	return
}

func (amp *AVMediaPlayer) geterror(n int) error {
	buffer := (*C.char)(C.malloc(1000))
	defer C.free(unsafe.Pointer(buffer))

	C.av_strerror(C.int(n), buffer, 999)
	str := C.GoString(buffer)

	return fmt.Errorf("%s", str)
}
