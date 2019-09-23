package display

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	fps      gfx.FPSmanager

	Debug bool
}

func NewDisplay(title string) (d *Display, err error) {
	d = &Display{}

	d.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1024, 768, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		err = fmt.Errorf("Failed to create window: %v", err)
		return
	}

	d.renderer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		d.window.Destroy()

		err = fmt.Errorf("Failed to create renderer: %v", err)
		return
	}

	d.Debug = true
	gfx.InitFramerate(&d.fps)

	return
}

func (d *Display) Tick() {
	d.renderer.Clear()

	if d.Debug {
		date := time.Now().String()
		gfx.StringRGBA(d.renderer, 100, 100, date, 255, 255, 255, 255)
	}
	d.renderer.Present()
}

func (d *Display) Close() {
	if d.renderer != nil {
		d.renderer.Destroy()
	}
	if d.window != nil {
		d.window.Destroy()
	}
}

func (d *Display) EventLoop() {
	for {
		e := sdl.WaitEventTimeout(10)
		if e != nil && e.GetType() == sdl.QUIT {
			break
		}

		d.Tick()
	}
}
