package main

import (
	"os"
	"v2sdl/config"
	"v2sdl/display"
	"v2sdl/dmx"
	"v2sdl/dmx/artnet"
	"v2sdl/dmx/sacn"
	"v2sdl/web"

	flags "github.com/jessevdk/go-flags"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
)

type options struct {
	Config string `long:"config" short:"c" long:"location of config.json file"`
}

func main() {
	opts := options{}

	tty := isatty.IsTerminal(os.Stderr.Fd())
	if tty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	cfg, err := config.Load(opts.Config)
	if err != nil {
		panic(err)
	}
	defer cfg.StopAll()

	w := web.WebServer{}
	cfg.AddAndStartService(&w)

	var net dmx.NetDMX

	switch cfg.Protocol {
	default:
		net = sacn.NewService()

	case "artnet":
		net, err = artnet.NewService(&cfg)
	}

	if err != nil {
		panic(err)
	}

	cfg.AddAndStartService(net)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.Main(func() {})

	d := display.NewDisplay("Media Server", cfg)
	cfg.AddAndStartService(d)

	net.SetDisplay(d)

	d.EventLoop()
}
