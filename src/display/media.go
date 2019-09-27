package display

import (
	"v2sdl/config"

	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
)

type Content interface {
	Draw(r *sdl.Renderer, layer *MediaLayer)
	Destroy()
	Start()
	Stop()
}

type MediaLayer struct {
	Intensity  single
	Library    single
	File       single
	Volume     single
	XPosition  double
	YPosition  double
	ScaleX     double
	ScaleY     double
	RotateZ    double
	Brightness single
	Contrast   single
	Playmode   single
	Flip       single

	content Content
	texture *sdl.Texture
}

func (ml *MediaLayer) loadContent(r *sdl.Renderer) (err error) {
	go func() {
		if ml.content != nil {
			ml.content.Destroy()
			ml.content = nil
		}

		if ml.texture != nil {
			ml.texture.Destroy()
			ml.texture = nil
		}

		group := int(ml.Library.get())
		slot := int(ml.File.get())
		log.Debug().Int("group", group).Int("slot", slot).Msg("Media Change")
		item := config.Media.Get(group, slot)
		if item != nil {
			fn := item.Path()

			log.Debug().Str("fn", fn).Msg("Attempting to load file")

			ml.content, err = NewImageContent(fn, r)
			if err != nil {
				log.Error().Err(err).Msg("Error loading content")
			}
		} else {
			log.Debug().Msg("No item found")
		}
	}()
	return
}
