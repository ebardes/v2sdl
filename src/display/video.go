package display

import (
	"v2sdl/video"

	"github.com/veandco/go-sdl2/sdl"
)

type Video struct {
	Image
	filename string
	mp       *video.AVMediaPlayer
}

func NewVideoContent(fn string, r *sdl.Renderer) (i Content, err error) {
	mp, err := video.NewAVMediaPlayer()
	if err != nil {
		return
	}

	v := &Video{
		filename: fn,
		mp:       mp,
	}

	err = mp.Open(fn)
	if err != nil {
		return
	}

	v.rect = sdl.Rect{X: 0, Y: 0, W: int32(mp.Width), H: int32(mp.Height)}
	v.texture, err = r.CreateTexture(sdl.PIXELFORMAT_YV12, sdl.TEXTUREACCESS_STREAMING, int32(mp.Width), int32(mp.Height))
	i = v
	return
}

func (v *Video) Destroy() {
	v.mp.Dispose()
	v.Image.Destroy()
}

func (v *Video) Start() {
	if !v.mp.Running {
		v.mp.Running = true
		go v.mp.Run()
	}
}

func (v *Video) Stop() {
	v.mp.Running = false
}
