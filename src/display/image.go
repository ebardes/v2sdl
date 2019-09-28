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

	/*
	 * calculate aspect ratio adjustments
	 */
	calcH1 := dest.W * i.rect.H / i.rect.W
	calcW1 := dest.H * i.rect.W / i.rect.H
	if calcH1 < dest.H {
		diff := dest.H - calcH1
		dest.Y = diff / 2
		dest.H = calcH1
	} else if calcW1 < dest.W {
		diff := dest.W - calcW1
		dest.X = diff / 2
		dest.W = calcW1
	}

	/*
	 * calculate positioning offsets
	 */
	x1 := int32(layer.XPosition.value) / 32
	y1 := int32(layer.YPosition.value) / 32
	dest.X += x1
	dest.Y += y1

	/*
	 * calculate scaling
	 */

	sx1 := (int32(layer.ScaleX.value) + 32768)
	sy1 := (int32(layer.ScaleY.value) + 32768)

	dest.X -= dest.W * (sx1 - 16384) / 32768
	dest.W = dest.W * sx1 / 16384
	dest.Y -= dest.H * (sy1 - 16384) / 32768
	dest.H = dest.H * sy1 / 16384

	center := sdl.Point{X: (dest.W / 2) + dest.X, Y: (dest.H / 2) + dest.Y}
	angle := float64(layer.RotateZ.get()) * 180.0 / 32768.0
	i.texture.SetAlphaMod(layer.Intensity.get())
	r.CopyEx(i.texture, &i.rect, &dest, angle, &center, flip)
}
