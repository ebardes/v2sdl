package display

import (
	"fmt"
	"v2sdl/config"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

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
	texture   *sdl.Texture
}

type Display struct {
	config.Service
	window   *sdl.Window
	renderer *sdl.Renderer
	fps      gfx.FPSmanager
	master   MasterLayer
	layers   []MediaLayer
	packet   []byte
	current  int
	running  bool
	title    string
	Debug    bool
}

func NewDisplay(title string, cfg config.Config) (d *Display) {
	d = &Display{
		Debug:  cfg.DebugLevel > 2,
		layers: make([]MediaLayer, 5),
		title:  title,
	}

	for i := range d.layers {
		d.layers[i].Init()
	}

	return
}

func (d *Display) Start(*config.Config) (err error) {
	d.window, err = sdl.CreateWindow(d.title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1024, 768, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_ALLOW_HIGHDPI)
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

	gfx.InitFramerate(&d.fps)
	return
}

func (d *Display) Tick() {
	date := ""
	rend := d.renderer
	rect := rend.GetViewport()

	rend.SetDrawColor(d.master.Red.value, d.master.Green.value, d.master.Blue.value, 255)
	rend.FillRect(&rect)

	for n := len(d.layers); n > 0; {
		n--

		layer := d.layers[n]
		if layer.content != nil {
			layer.content.Draw(rend, &layer)
		}
	}

	if d.Debug {
		y := int32(10)
		gfx.StringRGBA(rend, 10, y, string(date), 0, 255, 255, 255)
	}
	rend.Present()
}

func (d *Display) Stop() {
	if d.renderer != nil {
		d.renderer.Destroy()
	}
	if d.window != nil {
		d.window.Destroy()
	}
}

func (d *Display) Name() string { return "Display" }

func (d *MasterLayer) String() string {
	a := []interface{}{
		d.Red,
		d.Green,
		d.Blue,
		d.XPosition,
		d.YPosition,
		d.ScaleX,
		d.ScaleY,
		d.RotateZ,
		d.Mode,
	}
	return fmt.Sprintf("%v", a)
}

func (d *Display) EventLoop() {
	for {
		e := sdl.WaitEventTimeout(25)
		if e != nil {
			et := e.GetType()
			if et == sdl.QUIT {
				break
			}
		}
		d.Tick()
	}
}

func (d *Display) OnFrame(b []byte) {
	if len(b) < 99 {
		return
	}

	d.packet = b
	d.current = 0

	d.master.OnFrame(d)
	for i := range d.layers {
		d.layers[i].OnFrame(d)
	}
	d.master.Background(d)

}

func (d *Display) next() byte {
	if d.current >= len(d.packet) {
		return 0
	}
	c := d.current
	d.current++
	return d.packet[c]
}

func (s *MasterLayer) OnFrame(in *Display) {
	s.Mode.from(in)
	s.XPosition.from(in)
	s.YPosition.from(in)
	s.ScaleX.from(in)
	s.ScaleY.from(in)
	s.RotateZ.from(in)
}

func (s *MasterLayer) Background(in *Display) {
	s.Red.from(in)
	s.Green.from(in)
	s.Blue.from(in)
}

func (s *MediaLayer) OnFrame(in *Display) {
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
	s.Flip.from(in)

	if s.Library.changed() || s.File.changed() {
		group := int(s.Library.get())
		slot := int(s.File.get())

		s.loadContent(group, slot, in.renderer)
	}
}
