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
	window   *sdl.Window
	renderer *sdl.Renderer
	fps      gfx.FPSmanager
	master   MasterLayer
	layers   []MediaLayer
	packet   []byte
	current  int

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
	date := ""
	rend := d.renderer
	rect := rend.GetViewport()

	rend.SetDrawColor(d.master.Red.value, d.master.Green.value, d.master.Blue.value, 255)
	rend.FillRect(&rect)

	for n := len(d.layers); n > 0; {
		n--

		layer := d.layers[n]
		intensity := layer.Intensity.get()
		c := layer.content
		if c != nil && intensity > 0 {
			s := c.Surface()
			if layer.texture == nil {
				layer.texture, _ = rend.CreateTextureFromSurface(s)
				layer.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
			}
			src := sdl.Rect{X: 0, Y: 0, W: s.W, H: s.H}
			t := layer.texture
			if intensity < 255 {
				t.SetAlphaMod(intensity)
			} else {
				t.SetBlendMode(sdl.BLENDMODE_NONE)
			}
			rend.Copy(t, &src, &rect)
		}
	}

	if d.Debug {
		y := int32(10)
		// date := d.master.String()
		gfx.StringRGBA(rend, 10, y, string(date), 0, 255, 255, 255)

		// for _, layer := range d.layers {
		// 	y += 20
		// 	date, _ := json.Marshal(layer)
		// 	gfx.StringRGBA(rend, 10, y, string(date), 255, 255, 255, 255)
		// }
	}
	rend.Present()
}

func (d *Display) Close() {
	if d.renderer != nil {
		d.renderer.Destroy()
	}
	if d.window != nil {
		d.window.Destroy()
	}
}

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
		e := sdl.WaitEventTimeout(100)
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

	if s.Library.changed() || s.File.changed() {
		s.loadContent()
	}
}
