package display

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Image struct {
	Content
	texture *sdl.Texture
	rect    sdl.Rect
}

func NewImageContent(fn string, r *sdl.Renderer) (i *Image, err error) {
	s, err := img.Load(fn)
	if err != nil {
		return
	}
	rect := sdl.Rect{X: 0, Y: 0, W: s.W, H: s.H}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return
	}
	t.SetBlendMode(sdl.BLENDMODE_BLEND)

	i = &Image{
		texture: t,
		rect:    rect,
	}
	return
}

func (i *Image) Destroy() {
	if i.texture != nil {
		i.texture.Destroy()
		i.texture = nil
	}
}

func (i *Image) Draw(r *sdl.Renderer, layer *MediaLayer) {
	flip := sdl.FLIP_NONE
	switch layer.Flip.value {
	case 1:
		flip = sdl.FLIP_HORIZONTAL
	case 2:
		flip = sdl.FLIP_VERTICAL
	case 3:
		flip = sdl.FLIP_HORIZONTAL | sdl.FLIP_VERTICAL
	}

	dest := r.GetViewport()
	center := sdl.Point{X: dest.W / 2, Y: dest.H / 2}
	angle := float64(layer.RotateZ.get()) * 180.0 / 32768.0
	i.texture.SetAlphaMod(layer.Intensity.get())
	r.CopyEx(i.texture, &i.rect, &dest, angle, &center, flip)
}
