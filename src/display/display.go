package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"v2sdl/config"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type single uint8
type double int16

type MasterLayer struct {
	Red       single
	Green     single
	Blue      single
	XPosition double
	YPosition double
	ScaleX    double
	ScaleY    double
	RotateZ   double
	Mode      single
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
}

type Display struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	fps      gfx.FPSmanager
	state    MasterLayer
	layers   []MediaLayer

	Debug bool
}

func NewDisplay(title string, cfg config.Config) (d *Display, err error) {
	d = &Display{Debug: cfg.DebugLevel > 2}
	d.layers = make([]MediaLayer, 5)

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

	d.Debug = cfg.DebugLevel > 0
	gfx.InitFramerate(&d.fps)

	return
}

func (d *Display) Tick() {
	d.renderer.Clear()

	if d.Debug {
		y := int32(10)
		date, _ := json.Marshal(d.state)
		gfx.StringRGBA(d.renderer, 10, y, string(date), 255, 255, 255, 255)

		for _, layer := range d.layers {
			y += 20
			date, _ := json.Marshal(layer)
			gfx.StringRGBA(d.renderer, 10, y, string(date), 255, 255, 255, 255)
		}
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

func (d *Display) OnFrame(b []byte) {
	if len(b) < 99 {
		return
	}

	in := bytes.NewReader(b)
	d.state.OnFrame(in)
	for i := range d.layers {
		d.layers[i].OnFrame(in)
	}
	d.state.Background(in)
}

func (s *MasterLayer) OnFrame(in *bytes.Reader) {
	s.Mode.from(in)
	s.XPosition.from(in)
	s.YPosition.from(in)
	s.ScaleX.from(in)
	s.ScaleY.from(in)
	s.RotateZ.from(in)
}

func (s *MasterLayer) Background(in *bytes.Reader) {
	s.Red.from(in)
	s.Green.from(in)
	s.Blue.from(in)
}

func (s *MediaLayer) OnFrame(in *bytes.Reader) {
	s.Intensity.from(in)
	s.Library.from(in)
	s.File.from(in)
	s.Volume.from(in)
	s.XPosition.from(in)
	s.YPosition.from(in)
	s.ScaleX.from(in)
	s.ScaleY.from(in)
	s.RotateZ.from(in)
	s.Brightness.from(in)
	s.Contrast.from(in)
	s.Playmode.from(in)
}

func (d *single) from(r *bytes.Reader) {
	b, _ := r.ReadByte()
	*d = single(b)
}

func (d *double) from(r *bytes.Reader) {
	ub, _ := r.ReadByte()
	lb, _ := r.ReadByte()
	w := int(ub)<<8 | int(lb)
	*d = double(w + 32768)
}
