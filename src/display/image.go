package display

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Image struct {
	Content
	surface *sdl.Surface
	Width   int32
	Height  int32
}

func NewImageContent(fn string) (i *Image, err error) {
	s, err := img.Load(fn)
	if err != nil {
		return
	}

	log.Debug().Msgf("surface{%d, %d}", s.W, s.H)
	i = &Image{surface: s}
	return
}
func (i *Image) Surface() *sdl.Surface { return i.surface }

func (i *Image) Destroy() {
	if i.surface != nil {
		i.surface.Free()
		i.surface = nil
	}
}
