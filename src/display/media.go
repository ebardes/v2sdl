package display

import (
	"sync"
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

	content   Content
	texture   *sdl.Texture
	loadmutex *sync.Mutex
}

// Init initializes some of the synchronization sensitive members
func (ml *MediaLayer) Init() {
	ml.loadmutex = &sync.Mutex{}
}

func (ml *MediaLayer) loadContent(group, slot int, r *sdl.Renderer) (err error) {
	go func() {
		ml.loadmutex.Lock()
		defer ml.loadmutex.Unlock()

		if ml.content != nil {
			ml.content.Destroy()
			ml.content = nil
		}

		if ml.texture != nil {
			ml.texture.Destroy()
			ml.texture = nil
		}

		item := config.Media.Get(group, slot)
		if item != nil {
			fn := item.Path()

			log.Debug().Str("fn", fn).Msg("Attempting to load file")

			ml.content, err = NewImageContent(fn, r)
			if err != nil {
				log.Error().Err(err).Msg("Error loading content")
			}
		} else {
			log.Debug().Int("group", group).Int("slot", slot).Msg("No item found")
		}
	}()
	return
}
