package main

import (
	"log"
	"v2sdl/config"
	"v2sdl/display"
	"v2sdl/dmx"
	"v2sdl/dmx/artnet"
	"v2sdl/dmx/sacn"
	"v2sdl/web"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	w := web.WebServer{}
	w.Start(&cfg)

	var net dmx.NetDMX

	switch cfg.Protocol {
	default:
		net, err = sacn.NewService(&cfg)

	case "artnet":
		net, err = artnet.NewService(&cfg)
	}

	if err != nil {
		panic(err)
	}

	go net.Run()
	defer net.Stop()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.Main(func() {})

	d, err := display.NewDisplay("Media Server", cfg)
	if err != nil {
		log.Panic(err)
	}
	defer d.Close()

	net.SetDisplay(d)

	d.EventLoop()
}
