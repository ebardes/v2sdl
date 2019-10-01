package display

import (
	"fmt"
	"unsafe"

	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/veandco/go-sdl2/sdl"
)

type Video struct {
	Image
	filename string

	pFormatContext *avformat.Context
	pCodecCtxOrig  *avformat.CodecContext
	pCodecCtx      *avcodec.Context
	pCodec         *avcodec.Codec
}

func NewVideoContent(fn string, r *sdl.Renderer) (v *Video, err error) {
	v = &Video{
		filename: fn,
	}

	// Open video file
	v.pFormatContext = avformat.AvformatAllocContext()
	if avformat.AvformatOpenInput(&v.pFormatContext, fn, nil, nil) != 0 {
		err = fmt.Errorf("Unable to open file %s\n", fn)
		return
	}

	// Retrieve stream information
	if v.pFormatContext.AvformatFindStreamInfo(nil) < 0 {
		err = fmt.Errorf("Couldn't find stream information")
		return
	}

	// Find the first video stream
	for i := 0; i < int(v.pFormatContext.NbStreams()); i++ {
		switch v.pFormatContext.Streams()[i].CodecParameters().AvCodecGetType() {
		case avformat.AVMEDIA_TYPE_VIDEO:

			// Get a pointer to the codec context for the video stream
			v.pCodecCtxOrig = v.pFormatContext.Streams()[i].Codec()

			// Find the decoder for the video stream
			v.pCodec = avcodec.AvcodecFindDecoder(avcodec.CodecId(v.pCodecCtxOrig.GetCodecId()))
			if v.pCodec == nil {
				err = fmt.Errorf("Unsupported codec!")
				return
			}
			// Copy context
			v.pCodecCtx = v.pCodec.AvcodecAllocContext3()
			if v.pCodecCtx.AvcodecCopyContext((*avcodec.Context)(unsafe.Pointer(v.pCodecCtxOrig))) != 0 {
				err = fmt.Errorf("Couldn't copy codec context")
				return
			}

			// Open codec
			if v.pCodecCtx.AvcodecOpen2(v.pCodec, nil) < 0 {
				err = fmt.Errorf("Could not open codec")
				return
			}
		}
	}

	return
}

func (v *Video) Destroy() {
	// Close the codecs
	if v.pCodecCtx != nil {
		v.pCodecCtx.AvcodecClose()
		v.pCodecCtx = nil
	}
	if v.pCodecCtxOrig != nil {
		(*avcodec.Context)(unsafe.Pointer(v.pCodecCtxOrig)).AvcodecClose()
		v.pCodecCtxOrig = nil
	}

	// Close the video file
	if v.pFormatContext != nil {
		v.pFormatContext.AvformatCloseInput()
		v.pFormatContext = nil
	}
}
