package video

import (
	"fmt"
	"unsafe"

	"github.com/rs/zerolog/log"
)

// #cgo CFLAGS: -I ../../av/include
// #cgo LDFLAGS: -L ../../av/lib -l bz2 -l z -l avdevice -l avutil -l avcodec -l avresample -l avfilter -l avformat -l swscale
// #include "av.h"
import "C"

type avdemux struct {
	index int
	codec *C.AVCodec
	ctx   *C.AVCodecContext
}

type AVMediaPlayer struct {
	Debug   bool
	Running bool
	afctx   *C.AVFormatContext
	video   avdemux
	audio   avdemux
	frame   *C.AVFrame
	packet  *C.AVPacket
	Width   int
	Height  int
}

func init() {
	C.av_register_all()
}

func NewAVMediaPlayer() (amp *AVMediaPlayer, err error) {
	amp = &AVMediaPlayer{}
	amp.afctx = C.avformat_alloc_context()
	amp.Running = true
	return
}

func (amp *AVMediaPlayer) Dispose() {
	amp.video.dispose()
	amp.audio.dispose()
	C.av_packet_free(&amp.packet)
	C.av_frame_free(&amp.frame)
	C.avformat_close_input(&amp.afctx)
	C.avformat_free_context(amp.afctx)
}

func (d *avdemux) dispose() {
	C.avcodec_free_context(&d.ctx)
}

// Open
func (amp *AVMediaPlayer) Open(fn string) (err error) {
	cfn := C.CString(fn)
	defer C.free(unsafe.Pointer(cfn))

	status := int(C.avformat_open_input(&amp.afctx, cfn, nil, nil))
	if status != 0 {
		err = amp.geterror("open input", status)
		return
	}

	ifmt := amp.afctx.iformat
	log.Info().Msgf("format = %s", C.GoString(ifmt.long_name))

	status = int(C.avformat_find_stream_info(amp.afctx, nil))
	if status != 0 {
		err = amp.geterror("find stream info", status)
		return
	}

	// loop though all the streams and print its main information
	for i := 0; i < int(amp.afctx.nb_streams); i++ {
		streams := amp.afctx.streams
		ptrPtr := (**C.AVStream)(unsafe.Pointer(uintptr(unsafe.Pointer(streams)) + uintptr(i)*unsafe.Sizeof(*streams)))
		lc := (*ptrPtr).codecpar

		var common *avdemux

		switch lc.codec_type {
		case C.AVMEDIA_TYPE_VIDEO:
			log.Debug().Int("stream", i).Msgf("Params: Video %dx%d", lc.width, lc.height)
			common = &amp.video

			amp.Width = int(lc.width)
			amp.Height = int(lc.height)

		case C.AVMEDIA_TYPE_AUDIO:
			log.Debug().Int("stream", i).Msgf("Params: Audio %d channels %d samples", lc.channels, lc.sample_rate)
			common = &amp.audio

		default:
			continue
		}

		common.index = i
		common.codec = C.avcodec_find_decoder(lc.codec_id)
		common.ctx = C.avcodec_alloc_context3(common.codec)
		C.avcodec_parameters_to_context(common.ctx, lc)

		status = int(C.avcodec_open2(common.ctx, common.codec, nil))
		if status != 0 {
			err = amp.geterror("codec open", status)
			return
		}

	}
	amp.frame = C.av_frame_alloc()
	amp.packet = C.av_packet_alloc()

	return
}

func (amp *AVMediaPlayer) Run() error {
	frames := 1000
	for amp.Running {
		status := int(C.av_read_frame(amp.afctx, amp.packet))

		if status != 0 {
			return amp.geterror("read frame", status)
		}

		if int(amp.packet.stream_index) == amp.video.index {
			status = int(C.avcodec_send_packet(amp.video.ctx, amp.packet))
			if status != 0 {
				return amp.geterror("send packet", status)
			}

			for status >= 0 {
				status = int(C.avcodec_receive_frame(amp.video.ctx, amp.frame))
				if status == C.AVERROR_EOF {
					break
				}
			}
		}

		frames--
		if frames <= 0 {
			amp.Running = false
		}

		C.av_packet_unref(amp.packet)
	}
	return nil
}

func (amp *AVMediaPlayer) geterror(msg string, n int) error {
	buffer := (*C.char)(C.malloc(1000))
	defer C.free(unsafe.Pointer(buffer))

	C.av_strerror(C.int(n), buffer, 999)
	str := C.GoString(buffer)

	return fmt.Errorf("%s: %s", msg, str)
}
